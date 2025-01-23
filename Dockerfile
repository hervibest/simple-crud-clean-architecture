FROM docker.io/golang:1.23.4 AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./out/dist ./cmd

FROM gcr.io/distroless/static-debian12
COPY --from=build /app/out/dist /app/dist
COPY .env /app/.env
WORKDIR /app
# Set entrypoint untuk aplikasi
CMD ["/app/dist"]
