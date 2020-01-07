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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tomolink 

# final stage
FROM scratch
COPY --from=builder /app/tomolink /app/cmd/tomolink
EXPOSE 8080
# If you need default env vars for your environment: 
# ENV KEY1=VALUE1,KEY2=VALUE2
ENTRYPOINT ["/app/tomolink"]
