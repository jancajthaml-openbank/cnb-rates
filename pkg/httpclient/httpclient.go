// Copyright (c) 2016-2018, Jan Cajthaml <jan.cajthaml@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httpclient

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type HttpClient struct {
	client *http.Client
}

func New(secured bool) *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout:   5 * time.Second,
			Transport: getTransport(secured),
		},
	}
}

func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *HttpClient) Get(url string) (contents io.Reader, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Runtime Error %v", r)
		}

		if err != nil && resp != nil {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		} else if resp == nil && err != nil {
			err = fmt.Errorf("Runtime Error no response")
		}

		if err != nil {
			contents = bytes.NewReader(nil)
		} else {
			var data []byte
			data, err = ioutil.ReadAll(resp.Body)
			contents = bytes.NewReader(data)
			resp.Body.Close()
		}
	}()

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	resp, err = c.client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("%d", resp.StatusCode))
		return
	}

	return
}

func getTransport(secured bool) *http.Transport {
	if secured {
		return &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 5 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify:       false,
				MinVersion:               tls.VersionTLS12,
				PreferServerCipherSuites: false,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
					tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
					tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
					tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
					tls.TLS_RSA_WITH_AES_128_CBC_SHA,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
					tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				},
			},
		}
	}

	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
	}
}
