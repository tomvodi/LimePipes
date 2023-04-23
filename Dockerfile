FROM golang:1.20-bullseye as builder

CMD mkdir -p /app
WORKDIR /app

COPY go.mod go.sum banduslib.env.default ./
COPY ./cmd ./cmd
COPY ./internal ./internal

RUN go mod download
RUN go build -o /banduslib banduslib/cmd/banduslib

FROM golang:1.20-bullseye

RUN mkdir -p /app
WORKDIR /app
RUN mkdir -p /opt/banduslib

RUN groupadd -g 1234 banduslib && useradd -r -u 1234 -g banduslib banduslib

COPY --from=builder /banduslib /app
COPY --from=builder /app/banduslib.env.default /app/banduslib.env

RUN chown banduslib:banduslib /opt/banduslib

USER banduslib

EXPOSE 8080
CMD ["/app/banduslib"]

