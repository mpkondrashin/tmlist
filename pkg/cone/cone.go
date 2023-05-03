//////////////////////////////////////////////////////////////////////////
//
//  (c) TMList 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)
//  Copyright under MIT Lincese. Please see LICENSE file for details
//
//  c1ews.go - partial implementation of CloudOne Endpoint &
//  Workload Security API. For more details, refer to:
//  https://cloudone.trendmicro.com/docs/workload-security/api-reference/
//
//////////////////////////////////////////////////////////////////////////

package cone

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mpkondrashin/tmlist/pkg/c1ews"
)

var ErrWrongURNFormat = errors.New("worng URN format")

const Version = "v1"
const EntryPoint = "https://accounts.cloudone.trendmicro.com/api"

type Client struct {
	APIKey          string
	Host            string
	IgnoreTLSErrors bool
}

func NewClient(APIKey string) *Client {
	return &Client{
		APIKey:          APIKey,
		Host:            EntryPoint,
		IgnoreTLSErrors: false,
	}
}

func (c *Client) SetIgnoreTLSErrors(ignoreTLSErrors bool) *Client {
	c.IgnoreTLSErrors = ignoreTLSErrors
	return c
}

func (c *Client) NewWorkloadSecurity(Host string) *c1ews.Client {
	return c1ews.NewWorkloadSecurity(c.APIKey, Host)
}

type WSError struct {
	Message string `json:"message"`
}

type DescribeAPIKeyResponse struct {
	ID           string    `json:"id"`
	Alias        string    `json:"alias"`
	RoleID       string    `json:"roleID"`
	Locale       string    `json:"locale"`
	Timezone     string    `json:"timezone"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
	LastActivity time.Time `json:"lastActivity"`
	Enabled      bool      `json:"enabled"`
	URN          string    `json:"urn"`
}

func APIKeyID(APIKey string) string {
	return strings.Split(APIKey, ":")[0]
}

func (c *Client) DescribeAPIKey(ctx context.Context, APIKey string) (*DescribeAPIKeyResponse, error) {
	id := APIKeyID(APIKey)
	//url := fmt.Sprintf("/apikeys/%s", id)
	var response DescribeAPIKeyResponse
	err := c.query(ctx, "GET", "/apikeys/"+id, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) DescribeCurrentAPIKey(ctx context.Context) (*DescribeAPIKeyResponse, error) {
	return c.DescribeAPIKey(ctx, c.APIKey)
}

func (c *Client) APIKeyRegion(ctx context.Context, APIKey string) (string, error) {
	resp, err := c.DescribeAPIKey(ctx, APIKey)
	if err != nil {
		return "", err
	}
	urn := strings.Split(resp.URN, ":")
	if len(urn) < 3 {
		return "", fmt.Errorf("%s: %w", resp.URN, ErrWrongURNFormat)
	}
	return urn[3], nil
}

func (c *Client) CurrentAPIKeyRegion(ctx context.Context) (string, error) {
	return c.APIKeyRegion(ctx, c.APIKey)
}

func (c *Client) query(ctx context.Context,
	method string,
	url string,
	requestBody io.Reader,
	response any) error {
	uri := c.Host + url
	req, err := http.NewRequestWithContext(ctx, method, uri, requestBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "ApiKey "+c.APIKey)
	req.Header.Set("api-secret-key", c.APIKey)
	req.Header.Set("api-version", Version)
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.IgnoreTLSErrors}, //nolint
	}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request: %w", err)
	}
	if resp.StatusCode != 200 {
		var data bytes.Buffer
		if _, err := io.Copy(&data, resp.Body); err != nil {
			return fmt.Errorf("error body receive: %w", err)
		}
		var wse WSError
		if err := json.Unmarshal(data.Bytes(), &wse); err != nil {
			return fmt.Errorf("unmarshal error body: %w", err)
		}
		return fmt.Errorf("code %d: %s", resp.StatusCode, wse.Message)
	}
	defer resp.Body.Close()
	//io.Copy(os.Stdout, resp.Body)
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("response parse: %w", err)
	}
	return nil
}
