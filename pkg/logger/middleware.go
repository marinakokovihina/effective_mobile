package logger

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type CorrId struct{}

func NewResponseMDW(title string, logger *zap.Logger) resty.ResponseMiddleware {
	return func(_ *resty.Client, r *resty.Response) error {
		var msg json.RawMessage
		if resty.IsJSONType(r.Header().Get("Content-Type")) {
			msg = r.Body()
		}
		if len(msg) == 0 {
			msg = []byte("no response body")
		}
		logger.Info(title,
			zap.Any("corr_id", r.Request.Context().Value(CorrId{})),
			zap.String("http_status", r.Status()),
			zap.String("http_scheme", r.RawResponse.Proto),
			zap.String("http_response", string(msg)),
			zap.Duration("http_elapsed", r.Time()),
			zap.String("http_url", r.Request.URL))
		return nil
	}
}
