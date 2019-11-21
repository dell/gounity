package gounity

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"

	"github.com/dell/gounity/api"
	types "github.com/dell/gounity/payloads"
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

// Struct holds the configuration & REST Client.
type Client struct {
	configConnect *ConfigConnect
	api           api.Client
}

// Struct holds the endpoint & credential info.
type ConfigConnect struct {
	Endpoint string
	Username string
	Password string
}

var log *logrus.Logger

func SetLogger(l *logrus.Logger) {
	log = l
}

// Authenticate make a REST API call [/loginSessionInfo] to Unity to get authenticate the given credentials.
// The response contains the EMC-CSRF-TOKEN and the client caches it for further communication.
func (c *Client) Authenticate(configConnect *ConfigConnect) error {
	SetLogger(api.GetLogger())
	log.Info("Executing Authenticate REST client")
	c.configConnect = configConnect
	c.api.SetToken("")
	headers := make(map[string]string, 3)
	headers[api.AuthorizationHeader] = "Basic " + basicAuth(configConnect.Username, configConnect.Password)
	headers[api.XEmcRestClient] = "true"
	headers[api.HeaderKeyContentType] = api.HeaderValContentTypeJSON
	resp, err := c.api.DoAndGetResponseBody(context.Background(), http.MethodGet, api.UnityApiLoginSessionInfoUri, headers, nil)

	if err != nil {
		log.Errorf("Authentication error: %v", err)
		return err
	}

	if resp != nil {
		log.Info("Authentication response code:", resp.StatusCode)
		if err != nil {
			log.Errorf("Reading Authentication response body error:%v", err)
		}

		defer resp.Body.Close()

		switch {

		case resp.StatusCode >= 200 && resp.StatusCode <= 299:
			{
				log.Info("Authentication successful")
			}
		case resp.StatusCode == 401:
			{
				log.Error("Authentication failed")
				return status.Errorf(codes.Unauthenticated, "Unable to login to Unity. Verify username and password.")
			}
		default:
			log.Errorf("Authenticate error. Response: %v", c.api.ParseJSONError(resp))
			return c.api.ParseJSONError(resp)
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
func (c *Client) executeWithRetryAuthenticate(method, uri string, body, resp interface{}) error {
	headers := make(map[string]string, 2)
	headers[api.HeaderKeyAccept] = accHeader
	headers[api.HeaderKeyContentType] = conHeader
	headers[api.XEmcRestClient] = "true"
	log.Info("Invoking REST API server info Method: ", method, ", URI: ", uri)

	err := c.api.DoWithHeaders(context.Background(), method, uri, headers, body, resp)
	if err == nil {
		log.Info("Execution successful on Method: ", method, ", URI: ", uri)
		return nil
	}
	// check if we need to authenticate
	if e, ok := err.(*types.Error); ok {
		log.Info("Error in response", ", Method: ", method, ", URI: ", uri)
		log.WithError(err).Debugf("Got JSON error: %+v", e)
		if e.ErrorContent.HTTPStatusCode == 401 {
			log.Info("need to re-authenticate")
			// Authenticate then try again
			if err := c.Authenticate(c.configConnect); err != nil {
				return fmt.Errorf("authentication failure due to: %v", err)
			} else {
				log.Info("Authentication success")
			}
			return c.api.DoWithHeaders(context.Background(), method, uri, headers, body, resp)
		}
	} else {
		log.Error("Error is not a type of \"*types.Error\". Error:", err)
	}
	log.WithError(err).Error("failed to invoke Unity REST API server")

	return err
}

func (c *Client) SetToken(token string) {
	c.api.SetToken(token)
}

func (c *Client) GetToken() string {
	return c.api.GetToken()
}

// NewClient initialize the new REST Client with default options.
func NewClient() (client *Client, err error) {
	return NewClientWithArgs(
		os.Getenv("GOUNITY_ENDPOINT"),
		os.Getenv("GOUNITY_INSECURE") == "true",
		os.Getenv("GOUNITY_USECERTS") == "true")
}

// NewClientWithArgs initialize the new REST Client with the given arguments.
func NewClientWithArgs(endpoint string, insecure, useCerts bool) (client *Client, err error) {
	if showHTTP {
		debug = true
	}
	SetLogger(api.GetLogger())
	fields := map[string]interface{}{
		"endpoint": endpoint,
		"insecure": insecure,
		"useCerts": useCerts,
		"debug":    debug,
		"showHTTP": showHTTP,
	}

	log.WithFields(fields).Debug("unity client init")

	if endpoint == "" {
		log.WithFields(fields).Error("endpoint is required")
		return nil,
			withFields(fields, "endpoint is required")
	}

	opts := api.ClientOptions{
		Insecure: insecure,
		UseCerts: useCerts,
		ShowHTTP: showHTTP,
	}

	ac, err := api.New(context.Background(), endpoint, opts, debug)
	if err != nil {
		log.Errorf("Unable to create HTTP client %v", err)
		return nil, err
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
