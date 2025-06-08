FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY main.go go.mod ./

RUN CGO_ENABLED=0 go build -o /action -trimpath -ldflags="-s -w" .

# ---

FROM gcr.io/distroless/static-debian12

COPY --from=builder /action /action

CMD ["/action"]
