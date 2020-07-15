package models

type Player struct {
	Id  int    `json:"id" example:"15570"`
	Url string `json:"url" example:"https://transfermarkt.com/igor-akinfeev/profil/spieler/15570"`
}

type PlayerDetail struct {
	Name                  string   `json:"name" example:"Emanuel Mammana"`
	ImageUrl              string   `json:"image_url" example:"https://tmssl.akamaized.net/images/portrait/header/275512-1572532565.png?lm=1572532580"`
	BirthDate             *string  `json:"birth_date" example:"Feb 10, 1996"`         // nullable
	BirthCountry          *int     `json:"birth_country" example:"9"`                 // nullable
	Age                   *int     `json:"age, null" example:"24"`                    // nullable
	Height                *int     `json:"height, null" example:"183"`                // nullable
	Country               *int     `json:"country" example:"9"`                       // nullable
	Country2              *int     `json:"country_2" example:"75"`                    // nullable
	CurrentClub           *int     `json:"current_club" example:"964"`                // nullable
	CurrentRental         *int     `json:"current_rental" example:""`                 // nullable
	ContractExpires       *string  `json:"contract_expires"`                          // nullable
	ContractRentalExpires *string  `json:"contract_rental_expires"`                   // nullable
	Position              *string  `json:"position" example:"Defender - Centre-Back"` // nullable
	ShockFoot             *string  `json:"shock_foot" example:"right"`                // nullable
	Contacts              []string `json:"contracts" example:"http://www.instagram.com/emanuel_mammana24/"`
}
