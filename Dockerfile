FROM golang:latest as build-stage
WORKDIR /core
COPY go.mod .

COPY . .
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk add ca-certificates

COPY --from=build-stage /core/main /
ENTRYPOINT ["/main"]
