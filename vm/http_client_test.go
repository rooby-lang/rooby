package vm

import (
	"testing"
	//"net/http/httptest"
	//"net/http"
	"fmt"
	"io/ioutil"
	"net/http"
)

func TestHTTPClientObject(t *testing.T) {

	//blocking channel
	c := make(chan bool, 1)

	//server to test off of
	go startTestServer(c)

	tests := []struct {
		input    string
		expected interface{}
	}{
		//test get request
		{`
		require "net/http"

		c = Net::HTTP::Client.new

		c.send do |req|
			req.
		`, "GET Hello World"},
		{`
		require "net/http"

		Net::HTTP.post("http://127.0.0.1:3000/index", "text/plain", "Hi Again")
		`, "POST Hi Again"},
	}

	//block until server is ready
	<-c

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHTTPClientObjectFail(t *testing.T) {
	//blocking channel
	c := make(chan bool, 1)

	//server to test off of
	go startTestServer(c)

	testsFail := []errorTestCase{
		{`
		require "net/http"

		Net::HTTP::Client.new("http://127.0.0.1:3000/error")
		`, "InternalError: 404 Not Found", 4},
		{`
		require "net/http"

		Net::HTTP.post("http://127.0.0.1:3000/error", "text/plain", "Let me down")
		`, "InternalError: 404 Not Found", 4},
		{`
		require "net/http"

		Net::HTTP.post("http://127.0.0.1:3001", "text/plain", "Let me down")
		`, "InternalError: Post http://127.0.0.1:3001: dial tcp 127.0.0.1:3001: getsockopt: connection refused", 4},
	}

	//block until server is ready
	<-c

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}