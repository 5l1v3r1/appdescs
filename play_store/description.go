package main

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func AppDescription(id string) (string, error) {
	url := "https://play.google.com/store/apps/details?id=" + id + "&hl=en"
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}
	element, ok := scrape.Find(doc, scrape.ByClass("show-more-content"))
	if !ok {
		return "", errors.New("missing expected DOM element")
	}
	sub, ok := scrape.Find(element, func(n *html.Node) bool {
		return n.DataAtom == atom.Div && !scrape.ByClass("show-more-content")(n)
	})
	if ok {
		element = sub
	}
	var res bytes.Buffer
	child := element.FirstChild
	for child != nil {
		html.Render(&res, child)
		child = child.NextSibling
	}
	return res.String(), nil
}
