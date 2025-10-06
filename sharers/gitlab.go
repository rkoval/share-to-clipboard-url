package sharers

import (
	"context"
	"errors"
	"net/url"
	"os"
	"regexp"
	"strconv"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabInfo struct {
	owner           string
	project         string
	mergeRequestIID int
	commit          string
	noteID          int
}

func parseMergeRequestUrl(u *url.URL) (*GitlabInfo, error) {
	pathRegexp := regexp.MustCompile(`/(?P<owner>.+)/(?P<project>.+)/-/merge_requests/(?P<mergeRequestIID>\d+).*`)
	pathResults := FindNamedMatches(pathRegexp, u.Path)
	fragmentRegexp := regexp.MustCompile(`note_(?P<noteID>\d+)`)
	fragmentResults := FindNamedMatches(fragmentRegexp, u.Fragment)

	owner := pathResults["owner"]
	project := pathResults["project"]
	mergeRequestIID, err := strconv.Atoi(pathResults["mergeRequestIID"])
	if err != nil {
		return nil, err
	}
	noteID, err := strconv.Atoi(fragmentResults["noteID"])
	if err != nil {
		return nil, err
	}
	if owner == "" || project == "" || mergeRequestIID == 0 || noteID == 0 {
		return nil, errors.New("gitlab url, but not valid merge request comment")
	}

	return &GitlabInfo{
		owner:           owner,
		project:         project,
		mergeRequestIID: mergeRequestIID,
		noteID:          noteID,
	}, nil
}

func parseGitlabCommitCommentUrl(u *url.URL) (*GitlabInfo, error) {
	pathRegexp := regexp.MustCompile(`/(?P<owner>.+)/(?P<project>.+)/-/commit/(?P<commit>[a-f0-9]{40}).*`)
	pathResults := FindNamedMatches(pathRegexp, u.Path)
	fragmentRegexp := regexp.MustCompile(`note_(?P<noteID>\d+)`)
	fragmentResults := FindNamedMatches(fragmentRegexp, u.Fragment)

	owner := pathResults["owner"]
	project := pathResults["project"]
	commit := pathResults["commit"]
	noteID, err := strconv.Atoi(fragmentResults["noteID"])
	if err != nil {
		return nil, err
	}
	if owner == "" || project == "" || commit == "" || noteID == 0 {
		return nil, errors.New("gitlab url, but not valid commit comment")
	}

	return &GitlabInfo{
		owner:   owner,
		project: project,
		commit:  commit,
		noteID:  noteID,
	}, nil
}

func postGitlabComment(gitlabInfo *GitlabInfo, content string) (string, error) {
	envVar := "SHARE_TO_CLIPBOARD_URL_GITLAB_ACCESS_TOKEN"
	accessToken, success := os.LookupEnv(envVar)
	if !success {
		return "", errors.New(envVar + " env var was not set")
	}

	client, err := gitlab.NewClient(accessToken)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	projectID := gitlabInfo.owner + "/" + gitlabInfo.project

	var noteURL string
	if gitlabInfo.mergeRequestIID != 0 {
		// Create a reply note
		note, _, err := client.Notes.CreateMergeRequestNote(
			projectID,
			gitlabInfo.mergeRequestIID,
			&gitlab.CreateMergeRequestNoteOptions{Body: &content},
			gitlab.WithContext(ctx),
		)
		if err != nil {
			return "", err
		}
		// Get MR to construct the URL
		mr, _, err := client.MergeRequests.GetMergeRequest(projectID, gitlabInfo.mergeRequestIID, nil, gitlab.WithContext(ctx))
		if err != nil {
			return "", err
		}
		noteURL = mr.WebURL + "#note_" + strconv.Itoa(note.ID)
	} else {
		// Create a reply note on the commit
		_, _, err := client.Commits.PostCommitComment(
			projectID,
			gitlabInfo.commit,
			&gitlab.PostCommitCommentOptions{Note: &content},
			gitlab.WithContext(ctx),
		)
		if err != nil {
			return "", err
		}
		// Get the project to construct the URL
		project, _, err := client.Projects.GetProject(projectID, nil, gitlab.WithContext(ctx))
		if err != nil {
			return "", err
		}
		// Construct URL manually - commit comments don't return a direct URL
		noteURL = project.WebURL + "/-/commit/" + gitlabInfo.commit
	}

	return noteURL, nil
}

func ShareToGitlab(u *url.URL, content string) (string, error) {
	hostname := u.Hostname()
	extraHost := os.Getenv("SHARE_TO_CLIPBOARD_URL_EXTRA_GITLAB_HOST")

	isValidHost := hostname == "gitlab.com"
	if !isValidHost && extraHost != "" {
		isValidHost = hostname == extraHost
	}

	if !isValidHost {
		return "", nil
	}

	parsers := []func(u *url.URL) (*GitlabInfo, error){parseGitlabCommitCommentUrl, parseMergeRequestUrl}
	var err error
	var gitlabInfo *GitlabInfo
	for _, parser := range parsers {
		gitlabInfo, err = parser(u)
		if gitlabInfo != nil {
			break
		}
	}
	if err != nil {
		return "", err
	}
	result, err := postGitlabComment(gitlabInfo, content)
	if err != nil {
		return "", err
	}

	return result, nil
}
