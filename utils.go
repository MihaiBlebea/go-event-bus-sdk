package eventbus

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func post(endpoint string, payload interface{}, response interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer(b),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("request status is not 200")
	}

	if response != nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(body, response); err != nil {
			return err
		}
	}

	return nil
}

func get(endpoint string, response interface{}) error {
	req, err := http.NewRequest(
		"GET",
		endpoint,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("request status is not 200")
	}

	if err := json.Unmarshal(body, response); err != nil {
		return err
	}

	return nil
}
