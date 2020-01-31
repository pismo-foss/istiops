#!/usr/bin/env bash
for svc in $(kubectl get svc -n $env | awk '{print $1}'); do /run/istiops traffic clear -l app=$svc -n $env; done;
