package cmd

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type ResponseWriter struct {
	Headers    http.Header
	Body       bytes.Buffer
	StatusCode int
}

func (w *ResponseWriter) Header() http.Header {
	if w.Headers == nil {
		return http.Header{}
	}
	return w.Headers
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	return w.Body.Write(b)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func GatewayRequestToHttpRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (*http.Request, error) {
	u := url.URL{
		Path:     req.RawPath,
		RawQuery: req.RawQueryString,
	}

	var body *bytes.Reader
	if req.Body != "" {
		if req.IsBase64Encoded {
			decoded, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(req.Body)))
			if err != nil {
				return nil, err
			}
			body = bytes.NewReader(decoded)
		} else {
			body = bytes.NewReader([]byte(req.Body))
		}
	} else {
		body = bytes.NewReader([]byte{})
	}

	requestID := req.RequestContext.RequestID
	ctx = context.WithValue(ctx, "requestID", requestID)

	method := req.RequestContext.HTTP.Method
	httpReq, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	if requestID != "" {
		httpReq.Header.Set("X-Request-ID", requestID)
	}

	return httpReq, nil
}
