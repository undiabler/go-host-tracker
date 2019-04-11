package htracker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	// HostTrackerURL define global url for api
	HostTrackerURL = "https://www.host-tracker.com"
)

type authToken struct {
	Ticket             string `json:"ticket"`
	ExpirationUnixTime int64  `json:"expirationUnixTime"`
}

// HostTrackerClient extend http client with auth required for host-tracker
type HostTrackerClient struct {
	http.Client

	Login    string
	Password string

	token *authToken
}

var hostClient = &http.Client{
	Timeout: time.Second * 5,
}

// Do http request, but check token (and generate new if required) first
func (htc *HostTrackerClient) Do(req *http.Request) (*http.Response, error) {

	if htc.token == nil || htc.token.ExpirationUnixTime <= time.Now().UTC().Unix() {
		if err := htc.auth(); err != nil {
			return nil, err
		}
	}

	if htc.token == nil {
		return nil, fmt.Errorf("cannot get auth token")
	}

	req.Header.Add("Authorization", "bearer "+htc.token.Ticket)

	return htc.Client.Do(req)
}

func (htc *HostTrackerClient) auth() error {

	bodyMap := map[string]string{
		"login":    htc.Login,
		"password": htc.Password,
	}

	bodyByte, err := json.Marshal(bodyMap)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", HostTrackerURL+"/api/web/v1/users/token", bytes.NewBuffer(bodyByte))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := htc.Client.Do(req)
	if err != nil {
		return err
	}
	bodyResp, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var newToken *authToken
	if err := json.Unmarshal(bodyResp, &newToken); err != nil {
		return err
	}

	htc.token = newToken

	return nil
}

// NewHttpTask creates new task for host-tracker
func (htc *HostTrackerClient) NewHttpTask(params map[string]interface{}) (string, error) {

	bodyByte, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", HostTrackerURL+"/api/web/v1/tasks/http", bytes.NewBuffer(bodyByte))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := htc.Do(req)
	if err != nil {
		return "", err
	}
	bodyResp, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != 201 {
		return "", fmt.Errorf("invalid answer code: %d", resp.StatusCode)
	}

	type tempjs struct {
		ID string `json:"id"`
	}
	var resultID tempjs

	if err := json.Unmarshal(bodyResp, &resultID); err != nil {
		return "", err
	}

	return resultID.ID, nil
}
