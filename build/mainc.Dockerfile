FROM golang:1.21.6-alpine AS builder

COPY . /github.com/Andrest31/Coding
WORKDIR /github.com/Andrest31/Coding

RUN go mod download
RUN go clean --modcache
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./main.go

FROM scratch AS runner

WORKDIR /build_v1/

COPY --from=builder /github.com/Andrest31/Coding/.bin .

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip

EXPOSE 8020

ENTRYPOINT ["./.bin"]