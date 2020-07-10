package parsers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tebeka/selenium"
	"os"
	"strconv"
)

func Countries() ([]gin.H, error) {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	service, err := selenium.NewGeckoDriverService("/home/egor/geckodriver", port+1)
	if err != nil {
		return nil, err
	}
	defer service.Stop()
	driver, err := selenium.NewRemote(nil, fmt.Sprintf("http://localhost:%d", port+1))
	if err != nil {
		return nil, err
	}
	defer driver.Quit()
	if err := driver.Get("https://www.transfermarkt.com"); err != nil {
		return nil, err
	}
	elem, err := driver.FindElement(selenium.ByID, "land_select_breadcrumb_chzn")
	if err != nil {
		return nil, err
	}
	if err := elem.Click(); err != nil {
		return nil, err
	}
	countries := make([]gin.H, 0)
	elems, err := driver.FindElements(selenium.ByCSSSelector, "#land_select_breadcrumb > option")
	if err != nil {
		return nil, err
	}
	for _, elem := range elems {
		idStr, err := elem.GetAttribute("value")
		if err != nil {
			continue
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		name, _ := driver.ExecuteScript("return arguments[0].textContent", []interface{}{elem})
		countries = append(countries, map[string]interface{}{
			"country_id": id,
			"name":       name,
			"code":       name, //TODO: доделать
		})
	}
	return countries, nil
}
