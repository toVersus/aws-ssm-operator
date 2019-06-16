FROM alpine:3.8

RUN apk add --update --no-cache ca-certificates

USER nobody

ADD build/_output/bin/aws-ssm-operator /usr/local/bin/aws-ssm-operator
