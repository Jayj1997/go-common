package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// TODO: add configuration

func Get(url string) ([]byte, error) {

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 304 {
		return []byte(""), errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

func Post(url string, data interface{}) (string, error) {

	client := &http.Client{Timeout: 10 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result), err
}

func Put(url string, data interface{}) (string, error) {

	client := &http.Client{Timeout: 10 * time.Second}
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))

	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return "", err
	}

	resp, respErr := client.Do(req)

	if respErr != nil {
		return "", respErr
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	result, _ := ioutil.ReadAll(resp.Body)

	return string(result), err
}

func Delete(url string) (string, error) {

	req, _ := http.NewRequest("DELETE", url, nil)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	result, _ := ioutil.ReadAll(resp.Body)

	return string(result), err
}
