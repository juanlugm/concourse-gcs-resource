FROM google/cloud-sdk
RUN apk add --update jq bash
COPY resource /opt/resource/