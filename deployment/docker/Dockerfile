FROM golang:21-bookworm AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./
RUN mkdir
RUN CGO_ENABLED=0 go build -ldflags "-w -s" -o ./app

FROM golang:21-bookworm AS runner

RUN mkdir -p /app
COPY  ./configs/configs.env /app/configs.env
COPY --from=builder /app/app /app/app
WORKDIR /app
EXPOSE 4200
CMD ["./app"]