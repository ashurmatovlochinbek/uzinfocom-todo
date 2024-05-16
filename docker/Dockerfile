FROM golang:1.21-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app cmd/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app .

EXPOSE 5050

USER nonroot:nonroot
CMD ["chmod", "+x", "app"]
CMD ["./app"]


