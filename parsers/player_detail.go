package parsers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/utils"
	"github.com/egsam98/MegaScout/utils/slices"
	strings2 "github.com/egsam98/MegaScout/utils/strings"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func PlayerDetail(playerUrl string) (*models.PlayerDetail, error) {
	doc, err := utils.FetchHtml(playerUrl)
	if err != nil {
		return nil, err
	}

	info := doc.Find("table.auflistung tr")

	citizenships := make([]*int, 2)
	findByTh(info, "Citizenship").Find("img").Each(func(i int, img *goquery.Selection) {
		citizenships[i] = country(playerUrl, img)
	})

	birthPlace := findByTh(info, "Place of birth")
	imageUrl, _ := doc.Find(".dataBild > img").First().Attr("src")
	currentClub := team(findByTh(info, "Current club").Find("a").First())
	onLoanFrom := team(findByTh(info, "On loan from").Find("a").First())

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
			panic(err)
		}
		return url
	})

	return &models.PlayerDetail{
		Name:                  doc.Find(".dataName > h1").Text(),
		ImageUrl:              imageUrl,
		BirthDate:             birthDate,
		BirthCountry:          country(playerUrl, birthPlace.Find("img").First()),
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
		Transfers:             transfers(playerUrl, doc),
	}, nil
}

func findByTh(info *goquery.Selection, header string) *goquery.Selection {
	tr := info.FilterFunction(func(i int, s *goquery.Selection) bool {
		thText := strings.Trim(s.Find("th").Text(), "\t\n ")
		return strings.HasPrefix(thText, header)
	}).First()
	return tr.Find("td").First()
}

func transfers(playerUrl string, doc *goquery.Document) (transfers []models.Transfer) {
	doc.Find(".box.transferhistorie tbody > tr.zeile-transfer").Each(func(_ int, tr *goquery.Selection) {
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
				panic(err)
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
			panic(fmt.Errorf("%s: transfer to team is absent", playerUrl))
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
	})
	return transfers
}

func country(playerUrl string, img *goquery.Selection) *int {
	if img.Size() == 0 {
		return nil
	}
	src, exists := img.Attr("src")
	if !exists {
		panic(fmt.Errorf("%s: not an image element", playerUrl))
	}
	sep := slices.String_Last(strings.Split(src, "/"))
	id, err := strconv.Atoi(strings.Split(sep, ".")[0])
	if err != nil {
		panic(fmt.Errorf("%s: %v", playerUrl, err))
	}
	return &id
}

func team(a *goquery.Selection) *int {
	idStr, exists := a.Attr("id")
	if !exists {
		return nil
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(err)
	}
	return &id
}
