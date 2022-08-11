# share-to-clipboard-url

A simple CLI tool which will attempt to parse a URL from your clipboard and then post whatever is read from stdin as a reply to the parsed clipboard link.

To install, run:

```sh
# go must be installed
go get github.com/rkoval/share-to-clipboard-url
```

Example usage [here](https://github.com/rkoval/dotfiles/blob/69a6149b3992d29e1b3f3b1fbee58af708b82928/source/50_misc.sh#L70-L93)

For a list of supported sharers, see [this directory](./sharers)

- To use GitHub, [create a personal access token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token) with the `repo` scope and assign it to the environment variable `SHARE_TO_CLIPBOARD_URL_GITHUB_ACCESS_TOKEN` before running this program
