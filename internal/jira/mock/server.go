// Copyright © 2022 jesus m. rodriguez jmrodri@gmail.com
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

package mock

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

// EndpointPattern models the GitHub's API endpoints
type EndpointPattern struct {
	Pattern string // eg. "/repos/{owner}/{repo}/actions/artifacts"
	Method  string // "GET", "POST", "PATCH", etc
}

// MockBackendOption is used to configure the *mux.router
// for the mocked backend
type MockBackendOption func(*mux.Router)

// FIFOReponseHandler handler implementation that
// responds to the HTTP requests following a FIFO approach.
//
// Once all available `Responses` have been used, this handler will panic()!
type FIFOReponseHandler struct {
	lock         sync.Mutex
	Responses    [][]byte
	CurrentIndex int
}

// ServeHTTP implementation of `http.Handler`
func (srh *FIFOReponseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	srh.lock.Lock()
	defer srh.lock.Unlock()
	if srh.CurrentIndex > len(srh.Responses) {
		panic(fmt.Sprintf(
			"go-github-mock: no more mocks available for %s",
			r.URL.Path,
		))
	}

	defer func() {
		srh.CurrentIndex++
	}()

	w.Write(srh.Responses[srh.CurrentIndex])
}

// EnforceHostRoundTripper rewrites all requests with the given `Host`.
type EnforceHostRoundTripper struct {
	Host                 string
	UpstreamRoundTripper http.RoundTripper
}

// RoundTrip implementation of `http.RoundTripper`
func (efrt *EnforceHostRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	splitHost := strings.Split(efrt.Host, "://")
	r.URL.Scheme = splitHost[0]
	r.URL.Host = splitHost[1]

	return efrt.UpstreamRoundTripper.RoundTrip(r)
}

// NewMockedHTTPClient creates and configures an http.Client that points to
// a mocked GitHub's backend API.
//
// Example:
//
// mockedHTTPClient := NewMockedHTTPClient(
//
//	WithRequestMatch(
//		GetUsersByUsername,
//		github.User{
//			Name: github.String("foobar"),
//		},
//	),
//	WithRequestMatch(
//		GetUsersOrgsByUsername,
//		[]github.Organization{
//			{
//				Name: github.String("foobar123thisorgwasmocked"),
//			},
//		},
//	),
//	WithRequestMatchHandler(
//		GetOrgsProjectsByOrg,
//		func(w http.ResponseWriter, _ *http.Request) {
//			w.Write(MustMarshal([]github.Project{
//				{
//					Name: github.String("mocked-proj-1"),
//				},
//				{
//					Name: github.String("mocked-proj-2"),
//				},
//			}))
//		},
//	),
//
// )
//
// c := github.NewClient(mockedHTTPClient)
func NewMockedHTTPClient(options ...MockBackendOption) *http.Client {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteError(
			w,
			http.StatusNotFound,
			fmt.Sprintf("mock response not found for %s", r.URL.Path),
		)
	})

	for _, o := range options {
		o(router)
	}

	mockServer := httptest.NewServer(router)

	c := mockServer.Client()

	c.Transport = &EnforceHostRoundTripper{
		Host:                 mockServer.URL,
		UpstreamRoundTripper: mockServer.Client().Transport,
	}

	return c
}
