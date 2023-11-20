package util

import (
	"choccy/server/database"
	"choccy/server/database/model"
	"crypto/tls"
	"golang.org/x/net/http/httpproxy"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetHttpClient(timeout time.Duration) (*http.Client, error) {
	var setting model.Setting
	result := database.DB.Take(&setting)
	if result.Error != nil {
		return nil, result.Error
	}

	transport := &http.Transport{}
	if strings.TrimSpace(httpproxy.FromEnvironment().HTTPSProxy) != "" {
		proxyUrl, err := url.Parse(httpproxy.FromEnvironment().HTTPSProxy)
		if err != nil {
			return nil, err
		}
		transport = &http.Transport{
			ForceAttemptHTTP2: true,
			Proxy:             http.ProxyURL(proxyUrl),
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: setting.SkipVerifyTLS},
		}
	}
	httpClient := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
	return httpClient, nil
}
