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

package jira

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lister", func() {

	// Test out the ListerConfig struct and its methods
	Context("ClonerConfig", func() {
		Describe("getToken", func() {
			var (
				options       ClonerConfig
				originalToken string
			)
			BeforeEach(func() {
				options = ClonerConfig{}
			})
			BeforeEach(func() {
				originalToken = os.Getenv("JIRA_TOKEN")
				err := os.Setenv("JIRA_TOKEN", "blah-blah-blah")
				Expect(err).NotTo(HaveOccurred())
			})
			AfterEach(func() {
				err := os.Setenv("JIRA_TOKEN", originalToken)
				Expect(err).NotTo(HaveOccurred())
			})
			It("should return the token", func() {
				token, err := options.getToken()
				Expect(err).NotTo(HaveOccurred())
				Expect(token).To(Equal("blah-blah-blah"))
			})
			It("should return an error if no token", func() {
				err := os.Unsetenv("JIRA_TOKEN")
				Expect(err).NotTo(HaveOccurred())

				token, err := options.getToken()
				Expect(err).To(HaveOccurred())
				Expect(token).To(Equal(""))
			})
		})
	})

	Describe("Clone", func() {
	})
})
