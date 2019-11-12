package cloudfunction

import (
	"bytes"
	"encoding/base64"
	"fmt"
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
		method         string
		clubhouseToken string
		user           string
		password       string
		payload        string
		wantStatus     int
	}{
		"create ticket":                         {http.MethodPost, "MOCK_CLUBHOUSE", "", "", `{"title": "unit test", "id": 7777, "url": "http://unittest.io" }`, http.StatusCreated},
		"create ticket with auth":               {http.MethodPost, "MOCK_CLUBHOUSE", "unit-test", "YouShallNotPass!", `{"title": "unit test", "id": 7777, "url": "http://unittest.io" }`, http.StatusCreated},
		"create ticket with invalid payload":    {http.MethodPost, "MOCK_CLUBHOUSE", "unit-test", "YouShallNotPass!", `{}`, http.StatusBadRequest},
		"create ticket without clubhouse token": {http.MethodPost, "", "", "", `{"title": "unit test", "id": 7777, "url": "http://unittest.io" }`, http.StatusBadRequest},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			jsonPayload := bytes.NewBuffer([]byte(tt.payload))
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, "/", jsonPayload)
			os.Setenv("CH_TOKEN", tt.clubhouseToken)
			os.Setenv("AUTH_USER", tt.user)
			os.Setenv("AUTH_PASSWORD", tt.password)

			// Prepare BasicAuth payload
			if tt.user != "" || tt.password != "" {
				basicAuthPayload := base64.StdEncoding.EncodeToString([]byte(
					fmt.Sprintf("%s:%s", tt.user, tt.password),
				))
				r.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuthPayload))
			}

			ZendeskClubhouseAdapter(w, r)

			rw := w.Result()
			defer rw.Body.Close()

			if s := rw.StatusCode; s != tt.wantStatus {
				t.Fatalf("got: %d, want: %d", s, tt.wantStatus)
			}
		})
	}
}
