package parser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"strings"
)

// GetDoc return Document object of the HTML string
func GetDoc(html string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return doc, nil
}

// Title get title
func Title(doc *goquery.Document) string {
	var title string
	h1Elem := doc.Find("h1").First()
	h1Title, found := h1Elem.Attr("title")
	if !found {
		h1Title = h1Elem.Text()
	}
	title = strings.Replace(strings.TrimSpace(h1Title), "\n", "", -1)
	if title == "" {
		// Bilibili: Some movie page got no h1 tag
		title, _ = doc.Find("meta[property=\"og:title\"]").Attr("content")
	}
	if title == "" {
		title = doc.Find("title").Text()
	}
	return title
}
