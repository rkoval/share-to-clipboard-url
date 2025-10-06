package sharers

import (
	"fmt"
	"net/url"
	"testing"
)

var gitlabOrg string = "org"
var gitlabProject string = "project"
var mergeRequestIID int = 1
var gitlabCommit string = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
var noteID int = 11111111

func TestParseMergeRequestUrl(t *testing.T) {
	rawUrls := []string{
		fmt.Sprintf("https://gitlab.com/%s/%s/-/merge_requests/%d#note_%d", gitlabOrg, gitlabProject, mergeRequestIID, noteID),
		fmt.Sprintf("https://gitlab.com/%s/%s/-/merge_requests/%d/diffs#note_%d", gitlabOrg, gitlabProject, mergeRequestIID, noteID),
	}
	for _, rawUrl := range rawUrls {
		u, err := url.Parse(rawUrl)
		if err != nil {
			t.Errorf("%s was not a url", rawUrl)
		}
		t.Run(rawUrl, func(t *testing.T) {
			gitlabInfo, err := parseMergeRequestUrl(u)
			if err != nil {
				t.Errorf("couldn't parse gitlab info from url; got %s", rawUrl)
			} else if gitlabInfo.owner != gitlabOrg {
				t.Errorf("owner was not correct; got %s", gitlabInfo.owner)
			} else if gitlabInfo.project != gitlabProject {
				t.Errorf("project was not correct; got %s", gitlabInfo.project)
			} else if gitlabInfo.mergeRequestIID != mergeRequestIID {
				t.Errorf("mergeRequestIID was not correct; got %d", gitlabInfo.mergeRequestIID)
			} else if gitlabInfo.commit != "" {
				t.Errorf("commit was not correct; got %s", gitlabInfo.commit)
			} else if gitlabInfo.noteID != noteID {
				t.Errorf("noteID was not correct; got %d", gitlabInfo.noteID)
			}
		})
	}
}

func TestParseGitlabCommitCommentUrl(t *testing.T) {
	rawUrls := []string{
		fmt.Sprintf("https://gitlab.com/%s/%s/-/commit/%s#note_%d", gitlabOrg, gitlabProject, gitlabCommit, noteID),
	}
	for _, rawUrl := range rawUrls {
		u, err := url.Parse(rawUrl)
		if err != nil {
			t.Errorf("%s was not a url", rawUrl)
		}
		t.Run(rawUrl, func(t *testing.T) {
			gitlabInfo, err := parseGitlabCommitCommentUrl(u)
			if err != nil {
				t.Errorf("couldn't parse gitlab info from url; got %s", rawUrl)
			} else if gitlabInfo.owner != gitlabOrg {
				t.Errorf("owner was not correct; got %s", gitlabInfo.owner)
			} else if gitlabInfo.project != gitlabProject {
				t.Errorf("project was not correct; got %s", gitlabInfo.project)
			} else if gitlabInfo.commit != gitlabCommit {
				t.Errorf("commit was not correct; got %s", gitlabInfo.commit)
			} else if gitlabInfo.mergeRequestIID != 0 {
				t.Errorf("mergeRequestIID was not correct; got %d", gitlabInfo.mergeRequestIID)
			} else if gitlabInfo.noteID != noteID {
				t.Errorf("noteID was not correct; got %d", gitlabInfo.noteID)
			}
		})
	}
}
