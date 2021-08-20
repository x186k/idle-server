##
## Build
##

FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /docker-gs-ping

##
## Deploy
##

FROM restreamio/gstreamer:1.18.4.0-prod

WORKDIR /

COPY --from=build /docker-gs-ping /docker-gs-ping

EXPOSE 8080

#USER nonroot:nonroot

ENTRYPOINT ["/docker-gs-ping"]
