FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOSUMDB=off

# Move to working directory /build
WORKDIR /build

RUN apk add --no-cache git

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . /build

# Build the application
RUN go build -o main /build/cmd/

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Build a small image
FROM alpine

# for health check
RUN apk --no-cache add curl

COPY --from=builder /dist/main /

ADD ./scripts/run.sh /run.sh

ENV GIT_BRANCH="main"
ENV SERVICE_NAME="default"

CMD /run.sh /main /var/log/solar_$SERVICE_NAME.$GIT_BRANCH.log
