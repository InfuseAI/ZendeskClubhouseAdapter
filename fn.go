package cloudfunction

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type ZendeskTicket struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Organization string `json:"organization"`
	ID           int    `json:"id"`
	URL          string `json:"url"`
}

func createTicket(r *http.Request) error {
	var token = os.Getenv("CH_TOKEN")
	var zendeskTicket = ZendeskTicket{}
	var clubhouseStory = ClubHouseStory{}
	var currentIteration = ClubHouseIteration{}
	var clubhouse = ClubHouseBuilder(token)

	// Parse request body
	var decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&zendeskTicket)
	if err != nil {
		return err
	}

	// Prepare Clubhouse Story
	ZendeskToClubHouse(&zendeskTicket, &clubhouseStory)

	// Get current Clubhouse iteration
	err = clubhouse.CurrentIteration(&currentIteration)
	if err != nil {
		return err
	}
	clubhouseStory.IterationID = currentIteration.ID

	// Create Clubhouse Story
	err = clubhouse.CreateStory(&clubhouseStory)
	if err != nil {
		return err
	}

	return nil
}

func ZendeskClubhouseAdapter(w http.ResponseWriter, r *http.Request) {
	var method = r.Method
	var err error
	r.Header.Get("Authorization")

	if method == http.MethodPost {
		err = createTicket(r)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusCreated)
}
