package c1ews

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const Version = "v1"

type Client struct {
	APIKey          string
	Host            string
	IgnoreTLSErrors bool
}

func NewWorkloadSecurity(APIKey string, Host string) *Client {
	return &Client{
		APIKey: APIKey,
		Host:   Host,
	}
}

func (c *Client) SetIgnoreTLSErrors(ignoreTLSErrors bool) *Client {
	c.IgnoreTLSErrors = ignoreTLSErrors
	return c
}

type WSError struct {
	Message string `json:"message"`
}

type List struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Items       []string `json:"items,omitempty"`
}

type ListResponse struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Items       []string `json:"items"`
	ID          int      `json:"ID"`
}

func (c *Client) ListDirectoryLists(ctx context.Context) ([]ListResponse, error) {
	return c.listLists(ctx, "/directorylists")
}

func (c *Client) ListFileExtensionLists(ctx context.Context) ([]ListResponse, error) {
	return c.listLists(ctx, "/fileextensionlists")
}

func (c *Client) ListFileLists(ctx context.Context) ([]ListResponse, error) {
	return c.listLists(ctx, "/filelists")
}

func (c *Client) ListIPLists(ctx context.Context) ([]ListResponse, error) {
	return c.listLists(ctx, "/iplists")
}

func (c *Client) ListMACLists(ctx context.Context) ([]ListResponse, error) {
	return c.listLists(ctx, "/maclists")
}

func (c *Client) ListPortLists(ctx context.Context) ([]ListResponse, error) {
	return c.listLists(ctx, "/portlists")
}

func (c *Client) ModifyDirectoryList(ctx context.Context, id int, dirList *List) (*ListResponse, error) {
	return c.modifyList(ctx, "directorylists", id, dirList)
}

func (c *Client) ModifyFileExtensionList(ctx context.Context, id int, dirList *List) (*ListResponse, error) {
	return c.modifyList(ctx, "fileextensionlists", id, dirList)
}

func (c *Client) ModifyFileList(ctx context.Context, id int, dirList *List) (*ListResponse, error) {
	return c.modifyList(ctx, "filelists", id, dirList)
}

func (c *Client) ModifyIPList(ctx context.Context, id int, dirList *List) (*ListResponse, error) {
	return c.modifyList(ctx, "iplists", id, dirList)
}

func (c *Client) ModifyMACList(ctx context.Context, id int, dirList *List) (*ListResponse, error) {
	return c.modifyList(ctx, "maclists", id, dirList)
}

func (c *Client) modifyList(ctx context.Context, path string, id int, dirList *List) (*ListResponse, error) {
	url := fmt.Sprintf("/%s/%d", path, id)
	body, err := json.Marshal(dirList)
	if err != nil {
		return nil, err
	}
	var response ListResponse
	err = c.query(ctx, "POST", url, bytes.NewBuffer(body), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) listLists(ctx context.Context, url string) ([]ListResponse, error) {
	var response map[string][]ListResponse
	err := c.query(ctx, "GET", url, nil, &response)
	if err != nil {
		return nil, err
	}
	for _, r := range response {
		return r, nil
	}
	return nil, fmt.Errorf("missing response data for %s", url)
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
