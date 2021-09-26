FROM golang:1.17.1 as builder
WORKDIR /app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/ddosify main.go


FROM alpine:3.14.2
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/ddosify ./bin/
ENTRYPOINT ["ddosify"]  
