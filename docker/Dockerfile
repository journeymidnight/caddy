FROM journeymidnight/yig
ENV GOPATH="/root/go"

COPY . /work
WORKDIR /work

RUN rpm --rebuilddb && yum --enablerepo=epel-testing install -y golang && \
go get -v github.com/journeymidnight/yig-front-caddy && \
go get -v github.com/caddyserver/builds && \
go get -v github.com/prometheus/client_golang/prometheus && \
cd /root/go/src/github.com/journeymidnight/yig-front-caddy/caddy && \
go run build.go && \
cp ./caddy /work/docker/ && \
yum clean all && \
rm -rf /var/cache/yum

CMD cd /work/docker && ./caddy
