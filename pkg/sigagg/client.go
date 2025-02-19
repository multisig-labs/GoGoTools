package sigagg

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/go-resty/resty/v2"
	"github.com/multisig-labs/gogotools/pkg/utils"
)

// Package for talking to Glacier signature-aggregator service
// https://glacier-api.avax.network/api#/

// interface compliance
var (
	_ Client = (*client)(nil)
)

type Client interface {
	AggregateSignatures(msg *warp.UnsignedMessage, subnetID ids.ID, justification []byte) (*warp.Message, error)
}

type client struct {
	rc *resty.Client
}

func NewClient(url string, headers ...map[string]string) (Client, error) {
	rc := createRestyClient(url)
	// Add custom headers if provided
	if len(headers) > 0 {
		rc.SetHeaders(headers[0])
	}
	return &client{
		rc: rc,
	}, nil
}

// Use the Glacier API to aggregate signatures
func (c *client) AggregateSignatures(msg *warp.UnsignedMessage, subnetID ids.ID, justification []byte) (*warp.Message, error) {
	params := map[string]interface{}{
		"message":         utils.BytesToHex(msg.Bytes()),
		"signingSubnetId": subnetID.String(),
	}
	if justification != nil {
		params["justification"] = utils.BytesToHex(justification)
	}
	resp, err := post(c.rc, "", params)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w  Body: %s", err, resp.String())
	}
	signedMessage, ok := result["signedMessage"].(string)
	if !ok {
		return nil, fmt.Errorf("signedMessage key not found in response: %s", resp.String())
	}
	msgBytes := utils.HexToBytes(signedMessage)
	msgSigned, err := warp.ParseMessage(msgBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing signedMessage: %w", err)
	}
	return msgSigned, nil
}

func createRestyClient(baseURL string) *resty.Client {
	rc := resty.New()
	rc.SetBaseURL(baseURL)
	rc.SetTimeout(time.Duration(30 * time.Second))
	rc.SetRetryAfter(nil) // default is exponential backoff with jitter
	// For (very) verbose Resty logging
	if os.Getenv("DEBUG") == "1" {
		rc.SetDebug(true)
		rc.EnableGenerateCurlOnDebug()
	}
	rc.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
		"User-Agent":   "gogopool/gogotools",
	})
	return rc
}

func post(c *resty.Client, endpoint string, params map[string]interface{}) (*resty.Response, error) {
	resp, err := c.R().SetBody(params).Post(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode(), resp.String())
	}
	return resp, nil
}
