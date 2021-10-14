package cloudfunction

import (
	"github.com/jarcoal/httpmock"
	"testing"
)

func TestClubHouse_CurrentIteration(t *testing.T) {
	type fields struct {
		Token string
	}
	type args struct {
		currentIteration *ClubHouseIteration
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		responseBody string
		expectID     int
	}{
		{
			name:         "nil iteration obj",
			fields:       fields{Token: "test"},
			args:         args{currentIteration: nil},
			wantErr:      true,
			responseBody: `{}`,
			expectID:     0,
		},
		{
			name:         "Cet current iteration",
			fields:       fields{"test"},
			args:         args{new(ClubHouseIteration)},
			wantErr:      false,
			responseBody: `[{"id": 123, "status": "started", "name": "Fake iteration"}]`,
			expectID:     123,
		},
		{
			name:         "No any started iteration",
			fields:       fields{Token: "test"},
			args:         args{currentIteration: new(ClubHouseIteration)},
			wantErr:      true,
			responseBody: `[{"id": 123, "status": "end", "name": "Fake iteration 1"}, {"id": 234, "status": "end", "name": "Fake iteration 2"}]`,
			expectID:     0,
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClubHouse{
				Token: tt.fields.Token,
			}
			httpmock.RegisterResponder("GET", ClubHouseAPIURL + "/api/v3/iterations",
				httpmock.NewStringResponder(200, tt.responseBody))
			if err := c.CurrentIteration(tt.args.currentIteration); (err != nil) != tt.wantErr {
				t.Errorf("CurrentIteration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.currentIteration != nil && tt.expectID != tt.args.currentIteration.ID {
				t.Errorf("CurrentIteration ID should be %d not %d", tt.expectID, tt.args.currentIteration.ID)
			}
		})
	}
}

