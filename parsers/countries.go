package parsers

import (
	"fmt"
	"github.com/egsam98/MegaScout/models"
	"github.com/egsam98/MegaScout/utils"
	"github.com/pariz/gountries"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/firefox"
	"os"
	"strconv"
	"time"
)

var gountry = gountries.New()
var exceptionalCountries = map[string][2]string{
	"American Virgin Islands":        {"Virgin Islands (US)", "VI"},
	"Bonaire":                        {"Bonaire", "BQ"},
	"Bosnia-Herzegovina":             {"Bosnia and Herzegovina", "BA"},
	"Botsuana":                       {"Botswana", "BW"},
	"British India":                  {"British India", "IO"},
	"Brunei Darussalam":              {"Brunei", "BN"},
	"Chinese Taipei (Taiwan)":        {"Taiwan", "TW"},
	"Congo":                          {"DR Congo", "CD"},
	"Cookinseln":                     {"Cook Islands", "CK"},
	"Cote d'Ivoire":                  {"Ivory Coast", "CI"},
	"CSSR":                           {"CSSR", ""},
	"Curacao":                        {"Curaçao", "CW"},
	"DDR":                            {"GDR", "DE"},
	"DR Congo":                       {"DR Congo", "CD"},
	"Eswatini":                       {"Eswatini", "SZ"},
	"Hongkong":                       {"Hong Kong", "HK"},
	"Zaire":                          {"DR Congo", "CD"},
	"England":                        {"England", "GB-ENG"},
	"Federated States of Micronesia": {"Federated States of Micronesia", "FM"},
	"Jugoslawien (SFR)":              {"Yugoslavia", ""},
	"Yugoslavia (Republic)":          {"Yugoslavia", ""},
	"Korea, North":                   {"Korea, North", "KP"},
	"Korea, South":                   {"Korea, South", "KR"},
	"Kosovo":                         {"Kosovo", "XK"},
	"Macao":                          {"Macau", "MO"},
	"Mariana Islands":                {"Mariana Islands", "MP"},
	"Netherlands Antilles":           {"Antilles", "NL"},
	"Netherlands East India":         {"Netherlands East India", "NL"},
	"Neukaledonien":                  {"New Caledonia", "NC"},
	"North Macedonia":                {"North Macedonia", "MK"},
	"Northern Ireland":               {"Northern Ireland", "GB-NIR"},
	"Osttimor":                       {"East Timor", "TL"},
	"Palästina":                      {"Palestine", "PS"},
	"People's republic of the Congo": {"Congo", "CG"},
	"Saarland":                       {"Germany", "DE"},
	"Saint-Martin":                   {"St. Martin", "FR"},
	"Sao Tome and Principe":          {"São Tomé and Príncipe", "ST"},
	"Scotland":                       {"Scotland", "GB-SCT"},
	"Serbia and Montenegro":          {"Serbia and Montenegro", ""},
	"Southern Sudan":                 {"Southern Sudan", "SS"},
	"St. Kitts & Nevis":              {"St. Kitts & Nevis", "KN"},
	"St. Lucia":                      {"St. Lucia", "LC"},
	"St. Vincent & Grenadinen":       {"St. Vincent & Grenadinen", "VC"},
	"Tahiti":                         {"French Polynesia", "PF"},
	"The Gambia":                     {"Gambia", "GM"},
	"Tibet":                          {"China", "CN"},
	"Turks- and Caicosinseln":        {"Turks and Caicos", "TC"},
	"UdSSR":                          {"USSR", ""},
	"Vatican":                        {"Vatican", "VA"},
	"Wales":                          {"Wales", "GB-WLS"},
	"Zanzibar":                       {"Tanzania", "TZ"},
}

func Countries() (utils.Set, error) {
	service, driver, err := initSelenium()
	if err != nil {
		return nil, err
	}
	defer service.Stop()
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
	time.Sleep(1 * time.Second)
	countries := utils.NewSet()
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
		res, _ := driver.ExecuteScript("return arguments[0].textContent", []interface{}{elem})
		name, code, err := _ISOCode(res.(string))
		if err != nil {
			return nil, err
		}
		countries.Add(models.Country{
			Id:   id,
			Name: name,
			Code: code,
		}, "Name")
	}
	return countries, nil
}

func initSelenium() (*selenium.Service, selenium.WebDriver, error) {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	service, err := selenium.NewGeckoDriverService(os.Getenv("GECKODRIVER_PATH"), port+1)
	if err != nil {
		return nil, nil, err
	}

	caps := selenium.Capabilities{"browserName": "firefox"}
	firefoxCaps := firefox.Capabilities{
		Binary: os.Getenv("FIREFOX_BIN"),
		Args:   []string{"--headless"},
	}
	caps.AddFirefox(firefoxCaps)
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", port+1))
	return service, driver, err
}

func _ISOCode(countryName string) (string, string, error) {
	country, err := gountry.FindCountryByName(countryName)
	if err != nil {
		data, exists := exceptionalCountries[countryName]
		if !exists {
			return "", "", fmt.Errorf("Country %s's not found", countryName)
		}
		return data[0], data[1], nil
	}
	return countryName, country.Alpha2, nil
}
