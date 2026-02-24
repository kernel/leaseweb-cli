package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

type APIError struct {
	Method     string
	URL        string
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s %s: %s", e.Method, e.URL, e.Status)
}

type Client struct {
	http    *http.Client
	baseURL string
	apiKey  string
	debug   bool
}

func NewClient(cmd *cli.Command) (*Client, error) {
	apiKey, err := resolveAPIKey(cmd)
	if err != nil {
		return nil, err
	}
	return &Client{
		http:    http.DefaultClient,
		baseURL: resolveBaseURL(cmd),
		apiKey:  apiKey,
		debug:   cmd.Root().Bool("debug"),
	}, nil
}

// NewClientWithBaseURL creates a client with explicit base URL (useful for testing).
func NewClientWithBaseURL(baseURL, apiKey string) *Client {
	return &Client{
		http:    http.DefaultClient,
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (c *Client) Do(ctx context.Context, method, path string, body io.Reader) (gjson.Result, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("X-LSW-Auth", c.apiKey)
	req.Header.Set("User-Agent", fmt.Sprintf("lw-cli/%s", Version))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.debug {
		if dump, err := httputil.DumpRequest(req, true); err == nil {
			log.Printf("Request:\n%s\n", dump)
		}
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("reading response: %w", err)
	}

	if c.debug {
		if dump, err := httputil.DumpResponse(resp, false); err == nil {
			log.Printf("Response:\n%s\n%s\n", dump, string(respBody))
		}
	}

	if resp.StatusCode >= 400 {
		return gjson.Result{}, &APIError{
			Method:     method,
			URL:        url,
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       string(respBody),
		}
	}

	// For 204 No Content or empty bodies, return empty result
	if resp.StatusCode == 204 || len(respBody) == 0 {
		return gjson.Result{}, nil
	}

	return gjson.ParseBytes(respBody), nil
}

func (c *Client) Get(ctx context.Context, path string) (gjson.Result, error) {
	return c.Do(ctx, "GET", path, nil)
}

func (c *Client) Post(ctx context.Context, path string, jsonBody string) (gjson.Result, error) {
	var body io.Reader
	if jsonBody != "" {
		body = strings.NewReader(jsonBody)
	}
	return c.Do(ctx, "POST", path, body)
}

func (c *Client) Put(ctx context.Context, path string, jsonBody string) (gjson.Result, error) {
	var body io.Reader
	if jsonBody != "" {
		body = strings.NewReader(jsonBody)
	}
	return c.Do(ctx, "PUT", path, body)
}

func (c *Client) Delete(ctx context.Context, path string) (gjson.Result, error) {
	return c.Do(ctx, "DELETE", path, nil)
}

// DoRaw is like Do but returns the raw response body bytes (for binary responses like PDFs).
func (c *Client) DoRaw(ctx context.Context, method, path string) ([]byte, string, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("X-LSW-Auth", c.apiKey)
	req.Header.Set("User-Agent", fmt.Sprintf("lw-cli/%s", Version))

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, "", &APIError{
			Method:     method,
			URL:        url,
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       string(respBody),
		}
	}

	return respBody, resp.Header.Get("Content-Type"), nil
}

// BuildQueryString constructs a query string from key-value pairs, omitting empty values.
func BuildQueryString(params map[string]string) string {
	var parts []string
	for k, v := range params {
		if v != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v))
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return "?" + strings.Join(parts, "&")
}

// PostJSON sends a POST with a pre-built JSON body.
func (c *Client) PostJSON(ctx context.Context, path string, jsonBody []byte) (gjson.Result, error) {
	return c.Do(ctx, "POST", path, bytes.NewReader(jsonBody))
}

// PutJSON sends a PUT with a pre-built JSON body.
func (c *Client) PutJSON(ctx context.Context, path string, jsonBody []byte) (gjson.Result, error) {
	return c.Do(ctx, "PUT", path, bytes.NewReader(jsonBody))
}

// PatchJSON sends a PATCH with a pre-built JSON body.
func (c *Client) PatchJSON(ctx context.Context, path string, jsonBody []byte) (gjson.Result, error) {
	return c.Do(ctx, "PATCH", path, bytes.NewReader(jsonBody))
}

// DeleteWithBody sends a DELETE with a JSON body.
func (c *Client) DeleteWithBody(ctx context.Context, path string, jsonBody []byte) (gjson.Result, error) {
	return c.Do(ctx, "DELETE", path, bytes.NewReader(jsonBody))
}
