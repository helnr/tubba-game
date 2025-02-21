FROM node:23.1.0 AS frontend-builder

WORKDIR /app

COPY frontend/package*.json ./

RUN npm install

COPY frontend/ ./

ENV VITE_API_URL=/api
ENV VITE_API_GAME_URL=/api/game/join

RUN npm run build

FROM golang:1.23.5 AS backend-builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./

RUN go mod download

COPY backend/ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server main.go

FROM alpine:3 AS certs
RUN apk --no-cache add ca-certificates

FROM scratch

WORKDIR /app

COPY --from=backend-builder /app/server ./

COPY --from=frontend-builder /app/dist ./dist

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080

CMD ["/app/server"]