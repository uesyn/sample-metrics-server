FROM golang:1.15 as builder
COPY go.mod /app/go.mod
COPY go.sum /app/go.sum
WORKDIR /app
RUN go mod download
COPY . /app
RUN go build -o /sample-metrics-server

FROM gcr.io/distroless/base@sha256:0c5d357a80ab1315ef55f05be174a82e10d09fb2fd7dfcc3c44ebdde6f10c51e
COPY --from=builder /sample-metrics-server /bin/sample-metrics-server
ENTRYPOINT ["/bin/sample-metrics-server"]
