# builder
FROM golang:alpine as builder

RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o main .

# runner
FROM alpine

ENV PORT 8080
ENV MONGO_URL "mongodb://127.0.0.1:27017"
ENV MONGO_COLLECTION "acts-dev"
ENV GIN_MODE=release

RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
WORKDIR /app

EXPOSE $PORT
ENTRYPOINT ./main \
    -mongo-url=$MONGO_URL \
    -mongo-collection=$MONGO_COLLECTION \
    -port=$PORT