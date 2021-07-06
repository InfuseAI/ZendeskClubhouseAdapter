.PHONY: gomodgen deploy delete test coverage

FUNCTION_NAME ?= ZendeskClubhouseAdapter
CLUBHOUSE_STORY_TYPE ?= chore
CLUBHOUSE_PROJECT ?= Support
CLUBHOUSE_TEAM ?= Support
CLUBHOUSE_WORKFLOW ?= Support
CLUBHOUSE_PENDING_STATE ?= Suspend
CLUBHOUSE_COMPLETED_STATE ?= Resolved

require-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

gomodgen:
	GO111MODULE=on go mod init

deploy: require-CH_TOKEN require-GCP_PROJECT
	gcloud config set project $(GCP_PROJECT)
	gcloud functions deploy $(FUNCTION_NAME) --allow-unauthenticated --runtime=go111 --entry-point ZendeskClubhouseAdapter --trigger-http \
	--set-env-vars CH_TOKEN="$(CH_TOKEN)",AUTH_USER="$(AUTH_USER)",AUTH_PASSWORD="$(AUTH_PASSWORD)",CLUBHOUSE_STORY_TYPE="$(CLUBHOUSE_STORY_TYPE)",CLUBHOUSE_PROJECT="$(CLUBHOUSE_PROJECT)",CLUBHOUSE_WORKFLOW="$(CLUBHOUSE_WORKFLOW)",CLUBHOUSE_PENDING_STATE="$(CLUBHOUSE_PENDING_STATE)",CLUBHOUSE_COMPLETED_STATE="$(CLUBHOUSE_COMPLETED_STATE)"

test:
	go test

coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

delete:
	gcloud functions delete $(FUNCTION_NAME)
