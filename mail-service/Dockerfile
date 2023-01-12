FROM golang:1.19.3-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o mail-service /app/main.go


FROM alpine

RUN mkdir /app

COPY --from=builder /app/mail-service /app

CMD [ "/app/mail-service" ]