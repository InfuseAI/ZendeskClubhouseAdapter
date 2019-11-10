package main

import (
	".."
	"encoding/json"
	"log"
	"os"
)


func main() {
	token := os.Getenv("CH_TOKEN")
	c := cloudfunction.ClubHouse{token}

	iteration := cloudfunction.ClubHouseIteration{}
	err := c.CurrentIteration(&iteration)
	if err != nil {
		panic(err)
	}


	zendesk := cloudfunction.ZendeskTicket{
		Title:        "Test",
		Description:  "Test for zendesk",
		Organization: "InfuseAI",
		ID:           777,
		URL:          "https://infuseai.io",
	}
	clubhouse := cloudfunction.ClubHouseStory{}

	cloudfunction.ZendeskToClubHouse(&zendesk, &clubhouse)
	clubhouse.IterationID = iteration.ID
	jsonByte, err := json.Marshal(clubhouse)
	log.Printf("%v\n", string(jsonByte))
}