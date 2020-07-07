package sharers

import (
	"context"
	"errors"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubInfo struct {
	owner             string
	repo              string
	pullRequestNumber int
	commit            string
	commentID         int64
}

var reaction string = "+1"

func parsePullRequestUrl(u *url.URL) (*GithubInfo, error) {
	pathRegexp := regexp.MustCompile(`/(?P<owner>.+)/(?P<repo>.+)/pull/(?P<pullRequestNumber>\d+).*`)
	pathResults := FindNamedMatches(pathRegexp, u.Path)
	fragmentRegexp := regexp.MustCompile(`(discussion_)?r(?P<commentID>\d+)`)
	fragmentResults := FindNamedMatches(fragmentRegexp, u.Fragment)

	owner := pathResults["owner"]
	repo := pathResults["repo"]
	pullRequestNumber, err := strconv.Atoi(pathResults["pullRequestNumber"])
	if err != nil {
		return nil, err
	}
	commentID, err := strconv.Atoi(fragmentResults["commentID"])
	if err != nil {
		return nil, err
	}
	if owner == "" || repo == "" || pullRequestNumber == 0 || commentID == 0 {
		return nil, errors.New("github url, but not valid pull request comment")
	}

	return &GithubInfo{
		owner:             owner,
		repo:              repo,
		pullRequestNumber: pullRequestNumber,
		commentID:         int64(commentID),
	}, nil
}

func parseCommitCommentUrl(u *url.URL) (*GithubInfo, error) {
	pathRegexp := regexp.MustCompile(`/(?P<owner>.+)/(?P<repo>.+)/commit/(?P<commit>[a-f0-9]{40}).*`)
	pathResults := FindNamedMatches(pathRegexp, u.Path)
	fragmentRegexp := regexp.MustCompile(`(r|commitcomment-)(?P<commentID>\d+)`)
	fragmentResults := FindNamedMatches(fragmentRegexp, u.Fragment)

	owner := pathResults["owner"]
	repo := pathResults["repo"]
	commit := pathResults["commit"]
	commentID, err := strconv.Atoi(fragmentResults["commentID"])
	if err != nil {
		return nil, err
	}
	if owner == "" || repo == "" || commit == "" || commentID == 0 {
		return nil, errors.New("github url, but not valid pull request comment")
	}

	return &GithubInfo{
		owner:     owner,
		repo:      repo,
		commit:    commit,
		commentID: int64(commentID),
	}, nil
}

func postComment(githubInfo *GithubInfo, content string) (string, error) {
	envVar := "SHARE_TO_CLIPBOARD_URL_GITHUB_ACCESS_TOKEN"
	AccessToken, success := os.LookupEnv(envVar)
	if !success {
		return "", errors.New(envVar + " env var was not set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	var HTMLURL string
	if githubInfo.pullRequestNumber != 0 {
		_, _, err := client.Reactions.CreatePullRequestCommentReaction(ctx, githubInfo.owner, githubInfo.repo, githubInfo.commentID, reaction)
		if err != nil {
			return "", err
		}

		comment, _, err := client.PullRequests.CreateCommentInReplyTo(ctx, githubInfo.owner, githubInfo.repo, githubInfo.pullRequestNumber, content, githubInfo.commentID)
		if err != nil {
			return "", err
		}
		HTMLURL = *comment.HTMLURL
	} else {
		_, _, err := client.Reactions.CreateCommentReaction(ctx, githubInfo.owner, githubInfo.repo, githubInfo.commentID, reaction)
		if err != nil {
			return "", err
		}

		commitComment, _, err := client.Repositories.GetComment(ctx, githubInfo.owner, githubInfo.repo, githubInfo.commentID)
		if err != nil {
			return "", err
		}

		commitComment.Body = &content
		comment, _, err := client.Repositories.CreateComment(ctx, githubInfo.owner, githubInfo.repo, githubInfo.commit, commitComment)
		if err != nil {
			return "", err
		}
		HTMLURL = *comment.HTMLURL
	}

	return HTMLURL, nil
}

func ShareToGithub(u *url.URL, content string) (string, error) {
	hostname := u.Hostname()
	if hostname != "github.com" {
		return "", errors.New("hostname is not github")
	}
	parsers := []func(u *url.URL) (*GithubInfo, error){parseCommitCommentUrl, parsePullRequestUrl}
	var err error
	var githubInfo *GithubInfo
	for _, parser := range parsers {
		githubInfo, err = parser(u)
		if githubInfo != nil {
			break
		}
	}
	if err != nil {
		return "", err
	}
	result, err := postComment(githubInfo, content)
	if err != nil {
		return "", err
	}

	return result, nil
}
