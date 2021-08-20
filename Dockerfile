##
## Build
##

FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /idle-server

##
## Deploy
##

FROM restreamio/gstreamer:1.18.4.0-prod

WORKDIR /

COPY --from=build /idle-server /idle-server

EXPOSE 8080

#USER nonroot:nonroot

ENTRYPOINT ["/idle-server"]
