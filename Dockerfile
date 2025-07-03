# syntax=docker/dockerfile:1

# Build stage
FROM node:20-alpine AS node_builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install
COPY . .
RUN npm run build

FROM golang:1.24-alpine AS go_builder
WORKDIR /app
COPY --from=node_builder /app/dist ./dist
COPY . .
RUN go build -o /emoji-match

# Final stage
FROM scratch
WORKDIR /
COPY --from=go_builder /emoji-match /emoji-match
EXPOSE 3001
USER 10001
ENTRYPOINT ["/emoji-match"]
