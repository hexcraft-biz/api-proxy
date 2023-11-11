# api-proxy
API-PROXY authorizes incoming HTTP requests. It can be the Policy Enforcement Point in your cloud architecture, i.e. a reverse proxy in front of your upstream resource API server that rejects unauthorized requests and forwards authorized ones to your resource server. 

## Set up
```bash
Set up env file
$ cp ./.env.example ./.env
# Don't forget change OAuth2 related info at .env file.
$ cp ./example.proxyMappings.json ./.proxyMappings.json
# Don't forget add your resource services proxy rule to .proxyMappings.json
```

## Deployment
- Standalone Testing Flow
```bash
$ sh ./run.sh for testing
```
- Integration testing
  - Please flollow [hexc-deploy](https://github.com/hexcraft-biz/hexc-deploy) README.md step.
