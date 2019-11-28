.PHONY: gomodgen deploy delete test coverage

require-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

gomodgen:
	GO111MODULE=on go mod init

deploy: require-CH_TOKEN require-GCP_PROJECT
	serverless deploy

test:
	go test

coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

delete:
	serverless remove
