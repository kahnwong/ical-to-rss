FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -ldflags "-w -s" -o /ical-to-rss

# hadolint ignore=DL3007
FROM gcr.io/distroless/static-debian11:latest AS deploy
COPY --from=build /ical-to-rss /

EXPOSE 3000
CMD ["/ical-to-rss"]
