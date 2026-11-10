FROM golang:1.26.5-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/time-management-api .

FROM alpine:3.23

RUN apk add --no-cache ca-certificates tzdata
COPY --from=build /out/time-management-api /usr/local/bin/time-management-api

EXPOSE 8080
ENTRYPOINT ["time-management-api"]
