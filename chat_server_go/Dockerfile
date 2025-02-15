# syntax=docker/dockerfile:1

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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -buildvcs=true -a -gcflags='all="-l"' -ldflags='-s -w -extldflags "-static"'  -o /app/webserver cmd/webserver/main.go

# use debug because it includes busybox
FROM gcr.io/distroless/static-debian11:debug-nonroot@sha256:55716e80a7d4320ce9bc2dc8636fc193b418638041b817cf3306696bd0f975d1
LABEL org.opencontainers.image.url='https://github.com/MKand/movie-guru/tree/cloud-movieguru' \
    org.opencontainers.image.documentation='https://github.com/MKand/movie-guru/tree/cloud-movieguru' \
    org.opencontainers.image.source='https://github.com/MKand/movie-guru/tree/cloud-movieguru' \
    org.opencontainers.image.vendor='Google LLC' \
    org.opencontainers.image.licenses='Apache-2.0'
COPY --from=builder /app/webserver /usr/local/bin/webserver    

ENTRYPOINT [ "webserver" ]