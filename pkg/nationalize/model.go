package nationalize

type NationalityResponse struct {
	Count   *int      `json:"count,omitempty"`
	Name    *string   `json:"name,omitempty"`
	Country []Country `json:"country"`
}

type Country struct {
	Country_id  string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
