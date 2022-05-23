# api-proxy
API-PROXY authorizes incoming HTTP requests. It can be the Policy Enforcement Point in your cloud architecture, i.e. a reverse proxy in front of your upstream resource API server that rejects unauthorized requests and forwards authorized ones to your resource server. 

```bash
1. $ cp ./.example.env ./.env
2. $ cp ./example.proxyMappings.json ./.proxyMappings.json and setup your proxy rule.
3. $ sh ./run.sh for testing
```
