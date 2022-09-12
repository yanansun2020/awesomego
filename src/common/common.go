package common

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

var httpClient *http.Client

const NameServer = "https://names.mcquay.me/api/v0/"
const JokeAPIServer = "http://api.icndb.com/jokes/random"

var secureTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
	TLSClientConfig:     &tls.Config{InsecureSkipVerify: false},
}
var unsecureTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
	TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
}

type HttpReq struct {
	Url     string
	Method  string
	Body    interface{}
	Headers http.Header
}
type HttpResp struct {
	Body    []byte
	Status  int
	Headers http.Header
}

func HttpClientInit() (err error) {
	//httpClient = GetHttpClient(false)
	httpClient = &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	if httpClient == nil {
		log.Println("httpClient is nil")
		return
	}
	return
}

func GetHttpClient(sslverify bool) *http.Client {
	if sslverify {
		return &http.Client{
			Timeout:   time.Second * 10,
			Transport: secureTransport,
		}

	} else {
		return &http.Client{
			Timeout:   time.Second * 10,
			Transport: unsecureTransport,
		}
	}
}

func SendHttpReq(in HttpReq) (out HttpResp, err error) {
	var inData *bytes.Reader
	if in.Body != nil {
		b, _ := json.Marshal(in.Body)
		inData = bytes.NewReader(b)
	} else {
		inData = bytes.NewReader(nil)
	}
	req, err := http.NewRequest(in.Method, in.Url, inData)
	if err != nil {
		log.Printf("error while create req %s", err.Error())
		return
	}
	if len(in.Headers) != 0 {
		for k, _ := range in.Headers {
			req.Header.Add(k, in.Headers[k][0])
		}
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Request failed. Err: %s", err.Error())
		return
	}

	out.Body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body. Err: %s", err.Error())
		return
	}
	resp.Body.Close()
	out.Headers = resp.Header
	out.Status = resp.StatusCode
	return
}
