package sharers

import (
	"fmt"
	"net/url"
	"testing"
)

var org string = "org"
var repo string = "repo"
var pullRequestNumber int = 1
var commit string = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
var commentID int64 = 11111111

func TestParsePullRequestUrl(t *testing.T) {
	rawUrls := []string{
		fmt.Sprintf("https://github.com/%s/%s/pull/%d#discussion_r%d", org, repo, pullRequestNumber, commentID),
		fmt.Sprintf("https://github.com/%s/%s/pull/%d/files#r%d", org, repo, pullRequestNumber, commentID),
	}
	for _, rawUrl := range rawUrls {
		u, err := url.Parse(rawUrl)
		if err != nil {
			t.Errorf("%s was not a url", rawUrl)
		}
		t.Run(rawUrl, func(t *testing.T) {
			githubInfo, err := parsePullRequestUrl(u)
			if err != nil {
				t.Errorf("couldn't parse github info from url; got %s", rawUrl)
			} else if githubInfo.owner != org {
				t.Errorf("owner was not correct; got %s", githubInfo.owner)
			} else if githubInfo.repo != repo {
				t.Errorf("repo was not correct; got %s", githubInfo.repo)
			} else if githubInfo.pullRequestNumber != pullRequestNumber {
				t.Errorf("pullRequestNumber was not correct; got %d", githubInfo.pullRequestNumber)
			} else if githubInfo.commit != "" {
				t.Errorf("commit was not correct; got %s", githubInfo.commit)
			} else if githubInfo.commentID != commentID {
				t.Errorf("commentID was not correct; got %d", githubInfo.commentID)
			}
		})
	}
}

func TestParseCommitCommentUrl(t *testing.T) {
	rawUrls := []string{
		fmt.Sprintf("https://github.com/%s/%s/commit/%s#r%d", org, repo, commit, commentID),
		fmt.Sprintf("https://github.com/%s/%s/commit/%s#commitcomment-%d", org, repo, commit, commentID),
	}
	for _, rawUrl := range rawUrls {
		u, err := url.Parse(rawUrl)
		if err != nil {
			t.Errorf("%s was not a url", rawUrl)
		}
		t.Run(rawUrl, func(t *testing.T) {
			githubInfo, err := parseCommitCommentUrl(u)
			if err != nil {
				t.Errorf("couldn't parse github info from url; got %s", rawUrl)
			} else if githubInfo.owner != org {
				t.Errorf("owner was not correct; got %s", githubInfo.owner)
			} else if githubInfo.repo != repo {
				t.Errorf("repo was not correct; got %s", githubInfo.repo)
			} else if githubInfo.commit != commit {
				t.Errorf("commit was not correct; got %s", githubInfo.commit)
			} else if githubInfo.pullRequestNumber != 0 {
				t.Errorf("pullRequestNumber was not correct; got %d", githubInfo.pullRequestNumber)
			} else if githubInfo.commentID != commentID {
				t.Errorf("commentID was not correct; got %d", githubInfo.commentID)
			}
		})
	}
}
