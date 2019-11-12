.PHONY: gomodgen deploy delete

require-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

gomodgen:
	GO111MODULE=on go mod init

deploy: require-CH_TOKEN
	serverless deploy

delete:
	serverless remove
