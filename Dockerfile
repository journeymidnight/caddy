FROM journeymidnight/go:1.12.4
WORKDIR /work
RUN go env
COPY . /work
RUN  export GOPROXY=https://goproxy.cn && \
cd caddy && go run build.go
CMD cd /work/caddy/ && ./caddy