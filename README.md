# Zendesk Cloubhouse Adpter
--
## Prepare
- Install Google Cloud SDK
```bash
curl https://sdk.cloud.google.com | bash

exec -l $SHELL

gcloud init
 
```

## How to deploy
```bash
make deploy GCP_PROJECT=<your-gcp-project-name> CH_TOKEN=<your-clubhouse-token> \
            [AUTH_USER=<http-auth-username>] [AUTH_PASSWORD=<http-auth-password>]
```

## How to run test
```bash
make test
# Code coverage
make coverage
```