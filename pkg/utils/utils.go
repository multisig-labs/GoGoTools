package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/multisig-labs/gogotools/pkg/constants"
	"github.com/tidwall/gjson"
)

func Fetch(url string, body string) (string, error) {
	client := resty.New()
	// client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(30 * time.Second)

	var resp *resty.Response
	var err error

	if body == "" {
		resp, err = client.R().
			EnableTrace().
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json").
			Get(url)
	} else {
		resp, err = client.R().
			EnableTrace().
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json").
			SetBody(body).
			Post(url)
	}

	return resp.String(), err
}

func FetchRPC(url string, method string, params string) (string, error) {
	client := resty.New()
	// client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(30 * time.Second)

	var resp *resty.Response
	var err error

	if params == "" {
		params = "{}"
	}

	body := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id"     : %d,
		"method" : "%s",
		"params" : %s
	}`, time.Now().Unix(), method, params)

	resp, err = client.R().
		EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(body).
		Post(url)

	if resp.IsError() {
		return "", fmt.Errorf("fetch error %d: %s %s", resp.StatusCode(), url, body)
	}
	return resp.String(), err
}

func FetchRPCGJSON(url string, method string, params string) (*gjson.Result, error) {
	s, err := FetchRPC(url, method, params)
	if err != nil {
		return nil, err
	}
	out := gjson.Parse(s)
	return &out, nil
}

func EnsureFileExists(path string) bool {
	match, err := filepath.Glob(path)
	if err == nil && match != nil {
		return true
	}
	fmt.Fprintf(os.Stderr, "File does not exist: %s\n", path)
	os.Exit(1)
	return false
}

func LinkFile(src, dest string) error {
	return os.Symlink(src, dest)
}

func CopyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	if err = out.Sync(); err != nil {
		return err
	}
	if err = out.Chmod(constants.DefaultPerms755); err != nil {
		return err
	}
	return nil
}

// func AvaKeyToEthKey(key *crypto.PrivateKeySECP256K1R) common.Address {
// 	pubk := key.ToECDSA().PublicKey
// 	addr := ethcrypto.PubkeyToAddress(pubk)
// 	return addr
// }
