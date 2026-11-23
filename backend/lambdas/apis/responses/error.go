package responses

import (
	"errors"
	"net/http"

	"apis/exceptions"
	"apis/logger"

	"github.com/labstack/echo/v5"
)

func HTTPErrorHandler(ctx *echo.Context, err error) {
	if r, _ := echo.UnwrapResponse(ctx.Response()); r != nil && r.Committed {
		return
	}

	status, apiErr := mapError(ctx, err)
	id := RequestID(ctx)
	if id != "" {
		ctx.Response().Header().Set(echo.HeaderXRequestID, id)
	}

	write(ctx, status, Envelope[any]{
		RequestID: id,
		Data:      nil,
		Error:     &apiErr,
	})
}

func mapError(c *echo.Context, err error) (int, APIError) {
	var verr *exceptions.ErrValidation
	if errors.As(err, &verr) {
		return http.StatusBadRequest, APIError{
			Message: "validation failed",
			Details: verr.Errors,
		}
	}

	var nferr *exceptions.ErrNotFound
	if errors.As(err, &nferr) {
		return http.StatusNotFound, APIError{Message: nferr.Error()}
	}

	var uniqueField *exceptions.ErrBadEntity
	if errors.As(err, &uniqueField) {
		return http.StatusBadRequest, APIError{Message: uniqueField.Error()}
	}

	var he *echo.HTTPError
	if errors.As(err, &he) {
		msg := httpErrorMessage(he)
		if he.Code >= 500 {
			logger.From(c).Error("request failed", logger.Err(err))
		}
		return he.Code, APIError{Message: msg}
	}

	// Echo's built-in sentinels (ErrNotFound, ErrMethodNotAllowed) implement
	// HTTPStatusCoder but are not *HTTPError.
	var statusCoder echo.HTTPStatusCoder
	if errors.As(err, &statusCoder) {
		code := statusCoder.StatusCode()
		if code != 0 {
			if code >= 500 {
				logger.From(c).Error("request failed", logger.Err(err))
			}
			msg := http.StatusText(code)
			if code == http.StatusNotFound {
				msg = "route not found"
			}
			return code, APIError{Message: msg}
		}
	}

	logger.From(c).Error("request failed", logger.Err(err))
	return http.StatusInternalServerError, APIError{Message: "internal server error"}
}

func httpErrorMessage(httpError *echo.HTTPError) string {
	if httpError.Message != "" {
		return httpError.Message
	}
	return http.StatusText(httpError.Code)
}
