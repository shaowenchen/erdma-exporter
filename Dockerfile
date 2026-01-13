# Build stage
FROM golang:1.21 AS builder

WORKDIR /build

RUN apt-get update && apt-get install -y git make && rm -rf /var/lib/apt/lists/*

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo -o erdma-exporter .

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    ca-certificates \
    wget \
    gnupg \
    && rm -rf /var/lib/apt/lists/*

RUN wget -qO - http://mirrors.aliyun.com/erdma/GPGKEY | gpg --dearmour -o /etc/apt/trusted.gpg.d/erdma.gpg && \
    echo "deb [ ] http://mirrors.aliyun.com/erdma/apt/ubuntu jammy/erdma main" | tee /etc/apt/sources.list.d/erdma.list && \
    apt-get update && \
    apt-get install -y \
    libibverbs1 \
    eadm \
    ibverbs-providers \
    ibverbs-utils \
    librdmacm1 \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /root

COPY --from=builder /build/erdma-exporter .

EXPOSE 9101

ENTRYPOINT ["./erdma-exporter"]

