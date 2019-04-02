FROM golang:1.11 as builder

WORKDIR /src

COPY . .

RUN go build -o build/manabie cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /src/
COPY --from=builder /src/build/manabie .
COPY --from=builder /src/resources/ ./resources
CMD ["./manabie"]


