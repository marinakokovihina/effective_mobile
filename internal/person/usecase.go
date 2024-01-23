package person

import (
	"context"
	"effective_mobile/internal/person/model"
)

type UC interface {
	CreatePerson(ctx context.Context, params model.CreatePersonParams) error
	UpdatePerson(ctx context.Context, params model.UpdatePersonParams) error
	DeletePerson(ctx context.Context, id int64) error
	FetchPersons(ctx context.Context, params model.FetchPersonsParams) ([]model.Person, error)
}
