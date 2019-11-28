.PHONY: gomodgen deploy delete test coverage

AUTH_USER ?= ""
AUTH_PASSWORD ?= ""
CLUBHOUSE_STORY_TYPE ?= "chore"
CLUBHOUSE_PROJECT ?= "Support"
CLUBHOUSE_WORKFLOW ?= "Dev"
CLUBHOUSE_COMPLETED_STATE ?= "Completed"

require-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

gomodgen:
	GO111MODULE=on go mod init

deploy: require-CH_TOKEN require-GCP_PROJECTmak
	serverless deploy

test:
	go test

coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

delete:
	serverless remove
