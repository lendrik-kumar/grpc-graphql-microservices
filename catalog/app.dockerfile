FROM golang:1.23-alpine AS build

RUN apk --no-cache add gcc g++ make ca-certificates git

WORKDIR /go/src/github.com/lendrik-kumar/graphql-grpc-go-microservices

COPY go.mod go.sum ./
COPY catalog/ catalog/

RUN go build -o /go/bin/app ./catalog/cmd/catalog

FROM alpine:3.11
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["./app"]