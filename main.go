package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"gopkg.in/yaml.v2"
)

type TestRequest struct {
	Headers []map[string]string
	Timeout string
	Method  string
	URL     string
	Body    string
}

type TestResponse struct {
	Headers    []map[string]string
	Timeout    string
	Method     string
	URL        string
	Body       string
	StatusCode int
}

type TestThen struct {
	CaseID string
	Args   []map[string]interface{}
}

type TestDefinition struct {
	Title    string `yaml: "title"`
	ID       string
	Request  TestRequest `yaml: "request"`
	Response TestResponse
	Then     TestThen
}

type TestCase struct {
	Case TestDefinition `yaml: "case"`
}

func BuildRequest(testRequest TestRequest) *http.Response {
	hClient := http.Client{}

	r, err := http.NewRequest(testRequest.Method, testRequest.URL, nil)
	if err != nil {
		panic(err)
	}

	addTimeout(&hClient, &testRequest)
	addHeaders(r, &testRequest)
	addUrl(r, &testRequest)

	res, err := hClient.Do(r)
	if err != nil {
		panic(err)
	}

	return res
}

func ValidateResponse(res *http.Response, tresponse *TestResponse) {
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resBody))
	fmt.Println(res.StatusCode)

	if tresponse.StatusCode != 0 {
		if res.StatusCode != tresponse.StatusCode {
			panic("status codes not matched " + tresponse.URL)
		}
	}

	if tresponse.Body != "" {
		var parsedResponseBody map[string]interface{}
		err := json.Unmarshal(resBody, &parsedResponseBody)
		if err != nil {
			panic(err)
		}

		var desiredResponseBody map[string]interface{}
		err = json.Unmarshal([]byte(tresponse.Body), &desiredResponseBody)
		if err != nil {
			panic(err)
		}

		for key := range desiredResponseBody {
			if parsedResponseBody[key] != desiredResponseBody[key] {
				panic("values are not matched " + key)
			}
		}
	}
}

func addTimeout(hClient *http.Client, tr *TestRequest) {
	if tr.Timeout != "" {
		duration := ""

		if duration != "" {
			duration, err := time.ParseDuration(tr.Timeout)
			if err != nil {
				panic(err)
			}
			hClient.Timeout = duration
		}
	}
}

func addHeaders(r *http.Request, tr *TestRequest) {
	for _, h := range tr.Headers {
		for key := range h {
			r.Header.Add(key, h[key])
		}
	}
}

func addUrl(r *http.Request, tr *TestRequest) {
	url, err := url.Parse(tr.URL)
	if err != nil {
		panic(err)
	}
	r.URL = url
}

func main() {
	content, ioerr := ioutil.ReadFile("test.yaml")
	if ioerr != nil {
		panic(ioerr)
	}

	testCases := []TestCase{}
	err := yaml.Unmarshal(content, &testCases)
	if err != nil {
		panic(err)
	}

	for _, cases := range testCases {
		req := BuildRequest(cases.Case.Request)
		ValidateResponse(req, &cases.Case.Response)
	}
}
