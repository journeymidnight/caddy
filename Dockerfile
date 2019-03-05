FROM journeymidnight/go:1.12.4
WORKDIR /work
COPY . /work
RUN  yum --enablerepo=epel-testing install -y lttng-ust make gcc git rpm-build && \
cd caddy && go run build.go && \
cp ./caddy /work/conf/

CMD cd /work/conf/ && ./caddy