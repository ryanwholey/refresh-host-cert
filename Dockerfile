FROM golang:alpine AS build

WORKDIR /src
COPY . ./

RUN go build

FROM alpine

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=build /src/stevedore-kubernetes ./

ENTRYPOINT ["./refresh-host-cert"]
