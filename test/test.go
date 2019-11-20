package main

import (
	".."
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

	//zendesk := cloudfunction.ZendeskTicket{
	//	Title:        "Test",
	//	Description:  "Test for zendesk",
	//	Organization: "InfuseAI",
	//	ID:           "777",
	//	URL:          "https://infuseai.io",
	//}
	//clubhouse := cloudfunction.ClubHouseStory{}

	//cloudfunction.ZendeskToClubHouse(&zendesk, &clubhouse)
	//clubhouse.IterationID = iteration.ID

	//err = c.CreateStory(&clubhouse)
	//if err != nil {
	//	panic(err)
	//}

	//jsonByte, err := json.Marshal(clubhouse)
	//log.Printf("%v\n", string(jsonByte))

	story := cloudfunction.ClubHouseStory{}
	err = c.GetStoryByExternalID("zendesk-539", &story)
	if err != nil {
		log.Printf("%v\n", err)
	}

	log.Printf("%v\n", story)
	//err = c.AddCommentOnStory(story.ID, "Test message from ZendeskClubhouseAdapter")
}
