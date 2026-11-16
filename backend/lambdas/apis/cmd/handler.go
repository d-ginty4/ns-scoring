package cmd

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	e "apis/echo"
	"apis/logger"
	"apis/repository"
	"apis/responses"

	"github.com/aws/aws-lambda-go/events"
	"github.com/labstack/echo/v5"
)

var echoInstance *echo.Echo = e.NewEchoInstance()

func Entry(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	requestID := req.RequestContext.RequestID
	slog.Debug("lambda invocation",
		"requestID", requestID,
		"method", req.RequestContext.HTTP.Method,
		"path", req.RawPath,
	)

	httpReq, err := GatewayRequestToHttpRequest(ctx, req)
	if err != nil {
		slog.Error("failed to convert API Gateway request", "requestID", requestID, logger.Err(err))
		body, marshalErr := json.Marshal(responses.Envelope[any]{
			RequestID: requestID,
			Data:      nil,
			Error:     &responses.APIError{Message: "unable to handle request"},
		})
		if marshalErr != nil {
			body = []byte(`{"requestID":"","data":null,"error":{"message":"internal server error"}}`)
		}

		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(body),
		}, err
	}

	repository.ConnectDB()
	rec := new(ResponseWriter)
	echoInstance.ServeHTTP(rec, httpReq)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: rec.StatusCode,
		Body:       rec.Body.String(),
		Headers:    req.Headers,
	}, nil
}
