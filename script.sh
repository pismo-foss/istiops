istiops traffic shift -b 666 -d api-istiops:8080 -e -H x-version=PR-133 -n ext -p fullname=api-istiops-xpto -l app=api-istiops
istiops traffic shift -b 777 -d api-istiops:8080 -e -H x-version=pod -n ext -p fullname=api-istiops-ext-non-pods -l app=api-istiops
istiops traffic shift -b 999 -d api-istiops:8080 -e -H x-version=2.0.0 -n ext -p fullname=api-istiops-ext-pr-193-1 -l app=api-istiops
istiops traffic show -l app=api-istiops -n ext
