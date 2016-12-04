package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "play_url out_dir")
		fmt.Fprintln(os.Stderr, "\nThe play_url argument specifies a page of apps, such as")
		fmt.Fprintln(os.Stderr,
			"https://play.google.com/store/apps/collection/topselling_free?hl=en")
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}

	if _, err := os.Stat(os.Args[2]); os.IsNotExist(err) {
		if err := os.Mkdir(os.Args[2], 0755); err != nil {
			fmt.Fprintln(os.Stderr, "Create out dir:", err)
			os.Exit(1)
		}
	}

	playURL := os.Args[1]
	listing, errChan := AppList(playURL)
	for item := range listing {
		destPath := filepath.Join(os.Args[2], item.AppID+".html")
		if _, err := os.Stat(destPath); err == nil {
			log.Println("Already have", item.Name, "("+item.AppID+")")
			continue
		}
		log.Println("Fetching", item.Name, "("+item.AppID+")")
		desc, err := AppDescription(item.AppID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Fetch description:", err)
			os.Exit(1)
		}
		if err := ioutil.WriteFile(destPath, []byte(desc), 0755); err != nil {
			fmt.Fprintln(os.Stderr, "Write description:", err)
			os.Exit(1)
		}
	}
	if err := <-errChan; err != nil {
		fmt.Fprintln(os.Stderr, "Fetch list:", err)
		os.Exit(1)
	}
}
