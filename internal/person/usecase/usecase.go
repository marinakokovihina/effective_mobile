package usecase

import (
	"context"
	"effective_mobile/config"
	"effective_mobile/internal/person"
	"effective_mobile/internal/person/model"
	"effective_mobile/internal/person/repository"
	"effective_mobile/pkg/agify"
	"effective_mobile/pkg/genderize"
	"effective_mobile/pkg/nationalize"
	"effective_mobile/pkg/postgres"
	"effective_mobile/status"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type uc struct {
	psql              postgres.Postgres
	repo              person.Repo
	cfg               *config.Config
	logger            *zap.Logger
	agifyClient       *agify.Client
	genderizeClient   *genderize.Client
	nationalityClient *nationalize.Client
}

func New(psql postgres.Postgres, cfg *config.Config, logger *zap.Logger) person.UC {
	return &uc{
		psql:   psql,
		repo:   repository.NewPostgres(psql),
		cfg:    cfg,
		logger: logger,
		agifyClient: agify.NewClient(agify.Config{
			Logger: logger,
			APIurl: cfg.AgifyURL,
		}),
		genderizeClient: genderize.NewClient(genderize.Config{
			Logger: logger,
			APIurl: cfg.GenderizeURL,
		}),
		nationalityClient: nationalize.NewClient(nationalize.Config{
			Logger: logger,
			APIurl: cfg.NationalizeURL,
		}),
	}
}

func (u *uc) CreatePerson(ctx context.Context, params model.CreatePersonParams) error {
	if params.Name == nil {
		return status.ErrMandatoryParams.Wrap(fmt.Errorf("name is empty"))
	}
	if params.Surname == nil {
		return status.ErrMandatoryParams.Wrap(fmt.Errorf("surname is empty"))
	}
	if params.Patronymic == nil {
		return status.ErrMandatoryParams.Wrap(fmt.Errorf("patronymic is empty"))
	}
	ctxlogger := u.logger.With(zap.String("name", *params.Name), zap.String("surname", *params.Surname), zap.String("patronymic", *params.Patronymic))

	corrId := uuid.NewString()
	ctxlogger.Info("create new person", zap.String("corr_id", corrId))

	ageResp, err := u.agifyClient.GetAge(ctx, corrId, *params.Name)
	if err != nil {
		ctxlogger.Error("error to get age from agify client", zap.Error(err))
		return status.ErrExternalService.Wrap(fmt.Errorf("age is undefined"))
	}
	genderResp, err := u.genderizeClient.GetGender(ctx, corrId, *params.Name)
	if err != nil {
		ctxlogger.Error("error to get gender from genderize client", zap.Error(err))
		return status.ErrExternalService.Wrap(fmt.Errorf("gender is undefined"))
	}
	nationalityResp, err := u.nationalityClient.GetNationality(ctx, corrId, *params.Name)
	if err != nil {
		ctxlogger.Error("error to get nationality from nationalize client", zap.Error(err))
		return status.ErrExternalService.Wrap(fmt.Errorf("nationality is undefined"))
	}

	_, err = u.repo.InsertPerson(ctx, model.Person{
		Name:        *params.Name,
		Surname:     *params.Surname,
		Patronymic:  params.Patronymic,
		Age:         &ageResp.Age,
		Gender:      (*model.PersonGender)(&genderResp.Gender),
		Nationality: &nationalityResp.Country[0].Country_id,
		Created:     time.Now().UTC(),
		Updated:     time.Now().UTC(),
		Removed:     false,
	})
	if err != nil {
		ctxlogger.Error("error to insert person to db", zap.Error(err))
		return status.ErrUnexpected.Wrap(fmt.Errorf("cannot create person"))
	}
	return nil
}

func (u *uc) UpdatePerson(ctx context.Context, params model.UpdatePersonParams) error {
	person, err := u.repo.GetPersonById(ctx, params.PersonId)
	if err != nil {
		u.logger.Error("cannot get person to update him cause db error", zap.Error(err))
		return status.ErrUnexpected.Wrap(fmt.Errorf("cannot update person"))
	}
	if person == nil {
		u.logger.Error("cannot update no person", zap.Error(fmt.Errorf("person not found")), zap.Int64("person_id", params.PersonId))
		return status.ErrNotFound.Wrap(fmt.Errorf("no person with id - %d", params.PersonId))
	}

	err = postgres.ExecTx(ctx, func(err error) {
		u.logger.Error("execute tx to update person throw error", zap.Error(err))
	}, u.psql, func(tx postgres.Tx) error {
		self := repository.NewPostgres(tx)
		if params.Age != nil {
			_, err := self.UpdatePersonAge(ctx, params.PersonId, *params.Age, time.Now().UTC())
			if err != nil {
				u.logger.Error("cannot update person age cause db error", zap.Error(err))
				return status.ErrUnexpected.Wrap(fmt.Errorf("cannot update person age"))
			}
		}
		if params.Gender != nil {
			_, err := self.UpdatePersonGender(ctx, params.PersonId, *params.Gender, time.Now().UTC())
			if err != nil {
				u.logger.Error("cannot update person gender cause db error", zap.Error(err))
				return status.ErrUnexpected.Wrap(fmt.Errorf("cannot update person gender"))
			}
		}
		if params.Nationality != nil {
			_, err := self.UpdatePersonNationality(ctx, params.PersonId, *params.Nationality, time.Now().UTC())
			if err != nil {
				u.logger.Error("cannot update person nationality cause db error", zap.Error(err))
				return status.ErrUnexpected.Wrap(fmt.Errorf("cannot update person nationality"))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *uc) DeletePerson(ctx context.Context, id int64) error {
	person, err := u.repo.GetPersonById(ctx, id)
	if err != nil {
		u.logger.Error("cannot get person to delete him cause db error", zap.Error(err))
		return status.ErrUnexpected.Wrap(fmt.Errorf("cannot delete person"))
	}
	if person == nil {
		u.logger.Error("cannot delete no person", zap.Error(fmt.Errorf("person not found")), zap.Int64("person_id", id))
		return status.ErrNotFound.Wrap(fmt.Errorf("no person with id - %d", id))
	}
	_, err = u.repo.SetPersonRemoved(ctx, id, time.Now().UTC())
	if err != nil {
		u.logger.Error("cannot delete person cause db error", zap.Error(err))
		return status.ErrUnexpected.Wrap(fmt.Errorf("cannot delete person"))
	}
	return nil
}

func (u *uc) FetchPersons(ctx context.Context, params model.FetchPersonsParams) ([]model.Person, error) {
	if params.Limit == 0 {
		return []model.Person{}, nil
	}
	res, err := u.repo.FetchPersons(ctx, params)
	if err != nil {
		u.logger.Error("cannot fetch persons cause db error", zap.Error(err))
		return nil, status.ErrUnexpected.Wrap(fmt.Errorf("cannot fetch persons"))
	}
	if res == nil {
		res = []model.Person{}
	}
	return res, nil
}
