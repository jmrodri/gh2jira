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
	"io"
	"os"

	"github.com/google/go-github/v47/github"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var expectedLong = `Issue:	123
State:	open

   Issue 1

`
var _ = Describe("Printer", func() {

	Describe("PrintGithubIssue", func() {
		var (
			issue     *github.Issue
			r, w, tmp *os.File
		)
		BeforeEach(func() {
			issue = &github.Issue{
				ID:     github.Int64(12213123),
				Number: github.Int(123),
				Title:  github.String("Issue 1"),
				State:  github.String("open"),
				Body:   github.String("body of the issue"),
				URL:    github.String("https://api.github.com/repos/foo/bar/issues/123"),
			}

			r, w, _ = os.Pipe()
			tmp = os.Stdout
			os.Stdout = w
		})
		AfterEach(func() {
			os.Stdout = tmp
		})
		It("should print issue in color & one line", func() {
			go func() {
				PrintGithubIssue(issue, true, true)
				w.Close()
			}()

			stdout, _ := io.ReadAll(r)

			expected := fmt.Sprintf("\033[33m%5d\033[0m \033[32m%s\033[0m %s\n",
				issue.GetNumber(), issue.GetState(), issue.GetTitle())

			Expect(expected).To(Equal(string(stdout)))
		})
		It("should ignore color flag if not printing in oneline", func() {
			go func() {
				PrintGithubIssue(issue, false, true)
				w.Close()
			}()

			stdout, _ := io.ReadAll(r)

			Expect(expectedLong).To(Equal(string(stdout)))
		})
		It("should print in oneline but not in color", func() {
			go func() {
				PrintGithubIssue(issue, true, false)
				w.Close()
			}()

			stdout, _ := io.ReadAll(r)

			expected := fmt.Sprintf("%5d %s %s\n", issue.GetNumber(),
				issue.GetState(), issue.GetTitle())

			Expect(expected).To(Equal(string(stdout)))
		})
		It("should print full output not in color", func() {
			go func() {
				PrintGithubIssue(issue, false, false)
				w.Close()
			}()

			stdout, _ := io.ReadAll(r)

			Expect(expectedLong).To(Equal(string(stdout)))
		})
	})
})
