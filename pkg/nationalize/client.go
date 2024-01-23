package nationalize

import (
	"context"
	"effective_mobile/pkg/logger"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Client struct {
	logger *zap.Logger
	rc     *resty.Client
}

type Config struct {
	Logger *zap.Logger
	APIurl string
}

func (c Config) Validate() error {
	if c.Logger == nil {
		return fmt.Errorf("nationalize client: logger is required")
	}
	if c.APIurl == "" {
		return fmt.Errorf("nationalize client: api url must be non-empty")
	}
	_, err := url.ParseRequestURI(c.APIurl)
	if err != nil {
		return fmt.Errorf("nationalize: validating api url: %v", err)
	}
	return nil
}

func NewClient(cfg Config) *Client {
	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	rc := resty.New()
	rc.SetBaseURL(cfg.APIurl)
	rc.OnAfterResponse(logger.NewResponseMDW("nationalize outgoing HTTP request finished", cfg.Logger))
	return &Client{
		logger: cfg.Logger,
		rc:     rc,
	}
}

func (c Client) GetNationality(ctx context.Context, corrId string, name string) (*NationalityResponse, error) {
	c.logger.Info("nationalize client: get nationality", zap.String("corr_id", corrId), zap.String("name", name))
	resp, err := c.rc.R().
		SetContext(context.WithValue(ctx, logger.CorrId{}, corrId)).
		SetQueryParam("name", name).
		SetResult(NationalityResponse{}).
		Get("/")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("nationalize api: [%d]; [corr_id: %s, resp: %s]", resp.StatusCode(), corrId, resp.Body())
	}
	return resp.Result().(*NationalityResponse), nil
}
