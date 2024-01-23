package genderize

type GenderResponse struct {
	Count       *int    `json:"count,omitempty"`
	Name        *string `json:"name,omitempty"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}
