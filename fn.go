package cloudfunction

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type ZendeskTicket struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Organization string `json:"organization"`
	ID           string `json:"id"`
	URL          string `json:"url"`
	Status 	     string `json:"status"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
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
		log.Fatalln("Zendesk ticket decode error")
		return err
	}

	if token == "" ||
		zendeskTicket.Title == "" ||
		zendeskTicket.ID == "" ||
		zendeskTicket.URL == "" {
		return os.ErrInvalid
	}

	// Prepare Clubhouse Story
	clubhouseStoryType := getEnv("CLUBHOUSE_STORY_TYPE", "chore")
	clubhouseProjectID, err := clubhouse.GetProjectByName(getEnv("CLUBHOUSE_PROJECT", "Support"))
	clubhouseTeamID, err := clubhouse.GetTeamByName(getEnv("CLUBHOUSE_TEAM", "Support"))

	clubhouseWorkflow := getEnv("CLUBHOUSE_WORKFLOW", "Support")
	clubhouseCreatedState := getEnv("CLUBHOUSE_CREATED_STATE", "Created" )
	clubhouseCreatedStateID, err := clubhouse.GetWorkflowStateByName(clubhouseWorkflow, clubhouseCreatedState)

	if err != nil {
		return err
	}
	ZendeskToClubHouse(&zendeskTicket, &clubhouseStory, clubhouseProjectID, clubhouseTeamID, clubhouseStoryType, clubhouseCreatedStateID)

	// Get current Clubhouse iteration
	err = clubhouse.CurrentIteration(&currentIteration)
	if err != nil {
		log.Fatalln("Fail to get current iteration")
		return err
	}
	clubhouseStory.IterationID = currentIteration.ID

	// Create Clubhouse Story
	err = clubhouse.CreateStory(&clubhouseStory)
	if err != nil {
		log.Fatalln("Fail to create story")
		return err
	}

	return nil
}

func updateTicket(r *http.Request) error {
	var token = os.Getenv("CH_TOKEN")
	var zendeskTicket = ZendeskTicket{}
	var story = ClubHouseStory{}
	var clubhouse = ClubHouseBuilder(token)

	// Parse request body
	var decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&zendeskTicket)
	if err != nil {
		log.Fatalln("Zendesk ticket decode error")
		return err
	}

	if token == "" ||
		zendeskTicket.ID == "" {
		return os.ErrInvalid
	}

	externalID := fmt.Sprintf("zendesk-%s", zendeskTicket.ID)
	err = clubhouse.GetStoryByExternalID(externalID, &story)
	if err != nil {
		return err
	}

	err = clubhouse.AddCommentOnStory(story.ID, zendeskTicket.Description)
	if err != nil {
		return err
	}

	if zendeskTicket.Status == "Pending" {
		workflow := getEnv("CLUBHOUSE_WORKFLOW", "Dev")
		pendingState := getEnv("CLUBHOUSE_PENDING_STATE", "Blocks" )
		pendingStateID, err := clubhouse.GetWorkflowStateByName(workflow, pendingState)
		if err != nil {
			return err
		}

		if pendingStateID != story.WorkflowStateID {
			err = clubhouse.UpdateStoryState(story.ID, pendingStateID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func closeTicket(r *http.Request) error {
	var token = os.Getenv("CH_TOKEN")
	var zendeskTicket = ZendeskTicket{}
	var story = ClubHouseStory{}
	var clubhouse = ClubHouseBuilder(token)

	// Parse request body
	var decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&zendeskTicket)
	if err != nil {
		log.Fatalln("Zendesk ticket decode error")
		return err
	}

	if token == "" ||
		zendeskTicket.ID == "" {
		return os.ErrInvalid
	}

	externalID := fmt.Sprintf("zendesk-%s", zendeskTicket.ID)
	err = clubhouse.GetStoryByExternalID(externalID, &story)
	if err != nil {
		return err
	}

	workflow := getEnv("CLUBHOUSE_WORKFLOW", "Dev")
	completedState := getEnv("CLUBHOUSE_COMPLETED_STATE", "Completed")
	completedStateID, err := clubhouse.GetWorkflowStateByName(workflow, completedState)
	if err != nil {
		return err
	}

	if completedStateID == story.WorkflowStateID {
		return nil
	}

	return clubhouse.UpdateStoryState(story.ID, completedStateID)
}

func verifyBasicAuth(w http.ResponseWriter, r *http.Request, user string, password string) bool {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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
	} else if method == http.MethodPut {
		err = updateTicket(r)
	} else if method == http.MethodDelete {
		err = closeTicket(r)
	} else {
		// Unsupported method
		w.WriteHeader(http.StatusTeapot)
		return
	}

	if err != nil {
		if err == os.ErrInvalid {
			w.WriteHeader(http.StatusBadRequest)
		} else if err == os.ErrNotExist {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Printf("[Error] %s", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
