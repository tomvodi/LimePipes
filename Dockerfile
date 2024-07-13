FROM golang:1.20-bullseye as builder

CMD mkdir -p /app
WORKDIR /app

COPY go.mod go.sum limepipes.env.default ./
COPY ./cmd ./cmd
COPY ./internal ./internal

RUN go mod download
RUN go build -o /limepipes limepipes/cmd/limepipes

FROM golang:1.20-bullseye

RUN mkdir -p /app
WORKDIR /app
RUN mkdir -p /opt/limepipes

RUN groupadd -g 1234 limepipes && useradd -r -u 1234 -g limepipes limepipes

COPY --from=builder /limepipes /app
COPY --from=builder /app/limepipes.env.default /app/limepipes.env

RUN chown limepipes:limepipes /opt/limepipes

USER limepipes

EXPOSE 8080
CMD ["/app/limepipes"]

