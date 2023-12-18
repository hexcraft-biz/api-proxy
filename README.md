# Drawbridge
Proxy authorizes incoming HTTP requests. It can be the Policy Enforcement Point in your cloud architecture, i.e. a reverse proxy in front of your upstream resource API server that rejects unauthorized requests and forwards authorized ones to your resource server. 

## Set up
```bash
Set up env file
$ cp ./.env.example ./.env
# Don't forget change OAuth2 & Dogmas related info at .env file.
```

## Deployment
- Standalone Testing Flow with docker-compose
```bash
$ docker-compose -f dev.yml up --build -d
```
- Integration testing
  - Please flollow [hexc-deploy](https://github.com/hexcraft-biz/hexc-deploy) README.md step.

## API Endpoint
### HealthCheck
#### GET /healthcheck/v1/ping
- Params
  - None
- Response
  - 200
	```json
	{
	  "message": "OK"
	}
	```
