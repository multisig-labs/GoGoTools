package glacier

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/caarlos0/env/v11"
	"github.com/go-resty/resty/v2"
)

// Self-contained package for getting data from Glacier.

type Config struct {
	BaseURL string `env:"GLACIER_URL" envDefault:"https://glacier-api.avax.network"`
	ApiKey  string `env:"GLACIER_API_KEY"`
	Network string `env:"GLACIER_NETWORK" envDefault:"mainnet"`
}

var defaultConfig Config

func init() {
	if err := env.Parse(&defaultConfig); err != nil {
		log.Fatalf("error binding config to env: %v", err)
	}
}

type Client struct {
	cfg Config
	rc  *resty.Client
}

func NewClient() *Client {
	return NewClientWithConfig(defaultConfig)
}

func NewClientWithConfig(cfg Config) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultConfig.BaseURL
	}
	if cfg.Network == "" {
		cfg.Network = defaultConfig.Network
	}
	return &Client{
		cfg: cfg,
		rc:  createRestyClient(cfg.BaseURL, cfg.ApiKey),
	}
}

func (c *Client) FetchValidators(subnetID ids.ID) ([]Validator, error) {
	endpoint := fmt.Sprintf("/v1/networks/%s/l1Validators", c.cfg.Network)
	params := map[string]string{
		"subnetId":                    subnetID.String(),
		"includeInactiveL1Validators": "true",
	}
	allValidatorsPages, err := getAll[ListValidators](c.rc, endpoint, params, 10)
	if err != nil {
		return nil, err
	}
	allValidators := combineValidatorPages(allValidatorsPages)
	return allValidators, nil
}

// Client for connecting to Glacier.
func createRestyClient(baseURL string, apiKey string) *resty.Client {
	rc := resty.New()
	rc.SetBaseURL(baseURL)
	rc.SetTimeout(time.Duration(30 * time.Second))
	// For (very) verbose Resty logging
	// client.SetDebug(true)
	rc.SetHeaders(map[string]string{
		"Content-Type":      "application/json",
		"Accept":            "application/json",
		"User-Agent":        "gogopool/rialto",
		"x-glacier-api-key": apiKey,
	})
	return rc
}

func combineValidatorPages(pages []ListValidators) []Validator {
	var out []Validator
	for _, p := range pages {
		out = append(out, p.Validators...)
	}
	return out
}

func get(c *resty.Client, endpoint string, params map[string]string, result interface{}) (*resty.Response, error) {
	req := c.R()
	req.SetQueryParams(params)
	if result != nil {
		req = req.SetResult(result)
	}
	resp, err := req.Get(endpoint)

	if err != nil {
		// Example error string "Get "https://blah.dev/info": dial tcp: lookup blah.dev: no such host"
		return nil, fmt.Errorf("[Resty] %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("[Resty] Get %s %s", endpoint, resp)
	}
	return resp, nil
}

// Use pagination to get all the pages (as an array of the pages) from an endpoint.
func getAll[T Pageable](c *resty.Client, endpoint string, params map[string]string, maxPages int) ([]T, error) {
	const pageSize = "100"
	params["pageSize"] = pageSize

	allPages := make([]T, 0, maxPages)

	// Get first page
	resp, err := get(c, endpoint, params, nil)
	if err != nil {
		return nil, fmt.Errorf("getAll err: %w", err)
	}

	var aPage T
	err = json.Unmarshal(resp.Body(), &aPage)
	if err != nil {
		return nil, fmt.Errorf("unmarshal err: %w", err)
	}

	allPages = append(allPages, aPage)
	params["pageToken"] = aPage.NextPage()

	// Keep getting the rest of the pages
	for curPage := 1; curPage < maxPages && params["pageToken"] != ""; curPage++ {
		resp, err := get(c, endpoint, params, nil)
		if err != nil {
			return nil, fmt.Errorf("getAll curPage: %d  err: %w", curPage, err)
		}

		var aPage T
		err = json.Unmarshal(resp.Body(), &aPage)
		if err != nil {
			return nil, fmt.Errorf("unmarshal err: %w", err)
		}

		allPages = append(allPages, aPage)
		params["pageToken"] = aPage.NextPage()
	}

	return allPages, nil
}
