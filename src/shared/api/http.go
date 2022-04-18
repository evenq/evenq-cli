package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var client *http.Client

const apiRoot = "https://api.evenq.io"
const authTokenHeader = "x-evenq-token"
const orgHeader = "x-evenq-org"

func getClient() *http.Client {
	if client == nil {
		client = &http.Client{
			Timeout: time.Second * 30,
		}
	}

	return client
}

func Get(ctx context.Context, path string, dest interface{}) error {
	token, ok := getToken(ctx)
	if !ok {
		return errors.New("context is missing auth token")
	}

	req, err := http.NewRequest("GET", apiRoot+path, nil)
	if err != nil {
		return err
	}

	req.Header.Set(authTokenHeader, token)
	req.Header.Set(orgHeader, getOrg(ctx))
	req.Header.Set("content-type", "application/json")

	resp, err := getClient().Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()

	// debugging stuff
	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)
	dec := json.NewDecoder(tee)

	if isDebug(ctx) {
		log.Println(path, token, getOrg(ctx), resp.StatusCode, buf.String())
	}

	return dec.Decode(&dest)
}

func Post(ctx context.Context, path string, data interface{}, dest interface{}) (*http.Response, error) {
	token, _ := getToken(ctx)

	body, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", apiRoot+path, r)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if token != "" {
		req.Header.Set(authTokenHeader, token)
	}
	req.Header.Set(orgHeader, getOrg(ctx))
	req.Header.Set("content-type", "application/json")

	resp, err := getClient().Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)
	dec := json.NewDecoder(tee)

	if isDebug(ctx) {
		log.Println(path, token, getOrg(ctx), resp.StatusCode, buf.String())
	}

	return resp, dec.Decode(&dest)
}

func IsSuccess(data map[string]interface{}) bool {
	if val, ok := data["success"].(bool); ok {
		return val == true
	}

	return false
}

func isDebug(ctx context.Context) bool {
	if val, ok := ctx.Value("isDebug").(bool); ok {
		return val
	}

	return false
}
