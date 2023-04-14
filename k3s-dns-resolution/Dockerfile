FROM arm64v8/golang:1.20-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o app



FROM arm64v8/golang:1.20-alpine

COPY --from=build /app/app app

ENTRYPOINT "app"

