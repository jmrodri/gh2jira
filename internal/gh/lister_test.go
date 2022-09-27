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
