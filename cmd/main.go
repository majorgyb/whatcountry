package main

import (
	"fmt"
	"log"

	"gitlab.com/gyb/whatcountry"
)

func main() {
	wc, err := whatcountry.LoadCountries("../data/ne_10m_admin_1_states_provinces.geojson")
	if err != nil {
		log.Fatal(err)
	}

	t := wc.FindPoint(12.553, 48.811)
	fmt.Println(t)
}
