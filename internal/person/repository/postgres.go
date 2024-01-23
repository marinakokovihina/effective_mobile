package repository

import (
	"context"
	"effective_mobile/internal/person"
	"effective_mobile/internal/person/model"
	db "effective_mobile/pkg/postgres"
	"time"
)

type postgres struct {
	psql db.Postgres
}

func NewPostgres(psql db.Postgres) person.Repo {
	return &postgres{
		psql: psql,
	}
}

func (p postgres) InsertPerson(ctx context.Context, person model.Person) (*int64, error) {
	var affected int64
	err := p.psql.QueryRow(ctx, `insert into public.person (name, surname, patronymic, age, gender, nationality, created, updated, removed) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`,
		person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality, person.Created, person.Updated, person.Removed).Scan(&affected)
	if err != nil {
		return nil, err
	}
	return &affected, nil
}

func (p postgres) GetPersonById(ctx context.Context, id int64) (*model.Person, error) {
	var data []model.Person
	err := p.psql.Select(ctx, &data, `select id, name, surname, patronymic, age, gender, nationality, created, updated, removed from public.person where id=$1 and removed=false limit 1`, id)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	return &data[0], nil
}

func (p postgres) SetPersonRemoved(ctx context.Context, id int64, at time.Time) (*int64, error) {
	var affected int64
	err := p.psql.QueryRow(ctx, `update public.person set removed=true, updated=$1 where id=$2 returning id`, at, id).Scan(&affected)
	if err != nil {
		return nil, err
	}
	return &affected, nil
}

func (p postgres) UpdatePersonAge(ctx context.Context, id int64, age int, at time.Time) (*int64, error) {
	var affected int64
	err := p.psql.QueryRow(ctx, `update public.person set age=$1, updated=$2 where id=$3 returning id`, age, at, id).Scan(&affected)
	if err != nil {
		return nil, err
	}
	return &affected, nil
}

func (p postgres) UpdatePersonGender(ctx context.Context, id int64, gender model.PersonGender, at time.Time) (*int64, error) {
	var affected int64
	err := p.psql.QueryRow(ctx, `update public.person set gender=$1, updated=$2 where id=$3 returning id`, gender, at, id).Scan(&affected)
	if err != nil {
		return nil, err
	}
	return &affected, nil
}

func (p postgres) UpdatePersonNationality(ctx context.Context, id int64, nationality string, at time.Time) (*int64, error) {
	var affected int64
	err := p.psql.QueryRow(ctx, `update public.person set nationality=$1, updated=$2 where id=$3 returning id`, nationality, at, id).Scan(&affected)
	if err != nil {
		return nil, err
	}
	return &affected, nil
}

func (p postgres) FetchPersons(ctx context.Context, ffs model.FetchPersonsParams) ([]model.Person, error) {
	var data []model.Person
	err := p.psql.Select(ctx, &data, `select id, name, surname, patronymic, age, gender, nationality, created, updated, removed from public.person where removed=false
		and (id=$1::bigint or $1 is null) and (name=$2::text or $2 is null) and (surname=$3::text or $3 is null) and (patronymic=$4::text or $4 is null) and (age=$5::int or $5 is null) and (gender=$6::gender or $6 is null)
		 and (nationality=$7::text or $7 is null) order by id desc limit $8 offset $9`,
		ffs.Id, ffs.Name, ffs.Surname, ffs.Patronymic, ffs.Age, ffs.Gender, ffs.Nationality, ffs.Limit, ffs.Offset)
	if err != nil {
		return nil, err
	}
	return data, nil
}
