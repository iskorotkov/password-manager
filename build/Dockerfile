FROM golang:1.17 AS build
WORKDIR /go/src

COPY ["go.mod", "go.sum", "./"]
RUN go get -d -v ./...

COPY . .
ENV CGO_ENABLED=0
RUN go build -o app cmd/passwordmanager/main.go

FROM alpine AS runtime
WORKDIR /

COPY --from=build /go/src/app ./
COPY --from=build /go/src/api ./api
COPY --from=build /go/src/static ./static
EXPOSE 8080/tcp
ENTRYPOINT ["./app"]
