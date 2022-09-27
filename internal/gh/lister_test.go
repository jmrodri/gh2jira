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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lister", func() {

	Describe("GetGithubOrg", func() {
		var (
			options ListerOptions
		)
		BeforeEach(func() {
			options = ListerOptions{}
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
			options ListerOptions
		)
		BeforeEach(func() {
			options = ListerOptions{}
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
})
