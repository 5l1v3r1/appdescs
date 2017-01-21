package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "cmd args...")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Available commands:")
		fmt.Fprintln(os.Stderr, " fetch <play_url> <out_dir>")
		fmt.Fprintln(os.Stderr, " walk <play_url> <out_dir>")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\nThe play_url argument specifies a page of apps, such as")
		fmt.Fprintln(os.Stderr,
			"https://play.google.com/store/apps/collection/topselling_free?hl=en")
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}

	if _, err := os.Stat(os.Args[3]); os.IsNotExist(err) {
		if err := os.Mkdir(os.Args[3], 0755); err != nil {
			fmt.Fprintln(os.Stderr, "Create out dir:", err)
			os.Exit(1)
		}
	}

	if os.Args[1] == "fetch" {
		fetchCommand()
	} else if os.Args[1] == "walk" {
		walkCommand()
	} else {
		fmt.Fprintln(os.Stderr, "Unrecognized command:", os.Args[1])
	}
}

func fetchCommand() {
	playURL := os.Args[2]
	listing, errChan := AppList(playURL)
	dumpListing(listing, errChan)
}

func walkCommand() {
	startID := os.Args[2]
	listing, errChan := RandomWalk(startID)
	dumpListing(listing, errChan)
}

func dumpListing(listing <-chan ListEntry, errChan <-chan error) {
	seenIDs := map[string]bool{}
	for item := range listing {
		seenIDs[item.AppID] = true
		destPath := filepath.Join(os.Args[3], item.AppID+".html")
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
