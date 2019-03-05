FROM journeymidnight/yig:latest

COPY . /work
WORKDIR /work

RUN git clone https://github.com/journeymidnight/yig-front-caddy && \
cd yig-front-caddy/caddy && \
go run build.go && \
cp ./caddy /work/conf

CMD cd /work/conf && ./caddy
