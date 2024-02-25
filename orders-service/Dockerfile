FROM golang:1.22-alpine AS build

WORKDIR /go/src/versta-2024/orders-service

COPY app ./app/

COPY go.mod go.sum ./

RUN go mod download

RUN go build -o ../../../bin/app ./app/cmd/app/main.go

FROM alpine
WORKDIR /go

COPY --from=build /go/bin/app /bin/app

CMD ["app"]