# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# build stage
FROM golang:1.13 as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tomolink cmd/httpserver.go

# final stage
#FROM gcr.io/distroless/static:nonroot
FROM debian
#COPY --from=builder --chown=nonroot /app/tomolink /app/
#COPY --chown=nonroot internal/config/tomolink_defaults.yaml /app/ 
COPY --from=builder /app/tomolink /app/
COPY internal/config/tomolink_defaults.yaml /app/ 
EXPOSE 8080

RUN ls -lahR

ENTRYPOINT ["/app/tomolink"]

# Docker Image Arguments
ARG BUILD_DATE
ARG VCS_REF
ARG BUILD_VERSION

MAINTAINER joeholley@google.com
# Standardized Docker Image Labels
# https://github.com/opencontainers/image-spec/blob/master/annotations.md
LABEL \
    org.opencontainers.image.created="${BUILD_TIME}" \
    org.opencontainers.image.authors="Google LLC <joeholley@google.com>" \
    org.opencontainers.image.url="https://github.com/joeholley/tomolink" \
    org.opencontainers.image.documentation="https://godoc.org/github.com/joeholley/tomolink" \
    org.opencontainers.image.source="https://github.com/joeholley/tomolink/README.md" \
    org.opencontainers.image.version="${BUILD_VERSION}" \
    org.opencontainers.image.revision="1" \
    org.opencontainers.image.vendor="Google LLC" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.ref.name="" \
    org.opencontainers.image.title="${IMAGE_TITLE}" \
    org.opencontainers.image.description="GCP-native gaming friends service" \
    org.label-schema.schema-version="1.0" \
    org.label-schema.build-date=$BUILD_DATE \
    org.label-schema.url="" \
    org.label-schema.vcs-url="https://github.com/joeholley/tomolink" \
    org.label-schema.version=$BUILD_VERSION \
    org.label-schema.vcs-ref=$VCS_REF \
    org.label-schema.vendor="Google LLC" \
    org.label-schema.name="${IMAGE_TITLE}" \
    org.label-schema.description="GCP-native gaming friends service" \
    org.label-schema.usage="https://github.com/joeholley/tomolink/README.md"
