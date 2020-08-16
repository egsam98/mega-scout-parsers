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

	defer res.Body.Close()

	if res.StatusCode == 500 {
		return nil, errors.NewTransfermarktError(url)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.NewFetchHtmlError(err)
	}

	return doc, nil
}

func RetryFetchHtml(url string, times int) (doc *goquery.Document, err error) {
	for i := 0; i < times; i++ {
		if doc, err = FetchHtml(url); err == nil {
			return doc, nil
		}
		if _, ok := err.(*errors.TransfermarktError); ok {
			return nil, err
		}
	}
	return nil, err
}
