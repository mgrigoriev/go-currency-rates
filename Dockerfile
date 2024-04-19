FROM golang:1.22-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN go build -o ./bin/app cmd/main.go

FROM alpine

ENV TEMPLATES_DIR=/app/templates
RUN mkdir -p /app/cmd /app/templates
COPY --from=builder /usr/local/src/bin/app /app/cmd/
COPY --from=builder /usr/local/src/templates /app/templates/

CMD ["/app/cmd/app"]
