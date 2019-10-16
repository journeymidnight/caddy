# ~/bin/bash
export CERTNAME=www.wojiushixiangshiyishi.xyz
mkdir /etc/caddy
cd /etc/caddy
go run github.com/FiloSottile/mkcert -install $CERTNAME
