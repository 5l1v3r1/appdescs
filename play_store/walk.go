package main

import (
	"errors"
	"math/rand"
	"time"
)

func RandomWalk(startPage string) (<-chan ListEntry, <-chan error) {
	rand.Seed(time.Now().UnixNano())
	errChan := make(chan error, 1)
	resChan := make(chan ListEntry, 1)
	go func() {
		defer close(errChan)
		defer close(resChan)
		start, err := readAll(AppList(startPage))
		if err != nil {
			errChan <- err
			return
		}
		seen := map[string]bool{}
		seenList := []string{}
		for _, x := range start {
			seen[x.AppID] = true
			seenList = append(seenList, x.AppID)
		}

		id := seenList[rand.Intn(len(seenList))]
		for {
			// Occasionally, go back to a random app we've seen.
			// This could prevent us from getting stuck in loops.
			if rand.Intn(10) == 0 {
				id = seenList[rand.Intn(len(seenList))]
			}

			similar, err := similarApps(id)
			if err != nil {
				errChan <- err
				return
			}
			if len(similar) == 0 {
				errChan <- errors.New("no similar apps")
				return
			}
			item := similar[rand.Intn(len(similar))]
			resChan <- item
			id = item.AppID

			if !seen[id] {
				seen[id] = true
				seenList = append(seenList, id)
			}
		}
	}()
	return resChan, errChan
}

func similarApps(id string) ([]ListEntry, error) {
	return readAll(AppList("https://play.google.com/store/apps/similar?id=" + id + "&hl=en"))
}

func readAll(list <-chan ListEntry, errCh <-chan error) ([]ListEntry, error) {
	var all []ListEntry
	seen := map[string]bool{}
	for entry := range list {
		if seen[entry.AppID] {
			break
		}
		all = append(all, entry)
	}
	if err := <-errCh; err != nil {
		return nil, err
	}
	return all, nil
}
