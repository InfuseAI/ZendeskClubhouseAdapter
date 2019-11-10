package cloudfunction

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestZendeskClubhouseAdapter(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := map[string]struct {
		method     string
		payload    string
		wantStatus int
	}{
		"create ticket": {http.MethodPost, "{}", http.StatusCreated},
	}
	os.Setenv("CH_TOKEN", "MOCK_CLUBHOUSE")
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			jsonPayload := bytes.NewBuffer([]byte(tt.payload))
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, "/", jsonPayload)

			ZendeskClubhouseAdapter(w, r)

			rw := w.Result()
			defer rw.Body.Close()

			if s := rw.StatusCode; s != tt.wantStatus {
				t.Fatalf("got: %d, want: %d", s, tt.wantStatus)
			}
		})
	}
}
