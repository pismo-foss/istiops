# Istio Traffic Shifter

Istio Traffic Shifter (a.k.a `istiops`) is a tool to manage traffic for microservices deployed via [Istio](https://istio.io/). It simplifies deployment strategies such as bluegreen or canary releases with no need of messing around with tons of `yamls` from kubernetes' resources.

<p align="center">
    <a href="https://asciinema.org/a/OHOd98DRgrBCAib8mUwWQptwh?autoplay=1"><img src="https://asciinema.org/a/OHOd98DRgrBCAib8mUwWQptwh.png" width="600"/></a>
</p>

## Documentation

* [Architecture](#architecure)
* [Running tests](#running-tests)
* [Building the CLI](#building-the-cli)
* [Prerequisites](#prerequisites)
* [How it works ?](#how-it-works-?)
    - [Traffic Shifting](#traffic-shifting)
* [Using CLI](#using-cli)
    - [Get current routes](#get-current-routes)
    - [Clear all routes](#clear-all-routes)
    - [Headers routing](#shift-to-request-headers-routing)
    - [Weight Routing](#shift-to-weight-routing)
* [Importing as a package](#importing-as-a-package)

## Architecture

<img src="https://github.com/pismo/istiops/blob/master/imgs/overview.png" alt="">

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

```sh
- apiGroups: ["networking.istio.io"]
  resources: ["virtualservices", "destinationrules"]
  verbs: ["get", "list", "patch","update"]
  ````

## How it works ?

Istiops creates routing rules into virtualservices & destination rules in order to manage traffic correctly. This is an example of a routing being managed by Istio, using as default routing rule any HTTP request which matches as URI the regular expression: `'.+'`:

<img src="https://github.com/pismo/istiops/blob/master/imgs/howitworks1.png" alt="">

We call this `'.+'` rule as **master-route**, which it will be served as the default routing rule.


### Traffic Shifting

A deeper in the details

1. Find the needed kubernetes' resources based on given `labels-selector`

2. Create associate route rules based on `pod-selector` (to match which pods the routing will be served) & destination information (such as `hostname` and `port`)

3. Attach to an existent route rule a `request-headers` match if given

<img src="https://github.com/pismo/istiops/blob/master/imgs/howitworks2.png" alt="">

4. Attach to an existent route rule a `weight` if given. In case of a `weight: 100` the balance-routing will be skipped.

<img src="https://github.com/pismo/istiops/blob/master/imgs/howitworks3.png" alt="">

## Using CLI

### Each operation list, creates or removes items from both the VirtualService and DestinationRule

### Get current routes

Get all current traffic rules (respecting routes order) for resources which matches `label-selector`

```bash
$ istiops traffic show \
    --label-selector environment=pipeline-go \
    --namespace default \
    --output beautified
```

Ex.

```bash
>
client -> request to ->  [api.domain.io]
  \_ Headers
      - x-email: contact@domain.io
      - x-version: PR-142
      \_ Destination [k8s service]
         - api-domain:5000

  \_ regex:".+"
      \_ Destination [k8s service]
         - api.domain.io:5000
           \_ 90 % of requests
               \_ pods with labels
                  - build: 24
                  - version: 1.3.3
                  - app: api
                ---
         - api.domain.io:5000
           \_ 10 % of requests
               \_ pods with labels
                  - version: 1.3.3
                  - app: api
                  - build: 24
```

### Clear all routes

2. Clear all traffic rules, except for **master-route** (default), from service api-domain

`$ istiops traffic clear -l app=api-domain`

### Shift to request-headers routing

3. Send requests with HTTP header `"x-cid: seu_madruga"` to pods with labels `app=api-domain,build=PR-10`

```bash
$ istiops traffic shift \
    --namespace "default" \
    --hostname "api.domain.io" \
    --port 5000 \
    --label-selector "app=api-domain" \
    --pod-selector "app=api-domain,build=PR-10" \
    --headers "x-cid=seu_madruga"
```

### Shift to weight routing
4. Send 20% of traffic to pods with labels `app=api-domain,build=PR-10`

```bash
$ istiops traffic shift \
    --namespace "default" \
    --hostname "api.domain.io" \
    --port 5000 \
    --label-selector "app=api-domain" \
    --pod-selector "app=api-domain,build=PR-10" \
    --weight 20
```

## Importing as a package

You can assemble `istiops` as an interface for your own Golang code, to do it you just have to initialize the needed struct-dependencies and call the interface directly. You can see proper examples at `./examples`