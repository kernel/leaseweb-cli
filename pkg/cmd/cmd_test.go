package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, handlers map[string]http.HandlerFunc) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	for pattern, handler := range handlers {
		mux.HandleFunc(pattern, handler)
	}
	return httptest.NewServer(mux)
}

func jsonResponse(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(body)
}

func runCLI(t *testing.T, args []string) (stdout string, stderr string, err error) {
	t.Helper()

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	fullArgs := append([]string{"lw"}, args...)
	err = Command.Run(context.Background(), fullArgs)

	wOut.Close()
	wErr.Close()

	var outBuf, errBuf bytes.Buffer
	_, _ = outBuf.ReadFrom(rOut)
	_, _ = errBuf.ReadFrom(rErr)

	return outBuf.String(), errBuf.String(), err
}

func TestDedicatedServersList(t *testing.T) {
	srv := newTestServer(t, map[string]http.HandlerFunc{
		"GET /bareMetals/v2/servers": func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "test-key", r.Header.Get("X-LSW-Auth"))

			jsonResponse(w, 200, map[string]any{
				"servers": []map[string]any{
					{
						"id":        "12345",
						"reference": "my-server",
						"location":  map[string]any{"site": "WDC-02"},
						"specs": map[string]any{
							"chassis": "HP DL385 G11",
							"cpu":     map[string]any{"type": "AMD EPYC 9334"},
							"ram":     map[string]any{"size": 384, "unit": "GB"},
						},
						"networkInterfaces": map[string]any{
							"public": map[string]any{"ip": "1.2.3.4"},
						},
					},
				},
				"_metadata": map[string]any{"totalCount": 1, "limit": 20, "offset": 0},
			})
		},
	})
	defer srv.Close()

	stdout, _, err := runCLI(t, []string{
		"--api-key", "test-key",
		"--base-url", srv.URL,
		"dedicated-servers", "list",
	})
	require.NoError(t, err)
	assert.Contains(t, stdout, "12345")
	assert.Contains(t, stdout, "my-server")
	assert.Contains(t, stdout, "WDC-02")
	assert.Contains(t, stdout, "AMD EPYC 9334")
}

func TestDedicatedServersGet(t *testing.T) {
	srv := newTestServer(t, map[string]http.HandlerFunc{
		"GET /bareMetals/v2/servers/12345": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, map[string]any{
				"id":        "12345",
				"reference": "my-server",
				"location":  map[string]any{"site": "WDC-02"},
			})
		},
	})
	defer srv.Close()

	stdout, _, err := runCLI(t, []string{
		"--api-key", "test-key",
		"--base-url", srv.URL,
		"--format", "raw",
		"dedicated-servers", "get", "12345",
	})
	require.NoError(t, err)
	assert.Contains(t, stdout, `"id":"12345"`)
}

func TestDedicatedServersGetMissingID(t *testing.T) {
	_, _, err := runCLI(t, []string{
		"--api-key", "test-key",
		"--base-url", "http://localhost",
		"dedicated-servers", "get",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "server ID required")
}

func TestIPsList(t *testing.T) {
	srv := newTestServer(t, map[string]http.HandlerFunc{
		"GET /ipMgmt/v2/ips": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, map[string]any{
				"ips": []map[string]any{
					{
						"ip":            "1.2.3.4",
						"version":       4,
						"type":          "PRIMARY",
						"reverseLookup": "host.example.com",
						"nullRouted":    false,
						"equipmentId":   "12345",
					},
				},
			})
		},
	})
	defer srv.Close()

	stdout, _, err := runCLI(t, []string{
		"--api-key", "test-key",
		"--base-url", srv.URL,
		"ips", "list",
	})
	require.NoError(t, err)
	assert.Contains(t, stdout, "1.2.3.4")
	assert.Contains(t, stdout, "PRIMARY")
	assert.Contains(t, stdout, "host.example.com")
}

func TestInvoicesList(t *testing.T) {
	srv := newTestServer(t, map[string]http.HandlerFunc{
		"GET /invoices/v1/invoices": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, map[string]any{
				"invoices": []map[string]any{
					{
						"id":       "INV001",
						"date":     "2026-01-15",
						"status":   "PAID",
						"total":    1234.56,
						"currency": "USD",
						"dueDate":  "2026-02-15",
					},
				},
			})
		},
	})
	defer srv.Close()

	stdout, _, err := runCLI(t, []string{
		"--api-key", "test-key",
		"--base-url", srv.URL,
		"invoices", "list",
	})
	require.NoError(t, err)
	assert.Contains(t, stdout, "INV001")
	assert.Contains(t, stdout, "PAID")
	assert.Contains(t, stdout, "1234.56")
}

func TestServicesList(t *testing.T) {
	srv := newTestServer(t, map[string]http.HandlerFunc{
		"GET /services/v1/services": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, map[string]any{
				"services": []map[string]any{
					{
						"id":        "SVC001",
						"reference": "prod-server",
						"productId": "DEDICATED_SERVER",
						"status":    "ACTIVE",
						"startDate": "2025-01-01",
					},
				},
			})
		},
	})
	defer srv.Close()

	stdout, _, err := runCLI(t, []string{
		"--api-key", "test-key",
		"--base-url", srv.URL,
		"services", "list",
	})
	require.NoError(t, err)
	assert.Contains(t, stdout, "SVC001")
	assert.Contains(t, stdout, "ACTIVE")
}

