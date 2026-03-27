FROM golang:1.25-alpine AS builder

ARG GOPROXY=https://goproxy.cn,direct

WORKDIR /app

ENV GOPROXY=${GOPROXY}

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o meaw-oa .

FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY --from=builder /app/meaw-oa .

EXPOSE 3000

ENTRYPOINT ["./meaw-oa"]
