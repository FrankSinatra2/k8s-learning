FROM arm64v8/golang:1.19.8-alpine AS build

WORKDIR /usr/src/app


COPY . .
RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]




