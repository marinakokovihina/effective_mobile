package status

import (
	"fmt"

	"github.com/aws/smithy-go/ptr"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	if bErr, ok := err.(Berror); ok {
		return ctx.Status(bErr.Code).JSON(HTTPresponse{
			Status:      FailedBstatus,
			Description: &bErr.Message,
		})
	}
	return ctx.Status(500).JSON(HTTPresponse{
		Status:      FailedBstatus,
		Description: ptr.String("unknown error"),
	})
}

type Berror struct {
	inner   error
	Message string
	Code    int
}

func (err Berror) Error() string {
	return err.Message
}

func (err Berror) Wrap(inner error) Berror {
	return Berror{
		Code:    err.Code,
		Message: fmt.Sprintf("%v: %v", err, inner),
		inner:   inner,
	}
}

func (err Berror) Unwrap() error {
	return err.inner
}

func (err Berror) Is(target error) bool {
	return target.Error() == err.Message
}

var (
	ErrUnexpected      = Berror{Message: "unexpected error", Code: 500}
	ErrInput           = Berror{Message: "input error", Code: 400}
	ErrMandatoryParams = Berror{Message: "mandatory params must be filled", Code: 400}
	ErrNoPrivilege     = Berror{Message: "no privilege", Code: 401}
	ErrNotFound        = Berror{Message: "not found", Code: 404}
	ErrNotImplemented  = Berror{Message: "not implemented", Code: 500}
	ErrExternalService = Berror{Message: "cannot get external dependency", Code: 422}
)
