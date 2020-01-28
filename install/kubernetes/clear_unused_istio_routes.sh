#!/usr/bin/env bash
#
kubectl -n ext get services
/run/istiops traffic show -l app=sec-bankaccounts -n ext -o beauty