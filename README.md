# Pipeline 3.0

## Running tests

`$ go test ./... -v`

## Commands on traffic shifiting

### Each operation creates or removes items from both the VirtualService and DestinationRule

Clear all traffic rules, except for main one, for service api-gateway

`istiops traffic --clear app=api-gateway`

Send requests with HTTP header x-cid:seu_madruga to pods with labels: app=api-accounts,build:PR-10

`istiops traffic --headers x-cid=seu_madruga --to app=api-accounts,build:PR-10`

Send 10% of traffic to pods with labels: app=api-gateway,version:1.0.0

`istiops traffic --percentage 10 --to app=api-gateway,version=1.0.0`

Removes all traffic (rollback), headers and percentage, for pods with labels: app=api-gateway,version:1.0.0

`istiops traffic --clear app=api-gateway,version:1.0.0`
