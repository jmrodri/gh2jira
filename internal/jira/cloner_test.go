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

	"github.com/migueleliasweb/go-github-mock/src/mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cloner", func() {

	// Test out the ClonerConfig struct and its methods
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

	Context("With Option methods", func() {
		var (
			options ClonerConfig
		)
		BeforeEach(func() {
			options = ClonerConfig{}
		})
		Describe("WithClient", func() {
			It("should set client to nil if passed nil", func() {
				opt := WithClient(nil)
				err := opt(&options)
				Expect(err).NotTo(HaveOccurred())
				Expect(options.client).To(BeNil())
			})
			It("should set the client if given one", func() {
				mc := mock.NewMockedHTTPClient()
				opt := WithClient(mc)
				err := opt(&options)
				Expect(err).NotTo(HaveOccurred())
				Expect(options.client).To(Equal(mc))
			})
		})
		Describe("WithDryRun", func() {
			It("should set the dryrun to false", func() {
				opt := WithDryRun(false)
				err := opt(&options)
				Expect(err).NotTo(HaveOccurred())
				Expect(options.dryRun).To(BeFalse())
			})
		})
		Describe("WithProject", func() {
			It("should set the project", func() {
				opt := WithProject("OSDK")
				err := opt(&options)
				Expect(err).NotTo(HaveOccurred())
				Expect(options.project).To(Equal("OSDK"))
			})
		})
		Describe("WithJiraURL", func() {
			It("should set the jira url", func() {
				url := "https://issues.jira.com"
				opt := WithJiraURL(url)
				err := opt(&options)
				Expect(err).NotTo(HaveOccurred())
				Expect(options.jiraURL).To(Equal(url))
			})
		})
	})

	Describe("getWebURL", func() {
		It("should convert the Github API URL to web URL", func() {
			apiurl := "https://api.github.com/repos/operator-framework/operator-sdk/issues/3447"
			expectedurl := "https://github.com/operator-framework/operator-sdk/issues/3447"
			Expect(getWebURL(apiurl)).To(Equal(expectedurl))
		})
		It("should leave url untouched if it is blank", func() {
			apiurl := ""
			expectedurl := ""
			Expect(getWebURL(apiurl)).To(Equal(expectedurl))
		})
		It("should leave url untouched if it does not have any matching strings", func() {
			apiurl := "http://www.google.com"
			expectedurl := "http://www.google.com"
			Expect(getWebURL(apiurl)).To(Equal(expectedurl))

			apiurl = "http://example.com"
			expectedurl = "http://example.com"
			Expect(getWebURL(apiurl)).To(Equal(expectedurl))

			apiurl = "http://github.com/operator-framework"
			expectedurl = "http://github.com/operator-framework"
			Expect(getWebURL(apiurl)).To(Equal(expectedurl))
		})
	})

	Describe("Clone", func() {
		var (
			originalToken string
		)
		BeforeEach(func() {
			originalToken = os.Getenv("JIRA_TOKEN")
			err := os.Setenv("JIRA_TOKEN", "blah-blah-blah")
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			err := os.Setenv("JIRA_TOKEN", originalToken)
			Expect(err).NotTo(HaveOccurred())
		})
		It("shoud return an error if there is no token", func() {
			err := os.Unsetenv("JIRA_TOKEN")
			Expect(err).NotTo(HaveOccurred())

			err = Clone(nil)
			Expect(err).To(HaveOccurred())
		})
		It("should print out issue when dryRun is true", func() {
		})
		It("should save issue in jira when dryRun is false", func() {
		})
		It("should return an error if jira client returns an error", func() {
		})
	})
})
