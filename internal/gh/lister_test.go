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
	"os"

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

			err = ListIssues()
			Expect(err).To(HaveOccurred())
		})
		// It("should initialize ListerOption if not set", func() {
		//     // force ListIssue to return early
		//     err := os.Unsetenv("GITHUB_TOKEN")
		//     Expect(err).NotTo(HaveOccurred())
		//
		//     Expect(lister.Options).To(BeNil())
		//     err = lister.ListIssues()
		//     Expect(err).To(HaveOccurred())
		//     Expect(lister.Options).NotTo(BeNil())
		// })
		It("should not return an error", func() {
			Skip("figure out how to use the mock go github library")
			mockedHTTPClient := mock.NewMockedHTTPClient()
			err := ListIssues(WithClient(mockedHTTPClient))
			Expect(err).To(HaveOccurred())
		})
	})
})
