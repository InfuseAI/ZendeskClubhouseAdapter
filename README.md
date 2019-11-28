# Zendesk Cloubhouse Adpter
--
## Prepare
- Install serverless framework
```bash
npm install -g serverless 
```

## How to deploy
```bash
make deploy GCP_PROJECT=<your-gcp-project-name> [AUTH_USER=<http-auth-username>] [AUTH_PASSWORD=<http-auth-password>]
```

## How to run test
```bash
make test
# Code coverage
make coverage
```