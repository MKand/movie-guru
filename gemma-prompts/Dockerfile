# syntax=docker/dockerfile:1
# Copyright 2025 Google LLC
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


FROM golang:1.23.2@sha256:a7f2fc9834049c1f5df787690026a53738e55fc097cd8a4a93faa3e06c67ee32 AS builder

# Set destination for COPY
WORKDIR /app

COPY go.mod go.sum ./

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . .

ENV GO111MODULE=on
RUN go mod tidy
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -buildvcs=true -a -gcflags='all="-l"' -ldflags='-s -w -extldflags "-static"'  -o /app/bin/gemma-prompts main.go

# use debug because it includes busybox
FROM gcr.io/distroless/static-debian11:debug-nonroot@sha256:55716e80a7d4320ce9bc2dc8636fc193b418638041b817cf3306696bd0f975d1
LABEL org.opencontainers.image.url='https://github.com/srinandan/movie-guru' \
    org.opencontainers.image.documentation='https://github.com/srinandan/movie-guru' \
    org.opencontainers.image.source='https://github.com/srinandan/movie-guru' \
    org.opencontainers.image.vendor='Google LLC' \
    org.opencontainers.image.licenses='Apache-2.0'
COPY --from=builder /app/bin/gemma-prompts /usr/local/bin/gemma-prompts    

ENTRYPOINT [ "gemma-prompts" ]
