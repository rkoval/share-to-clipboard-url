package main

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/0xAX/notificator"
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/rkoval/share-to-clipboard-url/sharers"
)

func parseText(rawUrl, content string) (bool, error) {
	handlers := []func(*url.URL, string) (string, error){
		sharers.ShareToGithub,
		sharers.ShareToGitlab,
	}
	u, err := url.Parse(rawUrl)
	if err != nil {
		return false, err
	}

	notify := notificator.New(notificator.Options{
		AppName: "share-to-clipboard-url",
	})

	for _, handler := range handlers {
		result, err := handler(u, content)
		if err != nil {
			return false, err
		} else if result == "" {
			continue
		}
		if result != "" {
			fmt.Println(result)
			notify.Push("✅ Success", result, "", notificator.UR_NORMAL)
			// set clipboard to commit we just posted so that we don't accidentally post to a previous comment if we forgot to copy a new one
			err := clipboard.WriteAll(content)
			if err != nil {
				return false, err
			}
			return false, nil
		}
	}
	return true, nil
}

func readClipboard() string {
	clipboardText, err := clipboard.ReadAll()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return clipboardText
}

func main() {
	var content string
	flag.StringVar(&content, "content", "", "content override if nothing passed to stdin")
	var url string
	flag.StringVar(&url, "url", "", "url override if not using clipboard")
	flag.Parse()

	if content == "" {
		rawInput, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		content = string(rawInput)
	}

	if url != "" {
		fmt.Fprintln(os.Stderr, color.BlackString("not reading from clipboard; url override argument provided:\n"), color.BlackString(url))
		_, err := parseText(url, content)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		os.Exit(0)
	}

	clipboardText := readClipboard()
	shouldRetry, err := parseText(clipboardText, content)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	for shouldRetry {
		fmt.Fprintln(os.Stderr, "clipboard did not have a supported url:\n", color.BlackString(clipboardText))
		time.Sleep(1 * time.Second)
		clipboardText := readClipboard()
		shouldRetry, err = parseText(clipboardText, content)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}
}
