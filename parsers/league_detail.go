package parsers

import (
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/utils"
	"github.com/pkg/errors"
)

func LeagueDetail(url string) (*models.LeagueDetail, error) {
	doc, err := utils.RetryFetchHtml(url, 5)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var logo *string
	if result, exists := doc.Find(".headerfoto > img").Attr("src"); exists {
		logo = &result
	} else {
		if result, exists := doc.Find(".dataBild > img").Attr("src"); exists {
			logo = &result
		}
	}

	name := doc.Find(".spielername-profil").Text()
	if name == "" {
		name = doc.Find(".dataName > h1").Text()
	}

	return &models.LeagueDetail{
		Name: name,
		Logo: logo,
	}, nil
}
