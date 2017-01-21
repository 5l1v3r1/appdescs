package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/yhat/scrape"

	"golang.org/x/net/html"
)

const pageListCount = 60

type ListEntry struct {
	Name  string
	AppID string
}

func AppList(page string) (<-chan ListEntry, <-chan error) {
	page = ensureHasLanguage(page)

	resList := make(chan ListEntry)
	resErr := make(chan error, 1)
	go func() {
		defer close(resList)
		defer close(resErr)
		var idx int
		seenIDs := map[string]bool{}
		for {
			listing, err := loadListing(idx, page)
			if err != nil {
				resErr <- err
				return
			}
			if len(listing) == 0 {
				return
			}
			for _, x := range listing {
				if seenIDs[x.AppID] {
					return
				}
				seenIDs[x.AppID] = true
				resList <- x
			}
			idx += len(listing)
		}
	}()
	return resList, resErr
}

func loadListing(idx int, page string) ([]ListEntry, error) {
	req := url.Values{}
	req.Set("cctcss", "square-cover")
	req.Set("cllayout", "NORMAL")
	req.Set("hl", "en")
	req.Set("ipf", "1")
	req.Set("num", strconv.Itoa(pageListCount))
	req.Set("numChildren", "0")
	req.Set("start", strconv.Itoa(idx))
	req.Set("xhr", "1")
	newURL := page + "&authuser=0"
	resp, err := http.PostForm(newURL, req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.New("request list page: " + err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read list page: " + err.Error())
	}

	expr := regexp.MustCompilePOSIX(`\<a class="card-click-target" href="/store/apps/details` +
		`\?id=([^"]*)" aria-label=" *([0-9]*\.)? *([^"]*) *"\>`)
	matches := expr.FindAllStringSubmatch(string(body), -1)
	var res []ListEntry
	for _, match := range matches {
		res = append(res, ListEntry{Name: unescapeHTML(match[3]), AppID: match[1]})
	}
	return res, nil
}

func unescapeHTML(raw string) string {
	d, err := html.Parse(bytes.NewReader([]byte(raw)))
	if err != nil {
		return raw
	}
	return scrape.Text(d)
}

func ensureHasLanguage(page string) string {
	parsed, err := url.Parse(page)
	if err != nil {
		return page
	}
	if parsed.Query().Get("hl") == "" {
		q := parsed.Query()
		q.Set("hl", "en")
		parsed.RawQuery = q.Encode()
		return parsed.String()
	}
	return page
}
