package person

import (
	"context"
	"effective_mobile/internal/person/model"
	"time"
)

type Repo interface {
	InsertPerson(ctx context.Context, p model.Person) (*int64, error)
	GetPersonById(ctx context.Context, id int64) (*model.Person, error)
	SetPersonRemoved(ctx context.Context, id int64, at time.Time) (*int64, error)
	UpdatePersonAge(ctx context.Context, id int64, age int, at time.Time) (*int64, error)
	UpdatePersonGender(ctx context.Context, id int64, gender model.PersonGender, at time.Time) (*int64, error)
	UpdatePersonNationality(ctx context.Context, id int64, nationality string, at time.Time) (*int64, error)
	FetchPersons(ctx context.Context, filters model.FetchPersonsParams) ([]model.Person, error)
}
