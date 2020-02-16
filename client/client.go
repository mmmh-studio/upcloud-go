package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

const (
	defaultTimeout = 10
)

type Endpoint interface {
	// Path for the API endpoint.
	Path() string
}

// Payload send as request body.
type Payload interface {
	// Validate values of all fields.
	Validate() []FieldError
}

type EndpointPayloadMarshaler interface {
	Endpoint
	Payload
	json.Marshaler
}

type Client struct {
	http *http.Client

	user     string
	password string
}

func New(user, password string) *Client {
	httpClient := cleanhttp.DefaultClient()
	httpClient.Timeout = time.Second * defaultTimeout

	return &Client{
		http:     httpClient,
		user:     user,
		password: password,
	}
}

func (c *Client) Delete(endpoint Endpoint) error {
	r, err := http.NewRequest("DELETE", endpoint.Path(), nil)
	if err != nil {
		return err
	}

	_, err = c.request(r)
	return err
}

func (c *Client) Get(endpoint Endpoint, response json.Unmarshaler) error {
	if err := checkResponse(response); err != nil {
		return err
	}
	req, err := http.NewRequest("GET", endpoint.Path(), nil)
	if err != nil {
		return err
	}

	resBody, err := c.request(req)
	if err != nil {
		return err
	}

	return json.Unmarshal(resBody, response)
}

func (c *Client) Post(
	payload EndpointPayloadMarshaler,
	response json.Unmarshaler,
) error {
	if errors := payload.Validate(); len(errors) != 0 {
		return &ValidationError{
			Name:        typeName(payload),
			FieldErrors: errors,
		}
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", payload.Path(), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	resBody, err := c.request(req)
	if err != nil {
		return err
	}

	if response == nil {
		return nil
	}

	return json.Unmarshal(resBody, response)
}

func (c *Client) request(req *http.Request) ([]byte, error) {
	req.SetBasicAuth(c.user, c.password)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(fmt.Sprintf("failed to read error body: %s", err))
		}

		return nil, &Error{code: res.StatusCode, message: res.Status, body: body}
	}

	return ioutil.ReadAll(res.Body)
}

func checkResponse(res interface{}) error {
	switch reflect.ValueOf(res).Kind() {
	case reflect.Ptr:
		return nil
	case reflect.Struct:
		return fmt.Errorf("response is non-pointer struct")
	default:
		return fmt.Errorf("response needs to be pointer to struct")
	}
}