func TestAPIErrorHandling(t *testing.T) {
	srv := newTestServer(t, map[string]http.HandlerFunc{
		"GET /bareMetals/v2/servers/bad-id": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 404, map[string]any{
				"errorCode":    "404",
				"errorMessage": "The requested resource was not found.",
			})
		},
	})
	defer srv.Close()

	_, _, err := runCLI(t, []string{
		"--api-key", "test-key",
		"--base-url", srv.URL,
		"dedicated-servers", "get", "bad-id",
	})
	require.Error(t, err)

	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, 404, apiErr.StatusCode)
	assert.Contains(t, apiErr.Body, "errorMessage")
}

func TestJSONFormat(t *testing.T) {
	srv := newTestServer(t, map[string]http.HandlerFunc{
		"GET /bareMetals/v2/servers/12345": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, map[string]any{
				"id":        "12345",
				"reference": "test",
			})
		},
	})
	defer srv.Close()

	stdout, _, err := runCLI(t, []string{
		"--api-key", "test-key",
		"--base-url", srv.URL,
		"--format", "json",
		"dedicated-servers", "get", "12345",
	})
	require.NoError(t, err)
	assert.Contains(t, stdout, `"id"`)
	assert.Contains(t, stdout, `"12345"`)
}

func TestYAMLFormat(t *testing.T) {
	srv := newTestServer(t, map[string]http.HandlerFunc{
		"GET /bareMetals/v2/servers/12345": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, map[string]any{
				"id":        "12345",
				"reference": "test",
			})
		},
	})
	defer srv.Close()

	stdout, _, err := runCLI(t, []string{
		"--api-key", "test-key",
		"--base-url", srv.URL,
		"--format", "yaml",
		"dedicated-servers", "get", "12345",
	})
	require.NoError(t, err)
	assert.Contains(t, stdout, "id:")
	assert.Contains(t, stdout, "reference:")
}

func TestTransform(t *testing.T) {
	srv := newTestServer(t, map[string]http.HandlerFunc{
		"GET /bareMetals/v2/servers/12345": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, map[string]any{
				"id":        "12345",
				"reference": "my-server",
				"location": map[string]any{
					"site": "WDC-02",
					"unit": "001",
				},
			})
		},
	})
	defer srv.Close()

	stdout, _, err := runCLI(t, []string{
		"--api-key", "test-key",
		"--base-url", srv.URL,
		"--format", "raw",
		"--transform", "location.site",
		"dedicated-servers", "get", "12345",
	})
	require.NoError(t, err)
	assert.Equal(t, "\"WDC-02\"\n", stdout)
}

func TestEmptyListOutput(t *testing.T) {
	srv := newTestServer(t, map[string]http.HandlerFunc{
		"GET /bareMetals/v2/servers": func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, 200, map[string]any{
				"servers":   []any{},
				"_metadata": map[string]any{"totalCount": 0},
			})
		},
	})
	defer srv.Close()

	_, stderr, err := runCLI(t, []string{
		"--api-key", "test-key",
		"--base-url", srv.URL,
		"dedicated-servers", "list",
	})
	require.NoError(t, err)
	assert.Contains(t, stderr, "No dedicated servers found.")
}

func TestPagination(t *testing.T) {
	srv := newTestServer(t, map[string]http.HandlerFunc{
		"GET /bareMetals/v2/servers": func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "5", r.URL.Query().Get("limit"))
			assert.Equal(t, "10", r.URL.Query().Get("offset"))
			jsonResponse(w, 200, map[string]any{
				"servers": []any{},
			})
		},
	})
	defer srv.Close()

	_, _, err := runCLI(t, []string{
		"--api-key", "test-key",
		"--base-url", srv.URL,
		"dedicated-servers", "list", "--limit", "5", "--offset", "10",
	})
	require.NoError(t, err)
}

func TestTableWriter(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTableWriter(&buf, "NAME", "STATUS", "IP")
	tw.AddRow("server1", "active", "1.2.3.4")
	tw.AddRow("server2", "stopped", "5.6.7.8")
	tw.Render()

	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	require.Len(t, lines, 3)
	assert.Contains(t, lines[0], "NAME")
	assert.Contains(t, lines[0], "STATUS")
	assert.Contains(t, lines[0], "IP")
	assert.Contains(t, lines[1], "server1")
	assert.Contains(t, lines[2], "server2")
}

func TestFormatTimeAgo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{"zero", "", "N/A"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTimeAgo(parseTime(tt.input))
			assert.Contains(t, result, tt.contains)
		})
	}
}

func TestMaskKey(t *testing.T) {
	assert.Equal(t, "74B1****************************1948", maskKey("74B196B1-1346-43F8-A7F3-B29971231948"))
	assert.Equal(t, "****", maskKey("abcd"))
}

func TestBuildQueryString(t *testing.T) {
	q := BuildQueryString(map[string]string{
		"foo": "bar",
		"baz": "",
	})
	assert.Contains(t, q, "foo=bar")
	assert.NotContains(t, q, "baz")
}

func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func TestConfigCommands(t *testing.T) {
	stdout, _, err := runCLI(t, []string{"config", "show"})
	require.NoError(t, err)
	_ = stdout
}
