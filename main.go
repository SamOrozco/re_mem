package main

import (
	"encoding/json"
	"re-mem/re_mem"
)

const jsonData = `{
  "data": [
    {
      "CountryName": "United States",
      "Description": "Texas",
      "StateCode": "TX"
    },
    {
      "CountryName": "United States",
      "Description": "United States Minor Outlying Islands (see also separate entry under UM)",
      "StateCode": "UM"
    },
    {
      "CountryName": "United States",
      "Description": "Utah",
      "StateCode": "UT"
    },
    {
      "CountryName": "United States",
      "Description": "Virginia",
      "StateCode": "VA"
    },
    {
      "CountryName": "United States",
      "Description": "Virgin Islands, U.S. (see also separate entry under VI)",
      "StateCode": "VI"
    },
    {
      "CountryName": "United States",
      "Description": "Vermont",
      "StateCode": "VT"
    },
    {
      "CountryName": "United States",
      "Description": "Washington",
      "StateCode": "WA"
    },
    {
      "CountryName": "United States",
      "Description": "Wisconsin",
      "StateCode": "WI"
    },
    {
      "CountryName": "United States",
      "Description": "West Virginia",
      "StateCode": "WV"
    },
    {
      "CountryName": "United States",
      "Description": "Wyoming",
      "StateCode": "WY"
    }
  ]
}`

func main() {
	store := re_mem.NewLocalStorage("/Users/samorozco/first_db")
	stateCollection, err := store.GetCollection("states")
	if err != nil {
		panic(err)
	}

	// parse json into map
	jsonMap := make(map[string]interface{}, 0)
	err = json.Unmarshal([]byte(jsonData), &jsonMap)
	if err != nil {
		panic(err)
	}

	// extact data field
	dataField := jsonMap["data"]

	records := dataField.([]interface{})
	for _, rec := range records {
		key, err := stateCollection.Create(rec)
		if err != nil {
			panic(err)
		}
		println(key)
	}
}
