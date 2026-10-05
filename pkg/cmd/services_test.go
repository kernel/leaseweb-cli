package cmd

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServicesCancel(t *testing.T) {
	for _, test := range []struct {
		name    string
		reason  string
		detail  string
		payload map[string]string
	}{
		{
			name:   "predefined reason",
			reason: "CANCEL_PRODUCT_SWITCH",
			payload: map[string]string{
				"reasonCode": "CANCEL_PRODUCT_SWITCH",
			},
		},
		{
			name:   "other reason with detail",
			reason: "CANCEL_OTHER",
			detail: "Do not renew at the end of the current contract term.",
			payload: map[string]string{
				"reasonCode": "CANCEL_OTHER",
				"reason":     "Do not renew at the end of the current contract term.",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			called := false
			srv := newTestServer(t, map[string]http.HandlerFunc{
				"POST /services/v1/services/12345/cancel": func(w http.ResponseWriter, r *http.Request) {
					called = true
					var payload map[string]string
					require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
					assert.Equal(t, test.payload, payload)
					w.WriteHeader(http.StatusNoContent)
				},
			})
			defer srv.Close()

			args := []string{"services", "cancel", "12345", "--reason", test.reason}
			if test.detail != "" {
				args = append(args, "--reason-detail", test.detail)
			}
			_, stderr, err := runCLI(t, srv.URL, args)
			require.NoError(t, err)
			assert.True(t, called)
			assert.Contains(t, stderr, "Cancelled service 12345")
		})
	}
}