func TestClubHouse_CreateStory(t *testing.T) {
	type fields struct {
		Token string
	}
	type args struct {
		story *ClubHouseStory
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Create clubhouse story",
			fields:  fields{"test"},
			args:    args{new(ClubHouseStory)},
			wantErr: false,
		},
		{
			name:    "nil story obj ",
			fields:  fields{"test"},
			args:    args{nil},
			wantErr: true,
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClubHouse{
				Token: tt.fields.Token,
			}
			httpmock.RegisterResponder("POST", ClubHouseAPIURL + "/api/v3/stories",
				httpmock.NewStringResponder(201, "{}"))
			if err := c.CreateStory(tt.args.story); (err != nil) != tt.wantErr {
				t.Errorf("CreateStory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClubHouse_GetStoryByExternalID(t *testing.T) {
	type fields struct {
		Token string
	}
	type args struct {
		externalID string
		story      *ClubHouseStory
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		responseBody string
		expectID     int
		wantErr      bool
	}{
		{
			name:         "no story obj",
			fields:       fields{"test"},
			args:         args{"zendesk-777", nil},
			responseBody: `{}`,
			expectID:     0,
			wantErr:      true,
		},
		{
			name:         "get story by external ID",
			fields:       fields{"test"},
			args:         args{"zendesk-777", new(ClubHouseStory)},
			responseBody: `[{"id": 777, "iteration_id": 123, "workflow_state_id": 123}]`,
			expectID:     777,
			wantErr:      false,
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClubHouse{
				Token: tt.fields.Token,
			}
			httpmock.RegisterResponder("POST", ClubHouseAPIURL + "/api/v3/stories/search",
				httpmock.NewStringResponder(201, tt.responseBody))
			if err := c.GetStoryByExternalID(tt.args.externalID, tt.args.story); (err != nil) != tt.wantErr {
				t.Errorf("GetStoryByExternalID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.story != nil && tt.expectID != tt.args.story.ID {
				t.Errorf("Story ID should be %d not %d", tt.expectID, tt.args.story.ID)
			}
		})
	}
}

func TestClubHouse_AddCommentOnStory(t *testing.T) {
	type fields struct {
		Token string
	}
	type args struct {
		storyID int
		text    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Add comment",
			fields:  fields{"test"},
			args:    args{777, "Unit test"},
			wantErr: false,
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClubHouse{
				Token: tt.fields.Token,
			}
			httpmock.RegisterResponder("POST", `=~^https://api\.clubhouse\.io/api/v3/stories/.*/comments`,
				httpmock.NewStringResponder(201, `{}`))
			if err := c.AddCommentOnStory(tt.args.storyID, tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("AddCommentOnStory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClubHouse_UpdateStoryState(t *testing.T) {
	type fields struct {
		Token string
	}
	type args struct {
		storyID    int
		workflowID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "close story",
			fields:  fields{"test"},
			args:    args{777, 123},
			wantErr: false,
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClubHouse{
				Token: tt.fields.Token,
			}
			httpmock.RegisterResponder("PUT", `=~^https://api\.clubhouse\.io/api/v3/stories/.*`,
				httpmock.NewStringResponder(200, `{}`))
			if err := c.UpdateStoryState(tt.args.storyID, tt.args.workflowID); (err != nil) != tt.wantErr {
				t.Errorf("UpdateStoryState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

var workflowResponse = `[
  {
    "auto_assign_owner": true,
    "created_at": "2016-12-31T12:30:00Z",
    "default_state_id": 123,
    "description": "foo",
    "entity_type": "foo",
    "id": 123,
    "name": "Dev",
    "project_ids": [123],
    "states": [{
      "color": "foo",
      "created_at": "2016-12-31T12:30:00Z",
      "description": "foo",
      "entity_type": "foo",
      "id": 123,
      "name": "completed",
      "num_stories": 123,
      "num_story_templates": 123,
      "position": 123,
      "type": "foo",
      "updated_at": "2016-12-31T12:30:00Z",
      "verb": "foo"
    }],
    "team_id": 123,
    "updated_at": "2016-12-31T12:30:00Z"
  }
]`
func TestClubHouse_GetWorkflowStateByName(t *testing.T) {
	type fields struct {
		Token string
	}
	type args struct {
		workflowName string
		stateName    string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		responseBody string
		want         int
		wantErr      bool
	}{
		{
			name:         "Get workflow state",
			fields:       fields{"test"},
			args:         args{"Dev", "completed"},
			responseBody: workflowResponse,
			want:         123,
			wantErr:      false,
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClubHouse{
				Token: tt.fields.Token,
			}
			httpmock.RegisterResponder("GET", ClubHouseAPIURL + "/api/v3/workflows",
				httpmock.NewStringResponder(200, tt.responseBody))
			got, err := c.GetWorkflowStateByName(tt.args.workflowName, tt.args.stateName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWorkflowStateByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetWorkflowStateByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

var projectsResponse = `[
  {
    "abbreviation": "foo",
    "archived": true,
    "color": "foo",
    "created_at": "2016-12-31T12:30:00Z",
    "days_to_thermometer": 123,
    "description": "foo",
    "entity_type": "foo",
    "external_id": "foo",
    "follower_ids": ["12345678-9012-3456-7890-123456789012"],
    "id": 123,
    "iteration_length": 123,
    "name": "foo",
    "show_thermometer": true,
    "start_time": "2016-12-31T12:30:00Z",
    "stats": {
      "num_points": 123,
      "num_stories": 123
    },
    "team_id": 123,
    "updated_at": "2016-12-31T12:30:00Z"
  },
  {
    "id": 55,
    "name": "Support"
  }
]`
func TestClubHouse_GetProjectByName(t *testing.T) {
	type fields struct {
		Token string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		responseBody string
		want    int
		wantErr bool
	}{
		{
			name:         "Get project name",
			fields:       fields{"test"},
			args:         args{"Support"},
			responseBody: projectsResponse,
			want:         55,
			wantErr:      false,
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClubHouse{
				Token: tt.fields.Token,
			}
			httpmock.RegisterResponder("GET", ClubHouseAPIURL + "/api/v3/projects",
				httpmock.NewStringResponder(200, tt.responseBody))
			got, err := c.GetProjectByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProjectByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetProjectByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}