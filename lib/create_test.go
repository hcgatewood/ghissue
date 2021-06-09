package lib_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-github/v35/github"
	"github.com/stretchr/testify/assert"

	"hcgatewood/ghissue/lib"
)

func TestCreate(t *testing.T) {
	tests := []testT{
		{
			input:     "err_empty",
			expectErr: true,
		},
		{
			input:     "err_longowner",
			expectErr: true,
		},
		{
			input:     "err_noissues",
			expectErr: true,
		},
		{
			input:     "err_noowner",
			expectErr: true,
		},
		{
			input:     "err_norepo",
			expectErr: true,
		},
		{
			input:     "err_shortowner",
			expectErr: true,
		},
		{
			input: "simple",
			expected: []github.IssueRequest{
				{
					Title: github.String("Simple title"),
				},
			},
		},
		{
			input: "single",
			expected: []github.IssueRequest{
				{
					Title:     github.String("Single title"),
					Labels:    &[]string{"label0", "label1"},
					Assignees: &[]string{"hcgatewood23"},
					Body:      github.String("Long body\nWith multiple lines\nAnd ending with the issue sep"),
				},
			},
		},
		{
			input: "multiple",
			expected: []github.IssueRequest{
				{
					Title:     github.String("Multiple title A"),
					Labels:    &[]string{"label0", "label1"},
					Assignees: &[]string{"hcgatewood23"},
					Body: github.String(
						"Long body\nWith multiple lines\nAnd additional info\nIncluding some triple hyphens --- " +
							"that aren't at beginning of the line\n--- And hyphens that aren't newline-anchored\n--\n" +
							"And hyphens that are newline-anchored but only 2\nAlso add some unicode for fun 👉🌟",
					),
				},
				{
					Title:  github.String("Multiple title B"),
					Labels: &[]string{"label2"},
					Body:   github.String("Short body"),
				},
				{
					Title: github.String("Multiple title C"),
				},
			},
		},
	}

	cfg := &lib.Config{DryRun: true}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			actual, err := lib.Create(cfg, test.readInput(t))
			if test.expectErr && err == nil {
				assert.Fail(t, "expected an error from Create but received nil")
			} else if !test.expectErr && err != nil {
				assert.Failf(t, "expected no error from Create but received non-nil", "%+v", err)
			} else {
				assert.Equal(t, test.expected, actual)
			}
		})
	}
}

type testT struct {
	input     string
	expected  []github.IssueRequest
	expectErr bool
}

func (s *testT) readInput(t *testing.T) string {
	bytes, err := ioutil.ReadFile(fmt.Sprintf("../testdata/%s.txt", s.input))
	assert.NoError(t, err)
	return lib.TrimInput(string(bytes))
}
