package lib

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/go-github/v35/github"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Config for creating an Issue Request.
type Config struct {
	Token string

	Byline bool

	DryRun bool
	Info   bool
	Open   bool
}

var (
	// BigSep separates Issues in the input.
	BigSep = "\n---\n"
	// SmallSep separates metadata in the input.
	SmallSep = "|"
	// ListSep separates list items in metadata.
	ListSep = ","
	// TargetSep separates owner/repo.
	TargetSep = "/"

	// Byline is the suffix to append to Issue bodies.
	Byline = "> 🙌 Bulk-uploaded by https://github.com/hcgatewood/ghissue"

	// IndexWait is how long to wait between creating Issues and opening their
	// URL.
	IndexWait = 3 * time.Second
	// OpenSince is how far back from now to set the search time when opening
	// an Issue search URL.
	OpenSince = 30 * time.Second
)

// Create GitHub Issues.
func Create(cfg *Config, input string) ([]github.IssueRequest, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	target, reqs, err := parse(cfg, input)
	if err != nil {
		return nil, err
	}
	if len(reqs) == 0 {
		return nil, errors.New("found no issues to create")
	}

	if cfg.DryRun {
		fmt.Printf("Would create the following %d issue(s)\n\n", len(reqs))
		printRequests(target, reqs)
		return reqs, nil
	}

	oathClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token}))
	g := github.NewClient(oathClient)
	nums, err := create(g, target, reqs)
	if err != nil {
		return nil, err
	}
	fmt.Println(strings.Join(nums, ","))

	since := time.Now().Add(-OpenSince).UTC().Format(time.RFC3339)
	urlQuery := url.QueryEscape(fmt.Sprintf("is:issue is:open sort:updated-desc created:>%s", since))
	urlFull := fmt.Sprintf("https://github.com/%s/%s/issues?q=%s", target.owner, target.repo, urlQuery)
	if cfg.Info {
		fmt.Println(urlFull)
	}
	if cfg.Open {
		time.Sleep(IndexWait)
		err = browser.OpenURL(urlFull)
		if err != nil {
			panic(err)
		}
	}

	return reqs, nil
}

// TrimInput trims the string input, pre-parsing.
// Callers should call this first on what they pass to Create.
func TrimInput(s string) string {
	// Only trim left. Trailing whitespace will be handled by issue parser.
	return strings.TrimLeftFunc(s, unicode.IsSpace)
}

func create(g *github.Client, target targetT, reqs []github.IssueRequest) ([]string, error) {
	var created []*github.Issue
	for _, req := range reqs {
		c, _, err := g.Issues.Create(context.Background(), target.owner, target.repo, &req)
		if err != nil {
			return nil, errors.Wrap(err, "create issue")
		}
		created = append(created, c)
	}
	nums := getIssueNumbers(created)
	return nums, nil
}

func getIssueNumbers(issues []*github.Issue) []string {
	var nums []string
	for _, issue := range issues {
		numPtr := issue.Number
		if numPtr != nil {
			nums = append(nums, strconv.Itoa(*numPtr))
		}
	}
	return nums
}

func parse(cfg *Config, inp string) (targetT, []github.IssueRequest, error) {
	inpSplit := fields(inp, BigSep)

	if len(inpSplit) < 2 {
		return targetT{}, nil, errors.New("issues file must contain a target and at least one issue")
	}
	targetInp, issuesInp := inpSplit[0], inpSplit[1:]
	if targetInp == "" {
		return targetT{}, nil, errors.New("couldn't find a target repo in input")
	}
	if len(issuesInp) == 0 {
		return targetT{}, nil, errors.New("issues to create was empty in input")
	}

	target, err := parseTarget(targetInp)
	if err != nil {
		return targetT{}, nil, err
	}

	reqs, err := parseIssues(cfg, issuesInp)
	if err != nil {
		return targetT{}, nil, err
	}

	return target, reqs, nil
}

func parseTarget(targetStr string) (targetT, error) {
	split := fields(targetStr, TargetSep)
	if len(split) != 2 {
		return targetT{}, errors.Errorf("parse target: couldn't parse target '%s' to 'owner/repo'", targetStr)
	}

	owner, repo := split[0], split[1]
	if owner == "" {
		return targetT{}, errors.New("parse target: owner can't be empty")
	}
	if repo == "" {
		return targetT{}, errors.New("parse target: repo can't be empty")
	}

	target := targetT{owner: owner, repo: repo}
	return target, nil
}

func parseIssues(cfg *Config, inps []string) ([]github.IssueRequest, error) {
	var reqs []github.IssueRequest
	for _, inp := range inps {
		trimmed := strings.TrimSpace(inp)
		if trimmed == "" {
			continue
		}
		req, err := parseIssue(cfg, trimmed)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}
	return reqs, nil
}

func parseIssue(cfg *Config, inp string) (github.IssueRequest, error) {
	split := strings.SplitN(inp, "\n", 2)
	var metadataInp, bodyInp string
	if len(split) == 1 {
		metadataInp = split[0]
	} else {
		metadataInp, bodyInp = split[0], split[1]
	}

	m := fields(metadataInp, SmallSep)
	if len(m) == 0 {
		return github.IssueRequest{}, errors.Errorf("parse issue: metadata '%s' must contain at least a title", metadataInp)
	}
	var title, labelsInp, assigneesInp string
	title = m[0]
	if len(m) > 1 {
		labelsInp = m[1]
	}
	if len(m) > 2 {
		assigneesInp = m[2]
	}
	if title == "" {
		return github.IssueRequest{}, errors.New("parse issue: title can't be empty")
	}

	issue := github.IssueRequest{
		Title:     &title,
		Body:      getBody(cfg, bodyInp),
		Labels:    parseList(labelsInp),
		Assignees: parseList(assigneesInp),
	}

	return issue, nil
}

func getBody(cfg *Config, inp string) *string {
	trimmed := strings.TrimSpace(inp)
	if trimmed == "" {
		return nil
	}
	if cfg.Byline {
		s := fmt.Sprintf("%s\n\n\n%s", trimmed, Byline)
		return &s
	}
	return &trimmed
}

func parseList(inp string) *[]string {
	split := fields(inp, ListSep)
	if len(split) == 0 {
		return nil
	}
	return &split
}

func fields(s string, sep string) []string {
	ss := strings.Split(s, sep)
	if len(ss) == 1 && ss[0] == "" {
		return nil
	}
	var trimmed []string
	for _, t := range ss {
		trimmed = append(trimmed, strings.TrimSpace(t))
	}
	return trimmed
}

func printRequests(target targetT, rr []github.IssueRequest) {
	ss := []string{target.String()}
	for _, r := range rr {
		ss = append(ss, fmtRequest(r))
	}
	s := strings.Join(ss, BigSep)
	fmt.Println(s)
}

func fmtRequest(r github.IssueRequest) string {
	s := fmt.Sprintf("%s | %v | %v\n%s", r.GetTitle(), r.GetLabels(), r.GetAssignees(), r.GetBody())
	return s
}

type targetT struct {
	owner string
	repo  string
}

func (t *targetT) String() string {
	return fmt.Sprintf("%s/%s", t.owner, t.repo)
}
