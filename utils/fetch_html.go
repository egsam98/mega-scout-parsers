package utils

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/egsam98/MegaScout/utils/errors"
	"net/http"
)

func FetchHtml(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.NewFetchHtmlError(err)
	}

	if res.StatusCode == 500 {
		return nil, errors.NewTransfermarktError(url)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.NewFetchHtmlError(err)
	}

	return doc, nil
}
