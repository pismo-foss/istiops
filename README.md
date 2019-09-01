# Istio Operator

Istio Operator (a.k.a `istiops`) is a tool to manage traffic for microservices deployed via [Istio](https://istio.io/). It simplifies deployment strategies such as bluegreen or canary releases with no need of messing around with tons of `yamls` from kubernetes' resources.

## Architecture

<img src="https://github.com/pismo/istiops/blob/master/imgs/overview.png">

## Running tests

`$ go test ./... -v`

## Building the CLI

To use istiops binary you can just `go build` it. It will generate a command line interface tool to work with.

`./run` or `go build -o build/istiops main.go`

You can then run it as: `./build/istiops version`

## Prerequisites

A kubernetes config at `~/.kube/config` which allows the binary to GET, UPDATE and LIST resources: `virtualservices` & `destinationrules`.
 If you are running the binary with a custom kubernetes' service account you can use this RBAC template to append to your roles:

```
- apiGroups: ["networking.istio.io"]
  resources: ["virtualservices", "destinationrules"]
  verbs: ["get", "list", "patch","update"]
  ````

## How it works ?

Istiops creates routing rules into virtualservices & destination rules in order to manage traffic correctly. This is an example of a routing being managed by Istio, using as default routing rule any HTTP request which matches as URI the regular expression: `'.+'`:

<img src="https://github.com/pismo/istiops/blob/master/imgs/howitworks1.png">

We call this `'.+'` rule as **master-route**, which it will be served as the default routing rule.


### Operator traffic steps

A deeper in the details

1. Find the needed kubernetes' resources based on given `labels-selector`

2. Create associate route rules based on `pod-selector` (to match which pods the routing will be served) & destination information (such as `hostname` and `port`)

3. Attach to an existent route rule a `request-headers` match if given

4. Attach to an existent route rule a `weight` if given


## Commands on traffic shifting

### Each operation creates or removes items from both the VirtualService and DestinationRule

1. Clear all traffic rules, except for **master-route** (default), from service api-domain

`istiops traffic clear --label-selector app=api-domain`

2. Send requests with HTTP header `"x-cid: seu_madruga"` to pods with labels `app=api-domain,build=PR-10`

```
$ istiops traffic headers \
    --namespace "default" \
    --hostname "api.domain.io" \
    --port 5000 \
    --label-selector "app=api-domain" \
    --pod-selector "app=api-domain,build=PR-10" \
    --headers "x-cid=seu_madruga" \
```

3. Send 20% of traffic to pods with labels `app=api-domain,build=PR-10`

```
$ istiops traffic weight \
    --namespace "default" \
    --hostname "api.domain.io" \
    --port 5000 \
    --label-selector "app=api-domain" \
    --pod-selector "app=api-domain,build=PR-10" \
    --weight 20 \
```

4. Removes all traffic (rollback), headers and percentage, for pods with labels: app=api-gateway,version:1.0.0

`istiops traffic rollback --label-selector --pod-selector app=api-domain,version:1.0.0`

## Importing as a package

You can assemble `istiops` as interface for your own Golang code, to do it you just have to initialize the needed struct-dependencies and call the interface directly. You can see examples at `./examples`