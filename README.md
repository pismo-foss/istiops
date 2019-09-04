# Istio Operator

Istio Process Status (a.k.a `istiops`) is a tool to manage traffic for microservices deployed via [Istio](https://istio.io/). It simplifies deployment strategies such as bluegreen or canary releases with no need of messing around with tons of `yamls` from kubernetes' resources.

## Documentation

* [Architecture](#architecure)
* [Running tests](#running-tests)
* [Building the CLI](#building-the-cli)
* [Prerequisites](#prerequisites)
* [How it works ?](#how-it-works-?)
    - [Traffic Shifting](#traffic-shifting)
* [Using CLI](#using-cli)
* [Importing as a package](#importing-as-a-package)

## Architecture

<img src="https://github.com/pismo/istiops/blob/master/imgs/overview.png">

## Running tests

`$ go test ./... -v`

## Building the CLI

To use istiops binary you can just `go build` it. It will generate a command line interface tool to work with.

`./run` or `go get && build -o build/istiops main.go`

You can then run it as: `./build/istiops version`

## Prerequisites

- `go` version `1.12.9`+
- A kubernetes config at `~/.kube/config` which allows the binary to `GET`, `PATCH`, `UPDATE` and `LIST` resources: `virtualservices` & `destinationrules`.
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


### Traffic Shifting

A deeper in the details

1. Find the needed kubernetes' resources based on given `labels-selector`

2. Create associate route rules based on `pod-selector` (to match which pods the routing will be served) & destination information (such as `hostname` and `port`)

3. Attach to an existent route rule a `request-headers` match if given

<img src="https://github.com/pismo/istiops/blob/master/imgs/howitworks2.png">

4. Attach to an existent route rule a `weight` if given. In case of a `weight: 100` the balance-routing will be skipped.

<img src="https://github.com/pismo/istiops/blob/master/imgs/howitworks3.png">

## Using CLI

### Each operation list, creates or removes items from both the VirtualService and DestinationRule

### Get current routes

Get all current traffic rules for resources which matches `label-selector`

```shell script
istiops traffic show \
    --label-selector environment=pipeline-go \
    --namespace default \
    --output beautified
```

Ex.

```
>
api-domain-virtualservice
Hosts:  [api.domain.io]
* Match
  \_ Headers
      - x-email: contact@domain.io
      - x-version: PR-142
      \_ Destination
         - api-xpto-4-default
* Match
  \_ regex:".+" 
      \_ Destination
         - api-xpto-1-default
             \_ 90 % of requests
         - api-xpto-2-default
             \_ 10 % of requests
```

2. Clear all traffic rules, except for **master-route** (default), from service api-domain

`istiops traffic clear -l app=api-domain`

3. Send requests with HTTP header `"x-cid: seu_madruga"` to pods with labels `app=api-domain,build=PR-10`

```
$ istiops traffic shift \
    --namespace "default" \
    --hostname "api.domain.io" \
    --port 5000 \
    --label-selector "app=api-domain" \
    --pod-selector "app=api-domain,build=PR-10" \
    --headers "x-cid=seu_madruga" \
```

4. Send 20% of traffic to pods with labels `app=api-domain,build=PR-10`

```
$ istiops traffic shift \
    --namespace "default" \
    --hostname "api.domain.io" \
    --port 5000 \
    --label-selector "app=api-domain" \
    --pod-selector "app=api-domain,build=PR-10" \
    --weight 20 \
```

5. Removes all traffic (rollback), headers and percentage, for pods with labels: app=api-gateway,version:1.0.0

`istiops traffic rollback --label-selector --pod-selector app=api-domain,version:1.0.0`

## Importing as a package

You can assemble `istiops` as an interface for your own Golang code, to do it you just have to initialize the needed struct-dependencies and call the interface directly. You can see proper examples at `./examples`