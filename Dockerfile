ARG GO_VERSION=1.14

FROM golang:${GO_VERSION}-alpine AS build

RUN apk add --no-cache git

ARG TAG_NAME
ARG SHORT_SHA

WORKDIR /src
COPY . .

ENV CGO_ENABLED=0 GOOS=linux

RUN go build -mod vendor -o /app \
    -ldflags "-X main.version=$TAG_NAME -X main.commit=$SHORT_SHA" \
    ./cmd/machinist

FROM gcr.io/distroless/base AS run

ADD https://pki.google.com/roots.pem /roots.pem
ENV MACHINIST_CERTS /roots.pem

COPY --from=build /app /machinist

CMD ["/machinist"]
