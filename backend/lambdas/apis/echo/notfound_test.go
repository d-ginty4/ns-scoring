package echo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUnknownRoute(t *testing.T) {
	e := NewEchoInstance()

	tests := []struct {
		method string
		path   string
		status int
		msg    string
	}{
		{method: http.MethodGet, path: "/does-not-exist", status: http.StatusNotFound, msg: "route not found"},
		{method: http.MethodPost, path: "/api/v1/missing", status: http.StatusNotFound, msg: "route not found"},
		{method: http.MethodDelete, path: "/", status: http.StatusNotFound, msg: "route not found"},
		{method: http.MethodPost, path: "/routes", status: http.StatusMethodNotAllowed, msg: "Method Not Allowed"},
		{method: http.MethodPatch, path: "/api/v1/regions", status: http.StatusMethodNotAllowed, msg: "Method Not Allowed"},
	}

	for _, tc := range tests {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != tc.status {
				t.Fatalf("status = %d, want %d, body = %s", rec.Code, tc.status, rec.Body.String())
			}

			var body struct {
				Data  any `json:"data"`
				Error *struct {
					Message string `json:"message"`
				} `json:"error"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode body: %v, body = %s", err, rec.Body.String())
			}
			if body.Error == nil || body.Error.Message != tc.msg {
				t.Fatalf("error message = %#v, want %q", body.Error, tc.msg)
			}
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/routes", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /routes status = %d, body = %s", rec.Code, rec.Body.String())
	}
}
