# build stage
FROM golang:alpine AS build-stage

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /go/bin/app .

# run stage
FROM alpine:latest as run-stage

WORKDIR /
COPY --from=build-stage /go/bin/app /spotify-weekly-backup
CMD ["/spotify-weekly-backup"]
