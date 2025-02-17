package uptime

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ava-labs/avalanchego/utils/rpc"
	"github.com/ava-labs/subnet-evm/plugin/evm/client"
)

type Client struct {
	validatorsUrl string

	validatorsRequester rpc.EndpointRequester
	queryParams         url.Values
}

// url should be like https://node.myl1.network/ext/bc/<blockchainID>/validators
func NewClient(validatorsURL string) (*Client, error) {
	parsedURL, err := url.Parse(validatorsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse validators URL: %w", err)
	}

	requester := rpc.NewEndpointRequester(validatorsURL)
	return &Client{
		validatorsUrl:       validatorsURL,
		validatorsRequester: requester,
		queryParams:         parsedURL.Query(),
	}, nil
}

func (c *Client) GetCurrentValidators(ctx context.Context, options ...rpc.Option) ([]client.CurrentValidator, error) {
	res := &client.GetCurrentValidatorsResponse{}

	for key := range c.queryParams {
		options = append(options, rpc.WithQueryParam(key, c.queryParams.Get(key)))
	}

	err := c.validatorsRequester.SendRequest(ctx, "validators.getCurrentValidators", &client.GetCurrentValidatorsRequest{}, res, options...)
	if err != nil {
		return nil, fmt.Errorf("error fetching from %s: %w", c.validatorsUrl, err)
	}
	return res.Validators, nil
}
