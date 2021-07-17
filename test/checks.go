package test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func checkDiff(t *testing.T, buf bytes.Buffer, expect string, formats ...interface{}) {
	expect = fmt.Sprintf(expect, formats...)
	if !cmp.Equal(expect, buf.String()) {
		t.Fatal(
			cmp.Diff(
				buf.String(),
				expect,
			),
		)
	}
}

func checkCreateIssue(t *testing.T, buf bytes.Buffer, endpoint string) string {
	out := strings.Split(buf.String(), " ")
	if len(out) < 3 {
		t.Fatalf("unexpected split count on create output: %v", buf.String())
	}
	issue := out[1]
	expect := fmt.Sprintf("OK %s %s/browse/%s\n", issue, endpoint, issue)
	if !cmp.Equal(expect, buf.String()) {
		t.Fatal(
			cmp.Diff(
				buf.String(),
				expect,
			),
		)
	}
	return issue
}

func checkEditIssue(t *testing.T, buf bytes.Buffer, issue, endpoint string) {
	out := strings.Split(buf.String(), " ")
	if len(out) < 3 {
		t.Fatalf("unexpected split count on create output: %v", buf.String())
	}
	editedIssue := out[1]
	expect := fmt.Sprintf("OK %s %s/browse/%s\n", issue, endpoint, issue)
	if !cmp.Equal(expect, buf.String()) {
		t.Fatal(
			cmp.Diff(
				buf.String(),
				expect,
			),
		)
	}
	if !cmp.Equal(editedIssue, issue) {
		t.Fatal(
			cmp.Diff(
				editedIssue,
				issue,
			),
		)
	}
}

func checkIssueInOutput(t *testing.T, buf bytes.Buffer, issue string) {
	if !strings.Contains(buf.String(), issue) {
		t.Fatalf("issue %s not located in stdout: %s", issue, buf.String())
	}
}

func checkIssueNotInOutput(t *testing.T, buf bytes.Buffer, issue string) {
	if strings.Contains(buf.String(), issue) {
		t.Fatalf("issue %s not located in stdout: %s", issue, buf.String())
	}
}

func checkBlockIssue(t *testing.T, buf bytes.Buffer, issue, blocker, endpoint string) {
	checkDualIssues(t, buf, blocker, issue, endpoint)
}

func checkDupIssue(t *testing.T, buf bytes.Buffer, issue, duplicate, endpoint string) {
	checkDualIssues(t, buf, issue, duplicate, endpoint)
}

func checkDualIssues(t *testing.T, buf bytes.Buffer, first, second, endpoint string) {
	lines := strings.Split(buf.String(), "\n")
	if len(lines) < 2 {
		t.Fatalf("unexpected split count on create output: %v", buf.String())
	}

	testBuf := bytes.NewBuffer([]byte(lines[0] + "\n"))
	checkEditIssue(t, *testBuf, first, endpoint)

	testBuf = bytes.NewBuffer([]byte(lines[1] + "\n"))
	checkEditIssue(t, *testBuf, second, endpoint)
}
