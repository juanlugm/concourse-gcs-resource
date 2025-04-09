FROM google/cloud-sdk:alpine
RUN apk add --update jq yq bash
COPY resource /opt/resource/