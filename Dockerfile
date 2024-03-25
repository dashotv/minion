############################
# STEP 1a build ui
############################
FROM node:20-alpine as ui-builder

WORKDIR /app/ui
COPY ui/package.json ui/yarn.lock ui/.yarn ./
RUN  --mount=type=cache,target=/app/ui/node_modules yarn install
COPY ui/ ./
RUN --mount=type=cache,target=/app/ui/node_modules yarn build

############################
# STEP 1b build go binary
############################
FROM golang:1.21-alpine AS builder

WORKDIR /go/src/app
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download

COPY . .
COPY --from=ui-builder /app/static ./static

RUN --mount=type=cache,target=/go/pkg/mod \
  go install ./server

############################
# STEP 2 build a small image
############################
FROM alpine
# Copy our static executable.
WORKDIR /root/
COPY --from=builder /go/bin/server .
COPY .env.vault .
CMD ["./server", "server"]
