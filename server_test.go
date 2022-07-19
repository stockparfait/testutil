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
	"bytes"
	"net/http"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func URLQuery(uri string, query url.Values) string {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		panic(err)
	}
	req.URL.RawQuery = query.Encode()
	return req.URL.String()
}

func TestTestServer(t *testing.T) {
	t.Parallel()

	Convey("Server works", t, func() {
		server := NewTestServer()
		defer server.Close()

		client := server.Client()
		URL := server.URL()
		query := make(url.Values)

		Convey("default value", func() {
			query["key"] = []string{"v1", "v2"}
			path := "/test/path"
			resp, err := client.Get(URLQuery(URL+path, query))
			So(err, ShouldBeNil)
			defer resp.Body.Close()

			var body bytes.Buffer
			body.ReadFrom(resp.Body)
			So(server.BodyWriteError, ShouldBeNil)
			So(body.String(), ShouldEqual, "")
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(server.RequestPath, ShouldEqual, path)
			So(server.RequestQuery, ShouldResemble, query)
		})

		Convey("sequence values", func() {
			server.ResponseStatus = []int{http.StatusNotFound, http.StatusOK}
			server.ResponseBody = []string{"404", "I'm OK"}

			r1, err := client.Get(URL)
			So(err, ShouldBeNil)
			defer r1.Body.Close()

			var b1 bytes.Buffer
			b1.ReadFrom(r1.Body)
			So(server.BodyWriteError, ShouldBeNil)
			So(b1.String(), ShouldEqual, "404")
			So(r1.StatusCode, ShouldEqual, http.StatusNotFound)

			r2, err := client.Get(URL)
			So(err, ShouldBeNil)
			defer r2.Body.Close()

			var b2 bytes.Buffer
			b2.ReadFrom(r2.Body)
			So(server.BodyWriteError, ShouldBeNil)
			So(b2.String(), ShouldEqual, "I'm OK")
			So(r2.StatusCode, ShouldEqual, http.StatusOK)

			// The last value repeats.
			r3, err := client.Get(URL)
			So(err, ShouldBeNil)
			defer r3.Body.Close()

			var b3 bytes.Buffer
			b3.ReadFrom(r3.Body)
			So(server.BodyWriteError, ShouldBeNil)
			So(b3.String(), ShouldEqual, "I'm OK")
			So(r3.StatusCode, ShouldEqual, http.StatusOK)
		})

		Convey("map values", func() {
			server.ResponseBodyMap["/"] = []string{"root1", "root2"}
			server.ResponseStatusMap["/"] = []int{http.StatusOK, http.StatusAccepted}
			server.ResponseBody = []string{"wrong!"}

			path := "/path"

			server.ResponseBodyMap[path] = []string{"path1", "path2"}
			server.ResponseStatusMap[path] = []int{
				http.StatusNotFound, http.StatusCreated}

			// Run each call in closures to close the request body.
			func() {
				r1, err := client.Get(URL)
				So(err, ShouldBeNil)
				defer r1.Body.Close()

				var b1 bytes.Buffer
				b1.ReadFrom(r1.Body)
				So(server.BodyWriteError, ShouldBeNil)
				So(b1.String(), ShouldEqual, "root1")
				So(r1.StatusCode, ShouldEqual, http.StatusOK)
			}()

			func() {
				r2, err := client.Get(URL)
				So(err, ShouldBeNil)
				defer r2.Body.Close()

				var b2 bytes.Buffer
				b2.ReadFrom(r2.Body)
				So(server.BodyWriteError, ShouldBeNil)
				So(b2.String(), ShouldEqual, "root2")
				So(r2.StatusCode, ShouldEqual, http.StatusAccepted)
				So(server.ResponseBodyMap["/"], ShouldResemble, []string{"root2"})
				So(server.ResponseStatusMap["/"], ShouldResemble,
					[]int{http.StatusAccepted})
			}()

			func() {
				r3, err := client.Get(URL + path)
				So(err, ShouldBeNil)
				defer r3.Body.Close()

				var b3 bytes.Buffer
				b3.ReadFrom(r3.Body)
				So(server.BodyWriteError, ShouldBeNil)
				So(b3.String(), ShouldEqual, "path1")
				So(r3.StatusCode, ShouldEqual, http.StatusNotFound)
			}()

			func() {
				r4, err := client.Get(URL + path)
				So(err, ShouldBeNil)
				defer r4.Body.Close()

				var b4 bytes.Buffer
				_, err = b4.ReadFrom(r4.Body)
				So(err, ShouldBeNil)
				So(server.BodyWriteError, ShouldBeNil)
				So(server.BodyWriteBytes, ShouldEqual, len("path2"))
				So(b4.String(), ShouldEqual, "path2")
				So(r4.StatusCode, ShouldEqual, http.StatusCreated)
				So(server.ResponseBodyMap[path], ShouldResemble, []string{"path2"})
				So(server.ResponseStatusMap[path], ShouldResemble,
					[]int{http.StatusCreated})
				So(server.Flushed, ShouldBeTrue)
			}()
		})
	})
}
