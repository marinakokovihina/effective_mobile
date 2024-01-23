package status

const (
	FailedBstatus  = "FAILED"
	SuccessBstatus = "SUCCESS"
)

type HTTPresponse struct {
	Status      string      `json:"status"`
	Description *string     `json:"description,omitempty"`
	Result      interface{} `json:"result,omitempty"`
}
