FROM golang:alpine AS build
RUN apk add git --no-cache
COPY . /go/src/refresh-host-cert
WORKDIR /go/src/refresh-host-cert

RUN go get refresh-host-cert

RUN go install

FROM alpine

RUN apk --no-cache add ca-certificates

COPY --from=build /go/bin/refresh-host-cert /refresh-host-cert

ENTRYPOINT ["/refresh-host-cert"]
