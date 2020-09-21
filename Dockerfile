FROM golang:1.15.2 AS builder

RUN mkdir -p $GOPATH/src/github.com/yoshiodeveloper/home-broker
ADD . $GOPATH/src/github.com/yoshiodeveloper/home-broker
WORKDIR $GOPATH/src/github.com/yoshiodeveloper/home-broker

# Build
ARG build
ARG version
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=${version} -X main.Build=${build}" -o /homebroker

# Final Step
FROM alpine:20200917

# Base packages
RUN apk update
RUN apk upgrade
RUN apk add --update tzdata
RUN rm -rf /var/cache/apk/*

# Copy binary from build step
COPY --from=builder /homebroker /home/

# Define timezone
ENV TZ=America/Sao_Paulo

# Define the ENTRYPOINT
WORKDIR /home
ENTRYPOINT ["./homebroker"]

# Document that the service listens on port 8080.
EXPOSE 8080 8081