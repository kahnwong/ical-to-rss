FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -ldflags "-w -s" -o /ical-to-rss

# hadolint ignore=DL3007
FROM gcr.io/distroless/static-debian13:latest AS deploy
COPY --from=build /ical-to-rss /

EXPOSE 3000
CMD ["/ical-to-rss"]
