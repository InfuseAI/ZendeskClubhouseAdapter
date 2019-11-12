package cloudfunction

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"
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

	if token == "" ||
		zendeskTicket.Title == "" ||
		zendeskTicket.ID == 0 ||
		zendeskTicket.URL == "" {
		return os.ErrInvalid
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

func verifyBasicAuth(w http.ResponseWriter, r *http.Request, user string, password string) bool {
	basicAuthPrefix := "Basic "
	auth := r.Header.Get("Authorization")

	if user == "" && password == "" {
		return true
	}

	// Decode auth payload by base64
	if strings.HasPrefix(auth, basicAuthPrefix) {
		payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])
		if err == nil {
			pair := bytes.SplitN(payload, []byte(":"), 2)
			if len(pair) == 2 &&
				bytes.Equal(pair[0], []byte(user)) &&
				bytes.Equal(pair[1], []byte(password)) {
				return true
			}
		}
	}

	return false
}

func ZendeskClubhouseAdapter(w http.ResponseWriter, r *http.Request) {
	var user = os.Getenv("AUTH_USER")
	var password = os.Getenv("AUTH_PASSWORD")
	var method = r.Method
	var err error

	// Check http authorization
	if verifyBasicAuth(w, r, user, password) == false {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if method == http.MethodPost {
		err = createTicket(r)
	}

	if err != nil {
		if err == os.ErrInvalid {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}
