package agify

type AgeResponse struct {
	Count *int    `json:"count,omitempty"`
	Name  *string `json:"name,omitempty"`
	Age   int     `json:"age"`
}
