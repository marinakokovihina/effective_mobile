package agify

import (
	"context"
	"fmt"
	"net/url"

	"effective_mobile/pkg/logger"

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
		return fmt.Errorf("agify client: logger is required")
	}
	if c.APIurl == "" {
		return fmt.Errorf("agify client: api url must be non-empty")
	}
	_, err := url.ParseRequestURI(c.APIurl)
	if err != nil {
		return fmt.Errorf("agify: validating api url: %v", err)
	}
	return nil
}

func NewClient(cfg Config) *Client {
	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	rc := resty.New()
	rc.SetBaseURL(cfg.APIurl)
	rc.OnAfterResponse(logger.NewResponseMDW("agify outgoing HTTP request finished", cfg.Logger))
	return &Client{
		logger: cfg.Logger,
		rc:     rc,
	}
}

func (c Client) GetAge(ctx context.Context, corrId string, name string) (*AgeResponse, error) {
	c.logger.Info("agify client: get age", zap.String("corr_id", corrId), zap.String("name", name))
	resp, err := c.rc.R().
		SetContext(context.WithValue(ctx, logger.CorrId{}, corrId)).
		SetQueryParam("name", name).
		SetResult(AgeResponse{}).
		Get("/")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("agify api: [%d]; [corr_id: %s, resp: %s]", resp.StatusCode(), corrId, resp.Body())
	}
	return resp.Result().(*AgeResponse), nil
}
