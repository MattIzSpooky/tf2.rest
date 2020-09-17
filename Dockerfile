FROM golang:alpine as builder

WORKDIR /app

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ADD . .

RUN go build

FROM alpine AS runner

EXPOSE 8080

COPY --from=builder /app/tf2.rest .

ENV PORT=8080

ENTRYPOINT ./tf2.rest