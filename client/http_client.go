package client

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type HTTPClient struct {
	Addr   string
	TLS    bool
	client *http.Client
}

func (hc *HTTPClient) Get(k string) (string, error) {
	q := url.Values{}
	q.Set("key", k)
	u := fmt.Sprintf("%s/get?%s", hc.Addr, q.Encode())
	resp, err := hc.http().Get(u)
	if err != nil {
		return "", err
	}
	return hc.strResponse(resp, k)
}

func (hc *HTTPClient) Set(k, v string) error {
	q := url.Values{}
	q.Set("key", k)
	q.Set("value", v)
	u := fmt.Sprintf("%s/set", hc.Addr)
	resp, err := hc.http().PostForm(u, q)
	if err != nil {
		return err
	}
	_, err = hc.strResponse(resp, k)
	return err
}

func (hc *HTTPClient) Del(k string) error {
	q := url.Values{}
	q.Set("key", k)
	u := fmt.Sprintf("%s/del", hc.Addr)
	resp, err := hc.http().PostForm(u, q)
	if err != nil {
		return err
	}
	_, err = hc.strResponse(resp, k)
	return err
}

func (hc *HTTPClient) strResponse(response *http.Response, k string) (string, error) {
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("错误：键[%s]%s", k, string(data))
	}
	return string(data), nil
}

func (hc *HTTPClient) http() *http.Client {
	if hc.client == nil {
		cli := &http.Client{Timeout: 30 * time.Second}
		if hc.TLS {
			cfg := &tls.Config{InsecureSkipVerify: true}
			trs := &http.Transport{
				TLSClientConfig:     cfg,
				IdleConnTimeout:     90 * time.Second,
				MaxIdleConns:        100,
				TLSHandshakeTimeout: 10 * time.Second,
			}
			cli.Transport = trs
		}
		hc.client = cli
	}
	return hc.client
}
