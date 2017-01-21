package main

import (
	"errors"
	"math/rand"
)

func RandomWalk(startID string) (<-chan ListEntry, <-chan error) {
	errChan := make(chan error, 1)
	resChan := make(chan ListEntry, 1)
	go func() {
		defer close(errChan)
		defer close(resChan)
		id := startID
		for {
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
		}
	}()
	return resChan, errChan
}

func similarApps(id string) ([]ListEntry, error) {
	list, errCh := AppList("https://play.google.com/store/apps/similar?id=" + id + "&hl=en")
	var all []ListEntry
	for entry := range list {
		all = append(all, entry)
	}
	if err := <-errCh; err != nil {
		return nil, err
	}
	return all, nil
}
