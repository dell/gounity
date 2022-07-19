package gounity

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/dell/gounity/util"

	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emcCsrfToken = "EMC-CSRF-TOKEN"
)

var (
	accHeader string
	conHeader string

	errNoLink   = errors.New("error: problem finding link")
	debug, _    = strconv.ParseBool(os.Getenv("GOUNITY_DEBUG"))
	showHTTP, _ = strconv.ParseBool(os.Getenv("GOUNITY_SHOWHTTP"))
)

//Client Struct holds the configuration & REST Client.
type Client struct {
	configConnect *ConfigConnect
	api           api.Client
}

//ConfigConnect Struct holds the endpoint & credential info.
type ConfigConnect struct {
	Endpoint string
	Username string
	Password string
}

// Authenticate make a REST API call [/loginSessionInfo] to Unity to get authenticate the given credentials.
// The response contains the EMC-CSRF-TOKEN and the client caches it for further communication.
func (c *Client) Authenticate(ctx context.Context, configConnect *ConfigConnect) error {
	log := util.GetRunIDLogger(ctx)
	log.Debug("Executing Authenticate REST client")
	c.configConnect = configConnect
	c.api.SetToken("")
	headers := make(map[string]string, 3)
	headers[api.AuthorizationHeader] = "Basic " + basicAuth(configConnect.Username, configConnect.Password)
	headers[api.XEmcRestClient] = "true"
	headers[api.HeaderKeyContentType] = api.HeaderValContentTypeJSON
	resp, err := c.api.DoAndGetResponseBody(ctx, http.MethodGet, api.UnityAPILoginSessionInfoURI, headers, nil)

	if err != nil {
		return fmt.Errorf("authentication error: %v", err)
	}

	if resp != nil {
		log.Debugf("Authentication response code: %d", resp.StatusCode)
		if err != nil {
			log.Errorf("Reading Authentication response body error:%v", err)
		}

		defer resp.Body.Close()

		switch {
		case resp.StatusCode >= 200 && resp.StatusCode <= 299:
			{
				log.Debug("Authentication successful")
			}
		case resp.StatusCode == 401:
			{
				return status.Errorf(codes.Unauthenticated, "Authentication failed. Unable to login to Unity. Verify username and password.")
			}
		default:
			return fmt.Errorf("authenticate error. Response: %v", c.api.ParseJSONError(ctx, resp))
		}

		c.api.SetToken(resp.Header.Get(emcCsrfToken))
	} else {
		log.Errorf("Authenticate error: Nil response received")
	}
	return nil
}

// basicAuth converts the given username & password to Base64 encoded string.
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// GetJSONWithRetry method responsible to make the given API call to Unity REST API Server.
// In case if the given EMC-CSRF-TOKEN becomes invalid, retries the same operation after performing authentication.
func (c *Client) executeWithRetryAuthenticate(ctx context.Context, method, uri string, body, resp interface{}) error {
	log := util.GetRunIDLogger(ctx)
	headers := make(map[string]string, 2)
	headers[api.HeaderKeyAccept] = accHeader
	headers[api.HeaderKeyContentType] = conHeader
	headers[api.XEmcRestClient] = "true"
	log.Debug("Invoking REST API server info Method: ", method, ", URI: ", uri)
	err := c.api.DoWithHeaders(ctx, method, uri, headers, body, resp)
	if err == nil {
		log.Debug("Execution successful on Method: ", method, ", URI: ", uri)
		return nil
	}
	// check if we need to authenticate
	if e, ok := err.(*types.Error); ok {
		log.Debugf("Error in response. Method:%s URI:%s Error: %v JSON Error: %+v", method, uri, err, e)
		if e.ErrorContent.HTTPStatusCode == 401 {
			log.Debug("need to re-authenticate")
			// Authenticate then try again
			if err := c.Authenticate(ctx, c.configConnect); err != nil {
				return fmt.Errorf("authentication failure due to: %v", err)
			}
			log.Debug("Authentication success")
			return c.api.DoWithHeaders(ctx, method, uri, headers, body, resp)
		}
	} else {
		log.Error("Error is not a type of \"*types.Error\". Error:", err)
	}
	log.WithError(err).Error("failed to invoke Unity REST API server")

	return err
}

//SetToken function sets token
func (c *Client) SetToken(token string) {
	c.api.SetToken(token)
}

//GetToken function gets token
func (c *Client) GetToken() string {
	return c.api.GetToken()
}

// NewClient initialize the new REST Client with default options.
func NewClient(ctx context.Context) (client *Client, err error) {
	insecure, _ := strconv.ParseBool(os.Getenv("GOUNITY_INSECURE"))
	return NewClientWithArgs(ctx, os.Getenv("GOUNITY_ENDPOINT"), insecure)
}

// NewClientWithArgs initialize the new REST Client with the given arguments.
func NewClientWithArgs(ctx context.Context, endpoint string, insecure bool) (client *Client, err error) {
	log := util.GetRunIDLogger(ctx)
	if showHTTP {
		debug = true
	}

	fields := map[string]interface{}{
		"endpoint": endpoint,
		"insecure": insecure,
		"debug":    debug,
		"showHTTP": showHTTP,
	}

	log.WithFields(fields).Debug("unity client init")

	if endpoint == "" {
		log.WithFields(fields).Error("endpoint is required")
		return nil, withFields(fields, "endpoint is required")
	}

	opts := api.ClientOptions{
		Insecure: insecure,
		ShowHTTP: showHTTP,
	}

	ac, err := api.New(ctx, endpoint, opts, debug)
	if err != nil {
		return nil, fmt.Errorf("unable to create HTTP client %v", err)
	}

	client = &Client{
		api:           ac,
		configConnect: &ConfigConnect{},
	}
	conHeader = api.HeaderValContentTypeJSON
	return client, nil
}

func withFields(fields map[string]interface{}, message string) error {
	return withFieldsE(fields, message, nil)
}

func withFieldsE(fields map[string]interface{}, message string, inner error) error {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	if inner != nil {
		fields["inner"] = inner
	}

	x := 0
	l := len(fields)

	var b bytes.Buffer
	for k, v := range fields {
		if x < l-1 {
			b.WriteString(fmt.Sprintf("%s=%v,", k, v))
		} else {
			b.WriteString(fmt.Sprintf("%s=%v", k, v))
		}
		x = x + 1
	}

	return fmt.Errorf("%s %s", message, b.String())
}
