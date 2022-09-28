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

package gh

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v47/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lister", func() {

	// Test out the ListerConfig struct and its methods
	Context("ListerConfig", func() {
		Describe("GetGithubOrg", func() {
			var (
				options ListerConfig
			)
			BeforeEach(func() {
				options = ListerConfig{}
			})
			It("should return org if given an org/repo formatted string", func() {
				options.Project = "operator-framework/operator-sdk"
				Expect(options.GetGithubOrg()).To(Equal("operator-framework"))
			})
			It("should return empty string if the project is empty", func() {
				options.Project = ""
				Expect(options.GetGithubOrg()).To(Equal(""))
			})
			It("should return entire string if the project has no /", func() {
				options.Project = "operator-framework"
				Expect(options.GetGithubOrg()).To(Equal("operator-framework"))
			})
			It("should return empty string if the project starts with /", func() {
				options.Project = "/operator-framework"
				Expect(options.GetGithubOrg()).To(Equal(""))
			})
		})
		Describe("GetGithubRepo", func() {
			var (
				options ListerConfig
			)
			BeforeEach(func() {
				options = ListerConfig{}
			})
			It("should return repo if given an org/repo formatted string", func() {
				options.Project = "operator-framework/operator-sdk"
				Expect(options.GetGithubRepo()).To(Equal("operator-sdk"))
			})
			It("should return empty string if the project is empty", func() {
				options.Project = ""
				Expect(options.GetGithubRepo()).To(Equal(""))
			})
			It("should return entire string if the project has no /", func() {
				options.Project = "operator-framework"
				Expect(options.GetGithubRepo()).To(Equal("operator-framework"))
			})
			It("should return the second string if the project starts with /", func() {
				options.Project = "/operator-framework"
				Expect(options.GetGithubRepo()).To(Equal("operator-framework"))
			})
		})
		Describe("getToken", func() {
			var (
				options       ListerConfig
				originalToken string
			)
			BeforeEach(func() {
				options = ListerConfig{}
			})
			BeforeEach(func() {
				originalToken = os.Getenv("GITHUB_TOKEN")
				err := os.Setenv("GITHUB_TOKEN", "blah-blah-blah")
				Expect(err).NotTo(HaveOccurred())
			})
			AfterEach(func() {
				err := os.Setenv("GITHUB_TOKEN", originalToken)
				Expect(err).NotTo(HaveOccurred())
			})
			It("should return the token", func() {
				token, err := options.getToken()
				Expect(err).NotTo(HaveOccurred())
				Expect(token).To(Equal("blah-blah-blah"))
			})
			It("should return an error if no token", func() {
				err := os.Unsetenv("GITHUB_TOKEN")
				Expect(err).NotTo(HaveOccurred())

				token, err := options.getToken()
				Expect(err).To(HaveOccurred())
				Expect(token).To(Equal(""))
			})
		})
	})

	Context("With Option methods", func() {
		var (
			options ListerConfig
		)
		BeforeEach(func() {
			options = ListerConfig{}
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
		Describe("WithMilestone", func() {
			It("should set the milestone", func() {
				opt := WithMilestone("47")
				err := opt(&options)
				Expect(err).NotTo(HaveOccurred())
				Expect(options.Milestone).To(Equal("47"))
			})
		})
		Describe("WithAssignee", func() {
			It("should set the assignee", func() {
				opt := WithAssignee("johndoe")
				err := opt(&options)
				Expect(err).NotTo(HaveOccurred())
				Expect(options.Assignee).To(Equal("johndoe"))
			})
		})
		Describe("WithProject", func() {
			It("should set the project", func() {
				opt := WithProject("operator-framework/operator-sdk")
				err := opt(&options)
				Expect(err).NotTo(HaveOccurred())
				Expect(options.Project).To(Equal("operator-framework/operator-sdk"))
			})
		})
		Describe("WithLabel", func() {
			It("should set the label", func() {
				labels := []string{"kind/bug", "documentation"}
				opt := WithLabel(labels)
				err := opt(&options)
				Expect(err).NotTo(HaveOccurred())
				Expect(options.Label).To(Equal(labels))
			})
		})
	})

	Describe("ListIssues", func() {
		var (
			originalToken string
		)
		BeforeEach(func() {
			originalToken = os.Getenv("GITHUB_TOKEN")
			err := os.Setenv("GITHUB_TOKEN", "blah-blah-blah")
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			err := os.Setenv("GITHUB_TOKEN", originalToken)
			Expect(err).NotTo(HaveOccurred())
		})
		It("should return an error if there is no token", func() {
			err := os.Unsetenv("GITHUB_TOKEN")
			Expect(err).NotTo(HaveOccurred())

			iss, err := ListIssues()
			Expect(iss).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
		It("should find open issues", func() {
			mockedHTTPClient := mock.NewMockedHTTPClient(
				mock.WithRequestMatch(mock.GetReposIssuesByOwnerByRepo,
					[]github.Issue{
						{
							ID:    github.Int64(123),
							Title: github.String("Issue 1"),
							State: github.String("open"),
						},
						{
							ID:    github.Int64(456),
							Title: github.String("Issue 2"),
							State: github.String("open"),
						},
					},
				),
			)
			iss, err := ListIssues(WithClient(mockedHTTPClient), WithProject("fakeorg/fakeproject"))
			Expect(iss).NotTo(BeNil())
			Expect(len(iss)).To(Equal(2))
			Expect(err).NotTo(HaveOccurred())
		})
		It("should return error if list fails", func() {
			// if our request returns an error ListIssues should return
			// that error
			mockedHTTPClient := mock.NewMockedHTTPClient(
				mock.WithRequestMatchHandler(
					mock.GetReposIssuesByOwnerByRepo,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						mock.WriteError(
							w,
							http.StatusInternalServerError,
							"github went belly up or something",
						)
					}),
				),
			)
			iss, err := ListIssues(WithClient(mockedHTTPClient), WithProject("fakeorg/fakeproject"))
			Expect(iss).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
		It("should return error if Options return an error", func() {
			_, err := ListIssues(func(c *ListerConfig) error {
				return fmt.Errorf("do you see me")
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("do you see me"))
		})
	})
})
