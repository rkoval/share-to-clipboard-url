# share-to-clipboard-link

Simple CLI tool which will attempt to parse a URL from your clipboard and then post whatever is read from stdin to as a reply to the parsed clipboard link.

For a list of supported sharers, see [this directory](./sharers)

- To use GitHub, [create a personal access token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token) with the `repo` scope and assign it to the environment variable `SHARE_TO_CLIPBOARD_LINK_GITHUB_ACCESS_TOKEN` before running this program