FROM google/cloud-sdk:alpine
RUN apk add --update jq bash
COPY resource /opt/resource/