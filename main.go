package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"github.com/0xAX/notificator"
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/rkoval/share-to-clipboard-url/sharers"
)

func parseText(clipboardText, content string) error {
	handlers := []func(*url.URL, string) (string, error){sharers.ShareToGithub}
	u, err := url.Parse(clipboardText)
	if err != nil {
		return err
	}

	notify := notificator.New(notificator.Options{
		AppName: "share-to-clipboard-url",
	})

	for _, handler := range handlers {
		result, err := handler(u, content)
		if err != nil {
			notify.Push("❌ Error", err.Error(), "", notificator.UR_CRITICAL)
			return err
		}
		if result != "" {
			fmt.Println(result)
			notify.Push("✅ Success", result, "", notificator.UR_NORMAL)
			// set clipboard to commit we just posted so that we don't accidentally post to a previous comment if we forgot to copy a new one
			err := clipboard.WriteAll(content)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
	rawInput, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	content := string(rawInput)

	clipboardText := readClipboard()
	err = parseText(clipboardText, content)
	for err != nil {
		fmt.Fprintln(os.Stderr, "clipboard did not have a supported url:\n", color.BlackString(clipboardText))
		time.Sleep(1 * time.Second)
		clipboardText := readClipboard()
		err = parseText(clipboardText, content)
	}
}
