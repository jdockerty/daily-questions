FROM golang:1.17 as build

WORKDIR /build

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o dailyq

FROM scratch

WORKDIR /app

COPY --from=build /build/dailyq /app/dailyq
COPY --from=build /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT [ "/app/dailyq" ]