package client

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
)

const APIUrl = "https://gapp.redislabs.com"

type Client struct {
	httpClient *http.Client
}

func New(email, password string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", APIUrl+"/Login", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("email", email)
	q.Add("password", password)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non-200 response %d from login.", resp.StatusCode)
	}
	return &Client{httpClient: client}, nil
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%q", dump)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	dump, err = httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%q", dump)

	return resp, nil
}
