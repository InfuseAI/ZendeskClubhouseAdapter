package cloudfunction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ClubHouseProject struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ClubHoseWorkflow struct {
	EntityType string                   `json:"entity_type"`
	States     []ClubHouseWorkflowState `json:"states"`
	Name       string                   `json:"name"`
	ID         int                      `json:"id"`
}
type ClubHouseWorkflowState struct {
	EntityType string `json:"entity_type"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	ID         int    `json:"id"`
}

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
	ID              int                       `json:"id,omitempty"`
	ProjectID       int                       `json:"project_id"`
	StoryType       string                    `json:"story_type"`
	Name            string                    `json:"name"`
	Description     string                    `json:"description"`
	ExternalTickets []ClubHouseExternalTicket `json:"external_tickets"`
	ExternalID      string                    `json:"external_id"`
	IterationID     int                       `json:"iteration_id"`
	WorkflowStateID int                       `json:"workflow_state_id,omitempty"`
}

type AbstractClubHouse interface {
	CurrentIteration(*ClubHouseIteration) error
	GetStoryByExternalID(string, *ClubHouseStory) error
	GetWorkflowStateByName(string, string) (int, error)
	CreateStory(*ClubHouseStory) error
	AddCommentOnStory(int, string) error
	CloseStory(int, int) error
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

	if currentIteration == nil {
		return fmt.Errorf("no iteration provided")
	}

	URL := "https://api.clubhouse.io/api/v3/iterations?token=" + c.Token
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&iterations)
	if err != nil {
		return err
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
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
		return fmt.Errorf(resp.Status)
	}
	return nil
}

func (c *MockClubHouse) CreateStory(story *ClubHouseStory) error {
	return nil
}

func (c *ClubHouse) AddCommentOnStory(storyID int, text string) error {
	URL := fmt.Sprintf("https://api.clubhouse.io/api/v3/stories/%d/comments?token=%s", storyID, c.Token)
	payload := map[string]interface{}{"text": text}
	requestBytes, err := json.Marshal(payload)
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

func (c *MockClubHouse) AddCommentOnStory(storyID int, text string) error {
	return nil
}

func (c *ClubHouse) CloseStory(storyID int, workflowID int) error {
	URL := fmt.Sprintf("https://api.clubhouse.io/api/v3/stories/%d?token=%s", storyID, c.Token)
	payload := map[string]interface{}{"workflow_state_id": workflowID}
	requestBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, URL, bytes.NewBuffer(requestBytes))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(resp.Status)
	}
	return nil
}

func (c *MockClubHouse) CloseStory(storyID int, workflowID int) error {
	return nil
}

func (c *ClubHouse) GetStoryByExternalID(externalID string, story *ClubHouseStory) error {
	if story == nil {
		return fmt.Errorf("no story provided")
	}

	URL := fmt.Sprintf("https://api.clubhouse.io/api/v3/stories/search?token=%s", c.Token)
	payload := map[string]interface{}{"external_id": externalID}
	requestBytes, err := json.Marshal(payload)
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

	stories := []ClubHouseStory{}
	err = json.NewDecoder(resp.Body).Decode(&stories)
	if err != nil {
		return err
	}
	if len(stories) == 0 {
		return os.ErrNotExist
	}

	*story = stories[0]
	return nil
}

func (c *MockClubHouse) GetStoryByExternalID(externalID string, story *ClubHouseStory) error {
	if externalID == "zendesk-NON_EXIST_ID" {
		return os.ErrNotExist
	}
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

func (c *ClubHouse) GetWorkflowStateByName(workflowName string, stateName string) (int, error) {
	workflows := new([]ClubHoseWorkflow)
	URL := fmt.Sprintf("https://api.clubhouse.io/api/v3/workflows?token=%s", c.Token)

	resp, err := http.Get(URL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&workflows)
	if err != nil {
		return 0, fmt.Errorf(resp.Status)
	}

	for _, workflow := range *workflows {
		if workflow.Name == workflowName {
			for _, state := range workflow.States {
				if state.Name == stateName {
					return state.ID, nil
				}
			}
		}
	}

	return 0, os.ErrNotExist
}

func (c *MockClubHouse) GetWorkflowStateByName(workflowName string, stateName string) (int, error) {
	return 500000011, nil
}

func (c *ClubHouse) GetProjectByName(name string) (int, error) {
	projects := new([]ClubHouseProject)
	URL := fmt.Sprintf("https://api.clubhouse.io/api/v3/projects?token=%s", c.Token)

	resp, err := http.Get(URL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&projects)
	if err != nil {
		return 0, fmt.Errorf(resp.Status)
	}

	for  _, project := range *projects {
		if project.Name == name {
			return project.ID, nil
		}
	}

	return 0, os.ErrNotExist
}

func (c *MockClubHouse) GetProjectByName(name string) (int, error)  {
	return 55, nil
}
