FROM alpine
RUN apk add --update bash && apk add --update curl && rm -rf /var/cache/apk/*

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
RUN mv kubectl /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/kubectl

COPY istiops /run/istiops
COPY clear_unused_istio_routes.sh /run/clear_unused_istio_routes.sh
RUN chmod +x /run/clear_unused_istio_routes.sh
RUN chmod +x /run/istiops
CMD ["bash", "/run/clear_unused_istio_routes.sh"]



