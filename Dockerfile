FROM public.ecr.aws/docker/library/golang:1.20.11-alpine3.17 as builder

WORKDIR /
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o main ./cmd/main.go

FROM scratch

COPY --from=builder /main /main
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

CMD ["/main"]
