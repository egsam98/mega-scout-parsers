package parsers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/utils"
	"github.com/egsam98/MegaScout/utils/slices"
	strings2 "github.com/egsam98/MegaScout/utils/strings"
	errors2 "github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func PlayerDetail(playerUrl string) (_ *models.PlayerDetail, err error) {
	doc, fetchHtmlErr := utils.FetchHtml(playerUrl)
	if fetchHtmlErr != nil {
		return nil, errors2.WithStack(fetchHtmlErr)
	}

	info := doc.Find("table.auflistung tr")

	citizenships := make([]*int, 2)

	var innerErr error
	findByTh(info, "Citizenship").Find("img").EachWithBreak(func(i int, img *goquery.Selection) bool {
		citizenships[i], err = country(playerUrl, img)
		if err != nil {
			innerErr = errors2.WithStack(err)
			return false
		}
		return true
	})
	if innerErr != nil {
		return nil, innerErr
	}

	birthPlace := findByTh(info, "Place of birth")
	imageUrl, _ := doc.Find(".dataBild > img").First().Attr("src")

	currentClub, err := team(findByTh(info, "Current club").Find("a").First())
	if err != nil {
		return nil, errors2.WithStack(err)
	}

	onLoanFrom, err := team(findByTh(info, "On loan from").Find("a").First())
	if err != nil {
		return nil, errors2.WithStack(err)
	}

	var contractExpires *string
	var currentRental *int
	var contractRentalExpires *string

	if result := strings2.Delete(findByTh(info, "Contract expires").Text(), "-"); result != "" {
		contractExpires = &result
		if onLoanFrom != nil {
			currentRental = onLoanFrom
			contractRentalExpires = new(string)
			*contractRentalExpires = *contractExpires
			*contractExpires = strings2.Delete(findByTh(info, "Contract there expires").Text(), "-")
		}
	}

	var birthDate *string
	if result := findByTh(info, "Date of birth").Text(); result != "" {
		birthDate = new(string)
		*birthDate = strings.Trim(result, "\n\t ")
	}

	var age *int
	if result, err := strconv.Atoi(findByTh(info, "Age").Text()); err == nil {
		age = &result
	}

	var height *int
	if result, err := strconv.Atoi(regexp.MustCompile(`[^0-9]`).
		ReplaceAllString(findByTh(info, "Height").Text(), "")); err == nil {
		height = &result
	}

	var position *string
	if result := strings.Trim(findByTh(info, "Position").Text(), "\n\t "); result != "" {
		position = &result
	}

	var shockFoot *string
	if result := findByTh(info, "Foot").Text(); result != "" {
		shockFoot = &result
	}

	contacts := findByTh(info, "Social media").Find("a").Map(func(_ int, a *goquery.Selection) string {
		url, exists := a.Attr("href")
		if !exists {
			innerErr = errors2.New("href attr's absent")
			return ""
		}
		return url
	})

	if innerErr != nil {
		return nil, innerErr
	}

	birthCountry, err := country(playerUrl, birthPlace.Find("img").First())
	if err != nil {
		return nil, errors2.WithStack(err)
	}

	transfers, err := transfers(playerUrl, doc)
	if err != nil {
		return nil, errors2.WithStack(err)
	}

	return &models.PlayerDetail{
		Name:                  doc.Find(".dataName > h1").Text(),
		ImageUrl:              imageUrl,
		BirthDate:             birthDate,
		BirthCountry:          birthCountry,
		Age:                   age,
		Height:                height,
		Country:               citizenships[0],
		Country2:              citizenships[1],
		CurrentClub:           currentClub,
		CurrentRental:         currentRental,
		ContractExpires:       contractExpires,
		ContractRentalExpires: contractRentalExpires,
		Position:              position,
		ShockFoot:             shockFoot,
		Contacts:              contacts,
		Transfers:             transfers,
	}, nil
}

func findByTh(info *goquery.Selection, header string) *goquery.Selection {
	tr := info.FilterFunction(func(i int, s *goquery.Selection) bool {
		thText := strings.Trim(s.Find("th").Text(), "\t\n ")
		return strings.HasPrefix(thText, header)
	}).First()
	return tr.Find("td").First()
}

func transfers(playerUrl string, doc *goquery.Document) (transfers []models.Transfer, _ error) {
	var innerErr error
	doc.Find(".box.transferhistorie tbody > tr.zeile-transfer").EachWithBreak(func(_ int, tr *goquery.Selection) bool {
		tds := tr.Find("td")

		transferType := 0
		var fee *string
		if result := tds.Eq(11).Text(); result != "" && result != "-" && result != "?" {
			fee = &result
			if *fee == "Loan" {
				transferType = 1
			}
		}

		var date *string
		dateFormatted := tds.Eq(1).Text()
		if dateFormatted != "" {
			result, err := time.Parse("Jan 2, 2006", dateFormatted)
			if err != nil {
				innerErr = errors2.WithStack(err)
				return false
			}
			formatted := result.Format("02-01-2006")
			date = &formatted
		}

		var cost *string
		if result := tds.Eq(10).Text(); result != "" && result != "-" && result != "?" {
			cost = &result
		}

		var fromTeam *string
		if result, exists := tds.Eq(2).Find("a").First().Attr("id"); exists {
			fromTeam = &result
		}

		toTeam, exists := tds.Eq(6).Find("a").First().Attr("id")
		if !exists {
			innerErr = errors2.Errorf("%s: transfer to team is absent", playerUrl)
			return false
		}

		transfers = append(transfers, models.Transfer{
			TransferType: transferType,
			Date:         date,
			Season:       tds.First().Text(),
			FromTeam:     fromTeam,
			ToTeam:       toTeam,
			Cost:         cost,
			Fee:          fee,
		})
		return true
	})

	if innerErr != nil {
		return nil, innerErr
	}
	return transfers, nil
}

func country(playerUrl string, img *goquery.Selection) (*int, error) {
	if img.Size() == 0 {
		return nil, nil
	}
	src, exists := img.Attr("src")
	if !exists {
		return nil, errors2.Errorf("%s: not an image element", playerUrl)
	}
	sep := slices.String_Last(strings.Split(src, "/"))
	id, err := strconv.Atoi(strings.Split(sep, ".")[0])
	if err != nil {
		return nil, errors2.Errorf("%s: %v", playerUrl, err)
	}
	return &id, nil
}

func team(a *goquery.Selection) (*int, error) {
	idStr, exists := a.Attr("id")
	if !exists {
		return nil, nil
	}
	id, err := strconv.Atoi(idStr)
	return &id, errors2.WithStack(err)
}
