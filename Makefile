.PHONY: gomodgen deploy delete

gomodgen:
	GO111MODULE=on go mod init

deploy:
	serverless deploy

delete:
	serverless remove
