package metalcloud

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ybbus/jsonrpc"
)

//DefaultEndpoint returns the default Bigstep Metalcloud endpoint
func DefaultEndpoint() string {
	return "https://api.bigstep.com/metal-cloud"
}

//Client sruct defines a metalcloud client
type Client struct {
	rpcClient jsonrpc.RPCClient
	user      string
	apiKey    string
	endpoint  string
}

//GetMetalcloudClient returns a metal cloud client
func GetMetalcloudClient(user string, apiKey string, endpoint string, loggingEnabled bool) (*Client, error) {

	if user == "" {
		return nil, errors.New("user cannot be an empty string! It is typically in the form of user's email address")
	}

	if apiKey == "" {
		return nil, errors.New("apiKey cannot be empty string")
	}

	if endpoint == "" {
		return nil, errors.New("endpoint cannot be an empty string! It is typically in the form of user's email address")
	}

	_, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return nil, err
	}

	transport := &signatureAdderRoundTripper{
		APIKey:         apiKey,
		LoggingEnabled: loggingEnabled,
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	rpcClient := jsonrpc.NewClientWithOpts(endpoint, &jsonrpc.RPCClientOpts{
		HTTPClient: httpClient,
	})

	return &Client{
		rpcClient: rpcClient,
		user:      user,
		apiKey:    apiKey,
		endpoint:  endpoint,
	}, nil

}

//GetUserEmail returns the user configured for this connection
func (c *Client) GetUserEmail() string {
	return c.user
}

//GetEndpoint returns the endpoint configured for this connection
func (c *Client) GetEndpoint() string {
	return c.endpoint
}

type signatureAdderRoundTripper struct {
	APIKey string
	http.RoundTripper
	LoggingEnabled bool
	DryRun         bool
}

func (c *signatureAdderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	components := strings.Split(c.APIKey, ":")

	var strKeyMetaData *string

	strKeyMetaData = nil

	if len(components) > 1 {
		strKeyMetaData = &components[0]
	}

	key := []byte(c.APIKey)

	// Read the content
	var message []byte
	if req.Body != nil {
		message, _ = ioutil.ReadAll(req.Body)
	}

	if c.LoggingEnabled {
		log.Println(string(message))
	}

	//force close connection. This will avoid the keep-alive related issues for go < 1.6 https://go-review.googlesource.com/c/go/+/3210
	req.Close = true

	// Restore the io.ReadCloser to its original state
	req.Body = ioutil.NopCloser(bytes.NewBuffer(message))

	hmac := hmac.New(md5.New, key)
	hmac.Write(message)

	var signature = hex.EncodeToString(hmac.Sum(nil))

	values, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		log.Fatal(err)
	}

	if strKeyMetaData != nil {
		signature = *strKeyMetaData + ":" + signature
	}

	values.Add("verify", signature)

	url := req.URL

	url.RawQuery = values.Encode()

	req.URL = url

	var resp *http.Response = nil

	if !c.DryRun {
		resp, err = http.DefaultTransport.RoundTrip(req)

	}

	if c.LoggingEnabled {
		//log the reply
		if resp.Body != nil {
			message, _ = ioutil.ReadAll(resp.Body)
		}

		log.Println(string(message))

		// Restore the io.ReadCloser to its original state
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(message))
	}

	return resp, err
}
