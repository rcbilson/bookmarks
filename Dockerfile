FROM golang:1.24-alpine AS build-server
RUN apk update && apk add gcc libc-dev
WORKDIR /src
COPY backend/go.mod backend/go.sum .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \ 
    go mod download && go mod verify
COPY backend .
# sqlite requires cgo
ARG CGO_ENABLED=1
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \ 
    go build --tags fts5 -o /bin/server ./cmd/server

FROM node:23-bullseye AS build-frontend
WORKDIR /src
COPY frontend/package.json frontend/yarn.lock .
RUN yarn config set network-timeout 300000
RUN yarn install
COPY frontend .
RUN yarnpkg run build

FROM alpine:latest
COPY --from=build-frontend /src/dist /app/frontend
COPY --from=build-server /bin /app/bin
ENV BOOKMARKSERVER_DBFILE=/app/data/bookmark.db
ENV BOOKMARKSERVER_FRONTENDPATH=/app/frontend
ENV BOOKMARKSERVER_PORT=9093
CMD ["/app/bin/server"]
