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
	"net/http"

	"github.com/gorilla/mux"
)

// WithRequestMatchHandler implements a request callback
// for the given `pattern`.
//
// For custom implementations, this handler usage is encouraged.
//
// Example:
//
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
//	)
func WithRequestMatchHandler(
	ep EndpointPattern,
	handler http.Handler,
) MockBackendOption {
	return func(router *mux.Router) {
		router.Handle(ep.Pattern, handler).Methods(ep.Method)
	}
}

// WithRequestMatch implements a simple FIFO for requests
// of the given `pattern`.
//
// Once all responses have been used, it shall panic()!
//
// Example:
//
//	WithRequestMatch(
//		GetUsersByUsername,
//		github.User{
//			Name: github.String("foobar"),
//		},
//	)
func WithRequestMatch(
	ep EndpointPattern,
	responsesFIFO ...interface{},
) MockBackendOption {
	responses := [][]byte{}

	for _, r := range responsesFIFO {
		responses = append(responses, MustMarshal(r))
	}

	return WithRequestMatchHandler(ep, &FIFOReponseHandler{
		Responses: responses,
	})
}
