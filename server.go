// Copyright 2022 Stock Parfait

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutil

import (
	"net/http"
	"net/http/httptest"
	"net/url"
)

// TestServer is the handle for a test server.  The server returns the status
// code and response body from the matching sequence, and then repeats the last
// one in the sequence.
//
// Note, that this creates a real HTTP server running locally, and the client
// fetches real URLs from the network. Always use the base URL provided by URL()
// method to avoid calling random servers on the Internet in your tests.
//
// Tip: pay attention to the status codes when setting a response body: some
// response codes do not allow a response body. Always check for BodyWriteError
// in your tests to catch this problem.
type TestServer struct {
	ResponseStatusMap map[string][]int    // URL path -> status code sequence
	ResponseBodyMap   map[string][]string // URL path -> response body sequence
	ResponseStatus    []int               // default status codes sequence
	ResponseBody      []string            // default response body sequnece
	RequestPath       string              // path in the request URL
	RequestQuery      url.Values          // query received by the server in the request
	Flushed           bool                // whether the body write was flushed
	BodyWriteBytes    int                 // number of body bytes written
	BodyWriteError    error               // error value from writing body
	Server            *httptest.Server
}

// Close the handle's test server.
func (h *TestServer) Close() {
	h.Server.Close()
}

// Client is the test server's HTTP client to be used in tests.
func (h *TestServer) Client() *http.Client {
	return h.Server.Client()
}

// URL returns the base test server URL.
func (h *TestServer) URL() string {
	return h.Server.URL
}

// NewTestServer creates and starts a new test server.
func NewTestServer() *TestServer {
	h := TestServer{
		ResponseStatus:    []int{http.StatusOK},
		ResponseBody:      []string{""},
		ResponseStatusMap: make(map[string][]int),
		ResponseBodyMap:   make(map[string][]string),
	}
	h.Server = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			statusSeq, statusInMap := h.ResponseStatusMap[r.URL.Path]
			if !statusInMap {
				statusSeq = h.ResponseStatus
			}
			status := http.StatusOK
			if len(statusSeq) > 0 {
				status = statusSeq[0]
			}
			if len(statusSeq) > 1 {
				statusSeq = statusSeq[1:]
				if statusInMap {
					h.ResponseStatusMap[r.URL.Path] = statusSeq
				} else {
					h.ResponseStatus = statusSeq
				}
			}
			bodySeq, bodyInMap := h.ResponseBodyMap[r.URL.Path]
			useResponseBody := !bodyInMap || len(bodySeq) == 0
			if useResponseBody {
				bodySeq = h.ResponseBody
			}
			body := ""
			if len(bodySeq) > 0 {
				body = bodySeq[0]
			}
			if len(bodySeq) > 1 {
				bodySeq = bodySeq[1:]
				if useResponseBody {
					h.ResponseBody = bodySeq
				} else {
					h.ResponseBodyMap[r.URL.Path] = bodySeq
				}
			}
			w.WriteHeader(status)
			h.RequestPath = r.URL.Path
			h.RequestQuery = r.URL.Query()
			h.BodyWriteBytes, h.BodyWriteError = w.Write([]byte(body))
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
				h.Flushed = true
			}
		}))
	return &h
}
