package healthchecks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IdeaEvolver/cutter-pkg/client"
)

type ServiceResponse struct {
	Status string `json:"status"`
}

type Client struct {
	Client *client.Client

	Platform    string
	Fulfillment string
	Crm         string
	Study       string
}

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// func New(c HttpClient, platform, fulfillment, crm, study string) *Client {
// 	return &Client{
// 		Client:      client.New(c),
// 		Platform:    platform,
// 		Fulfillment: fulfillment,
// 		Crm:         crm,
// 		Study:       study,
// 	}
// }

func (c *Client) do(ctx context.Context, req *client.Request, ret interface{}) error {
	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if ret != nil {
		return json.NewDecoder(res.Body).Decode(&ret)
	}

	return nil
}

func (c *Client) PlatformStatus(ctx context.Context) (*ServiceResponse, error) {
	fmt.Println("platform ", c.Platform)
	url := fmt.Sprintf("%s/healthcheck", c.Platform)
	fmt.Println("URL ", url)
	req, _ := client.NewRequestWithContext(ctx, "GET", url, nil)

	status := &ServiceResponse{}
	if err := c.do(ctx, req, &status); err != nil {
		return nil, err
	}

	return status, nil
}

func (c *Client) PlatformUIStatus(ctx context.Context) (*ServiceResponse, error) {
	url := "https://dev.cutter.live/sign-in"
	req, _ := client.NewRequestWithContext(ctx, "GET", url, nil)

	status := &ServiceResponse{}
	if err := c.do(ctx, req, &status); err != nil {
		return nil, err
	}

	return status, nil
}

func (c *Client) FulfillmentStatus(ctx context.Context) (*ServiceResponse, error) {
	url := fmt.Sprintf("%s/healthcheck", c.Fulfillment)
	req, _ := client.NewRequestWithContext(ctx, "GET", url, nil)

	status := &ServiceResponse{}
	if err := c.do(ctx, req, &status); err != nil {
		return nil, err
	}

	return status, nil
}

func (c *Client) CrmStatus(ctx context.Context) (*ServiceResponse, error) {
	url := fmt.Sprintf("%s/healthcheck", c.Crm)
	req, _ := client.NewRequestWithContext(ctx, "GET", url, nil)

	status := &ServiceResponse{}
	if err := c.do(ctx, req, &status); err != nil {
		return nil, err
	}

	return status, nil
}

func (c *Client) StudyStatus(ctx context.Context) (*ServiceResponse, error) {
	url := fmt.Sprintf("%s/healthcheck", c.Study)
	req, _ := client.NewRequestWithContext(ctx, "GET", url, nil)

	status := &ServiceResponse{}
	if err := c.do(ctx, req, &status); err != nil {
		return nil, err
	}

	return status, nil
}

func (c *Client) StudyUIStatus(ctx context.Context) (*ServiceResponse, error) {
	url := "https://study.dev.cutter.live/"
	req, _ := client.NewRequestWithContext(ctx, "GET", url, nil)

	status := &ServiceResponse{}
	if err := c.do(ctx, req, &status); err != nil {
		return nil, err
	}

	return status, nil
}
