############################
# STEP 1 build executable binary
############################
FROM golang as builder

WORKDIR $GOPATH/src/aarnaud/gke-maintenance-state-exporter/
COPY . .

RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/gke-maintenance-state-exporter -mod vendor main.go

############################
# STEP 2 build a small image
############################
FROM scratch

ENV GIN_MODE=release
WORKDIR /app/
# Import from builder.
COPY --from=builder /go/bin/gke-maintenance-state-exporter /app/gke-maintenance-state-exporter
ENTRYPOINT ["/app/gke-maintenance-state-exporter"]
EXPOSE 9723
