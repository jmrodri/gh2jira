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
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	gojira "github.com/andygrunwald/go-jira"
	"github.com/google/go-github/v47/github"
	jmock "github.com/jmrodri/gh2jira/internal/jira/mock"

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
				mc := jmock.NewMockedHTTPClient()
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

			_, err = Clone(nil)
			Expect(err).To(HaveOccurred())
		})
		It("should print out issue when dryRun is true", func() {
			mockedHTTPClient := jmock.NewMockedHTTPClient(
				jmock.WithRequestMatch(jmock.PostIssue),
			)

			// giving it a github issue
			ghissue := &github.Issue{
				ID:     github.Int64(12213123),
				Number: github.Int(123),
				Title:  github.String("Issue 1"),
				State:  github.String("open"),
				Body:   github.String("body of the issue"),
				URL:    github.String("https://api.github.com/repos/foo/bar/issues/123"),
			}

			// Capture stdout to verify the printing
			r, w, _ := os.Pipe()
			tmp := os.Stdout
			defer func() {
				os.Stdout = tmp
			}()
			os.Stdout = w
			go func() {
				// Test the clone function
				_, err := Clone(ghissue, WithClient(mockedHTTPClient),
					WithDryRun(true),
					WithJiraURL("http://localhost"),
				)
				w.Close()
				Expect(err).NotTo(HaveOccurred())
			}()
			stdout, _ := io.ReadAll(r)

			Expect(strings.Contains(string(stdout), "DRY RUN MODE")).To(BeTrue())
		})
		It("should return an error if jira client returns an error", func() {
			// if our request returns an error ListIssues should return
			// that error
			mockedHTTPClient := jmock.NewMockedHTTPClient(
				jmock.WithRequestMatchHandler(
					jmock.PostIssue,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						jmock.WriteError(
							w,
							http.StatusInternalServerError,
							"jira went belly up or something",
						)
					}),
				),
			)
			_, err := Clone(nil, WithClient(mockedHTTPClient),
				WithDryRun(false),
				WithJiraURL("http://localhost"),
			)
			Expect(err).To(HaveOccurred())
		})
		It("should return error if Options return an error", func() {
			_, err := Clone(nil, func(c *ClonerConfig) error {
				return fmt.Errorf("do you see me")
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("do you see me"))
		})
		It("should create a jira issue from the given github issue", func() {
			// expected return jira issue
			expectedissue := gojira.Issue{
				Fields: &gojira.IssueFields{
					Description: "body of the issue\n\nUpstream Github issue: " +
						"https://github.com/foo/bar/issues/123\n",
					Type: gojira.IssueType{
						Name: "Story",
					},
					Project: gojira.Project{
						Key: "OSDK",
					},
					Summary: "[UPSTREAM] Issue 1 #123",
				},
			}

			mockedHTTPClient := jmock.NewMockedHTTPClient(
				jmock.WithRequestMatch(jmock.PostIssue, expectedissue),
			)

			// giving it a github issue
			ghissue := &github.Issue{
				ID:     github.Int64(12213123),
				Number: github.Int(123),
				Title:  github.String("Issue 1"),
				State:  github.String("open"),
				Body:   github.String("body of the issue"),
				URL:    github.String("https://api.github.com/repos/foo/bar/issues/123"),
			}

			// Test the clone function
			jissue, err := Clone(ghissue, WithClient(mockedHTTPClient),
				WithDryRun(false),
				WithJiraURL("http://localhost"),
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(jissue).NotTo(BeNil())
			Expect(jissue.Fields.Description).To(Equal(expectedissue.Fields.Description))
			Expect(jissue.Fields.Type).To(Equal(expectedissue.Fields.Type))
			Expect(jissue.Fields.Project).To(Equal(expectedissue.Fields.Project))
			Expect(jissue.Fields.Summary).To(Equal(expectedissue.Fields.Summary))
		})
	})
})
