package model

type CreatePersonParams struct {
	Name       *string `json:"name"`
	Surname    *string `json:"surname"`
	Patronymic *string `json:"patronymic"`
}

type UpdatePersonParams struct {
	PersonId    int64         `json:"person_id"`
	Age         *int          `json:"age,omitempty"`
	Gender      *PersonGender `json:"gender,omitempty"`
	Nationality *string       `json:"nationality,omitempty"`
}

type FetchPersonsParams struct {
	Id          *int64        `query:"id,omitempty"`
	Name        *string       `query:"name,omitempty"`
	Surname     *string       `query:"surname,omitempty"`
	Patronymic  *string       `query:"patronymic,omitempty"`
	Age         *int          `query:"age,omitempty"`
	Gender      *PersonGender `query:"gender,omitempty"`
	Nationality *string       `query:"nationality,omitempty"`
	Limit       int64         `query:"limit"`
	Offset      int64         `query:"offset"`
}
