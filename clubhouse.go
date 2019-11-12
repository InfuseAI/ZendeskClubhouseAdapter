package cloudfunction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ClubHouseIteration struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Name   string `json:"name"`
}

type ClubHouseExternalTicket struct {
	ID  string `json:"external_id"`
	URL string `json:"external_url"`
}

type ClubHouseStory struct {
	ProjectID       int                       `json:"project_id"`
	StoryType       string                    `json:"story_type"`
	Name            string                    `json:"name"`
	Description     string                    `json:"description"`
	ExternalTickets []ClubHouseExternalTicket `json:"external_tickets"`
	ExternalID      string                    `json:"external_id"`
	IterationID     int                       `json:"iteration_id"`
}

type AbstractClubHouse interface {
	CurrentIteration(*ClubHouseIteration) error
	CreateStory(*ClubHouseStory) error
}

type ClubHouse struct {
	Token string
}

type MockClubHouse struct {
	Token string
}

func ClubHouseBuilder(token string) AbstractClubHouse {
	if token == "MOCK_CLUBHOUSE" {
		return &MockClubHouse{token}
	}
	return &ClubHouse{token}
}

func (c *ClubHouse) CurrentIteration(currentIteration *ClubHouseIteration) error {
	var iterations []ClubHouseIteration
	URL := "https://api.clubhouse.io/api/v3/iterations?token=" + c.Token
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&iterations)
	if err != nil {
		log.Fatal(err)
	}

	latestIterationID := 0
	for _, v := range iterations {
		if v.Status == "started" {
			if v.ID > latestIterationID {
				*currentIteration = v
				latestIterationID = v.ID
			}
		}
	}

	if latestIterationID == 0 {
		return fmt.Errorf("No started iterations")
	}
	return nil
}

func (c *MockClubHouse) CurrentIteration(currentIteration *ClubHouseIteration) error {
	return nil
}

func (c *ClubHouse) CreateStory(story *ClubHouseStory) error {
	if story == nil {
		return fmt.Errorf("no story provided")
	}
	URL := fmt.Sprintf("https://api.clubhouse.io/api/v3/stories?token=%s", c.Token)
	requestBytes, err := json.Marshal(*story)
	if err != nil {
		return err
	}
	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf(resp.Status)
	}
	return nil
}

func (c *MockClubHouse) CreateStory(story *ClubHouseStory) error {
	return nil
}

func (c *ClubHouse) UpdateStory() error {
	return nil
}

func ZendeskToClubHouse(zendeskTicket *ZendeskTicket, clubhouseTicket *ClubHouseStory) {
	if zendeskTicket == nil || clubhouseTicket == nil {
		return
	}

	clubhouseTicket.Name = fmt.Sprintf("[%s] %s", zendeskTicket.Organization, zendeskTicket.Title)
	clubhouseTicket.Description = zendeskTicket.Description
	clubhouseTicket.ProjectID = 55
	clubhouseTicket.StoryType = "bug"
	clubhouseTicket.ExternalTickets = append(clubhouseTicket.ExternalTickets, ClubHouseExternalTicket{
		ID:  zendeskTicket.ID,
		URL: zendeskTicket.URL,
	})
	clubhouseTicket.ExternalID = fmt.Sprintf("zendesk-%s", zendeskTicket.ID)
}
