module github.com/journeymidnight/yig-front-caddy

go 1.12

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190426145343-a29dc8fdc734
	golang.org/x/net => github.com/golang/net v0.0.0-20181201002055-351d144fa1fc
	golang.org/x/sync => github.com/golang/sync v0.0.0-20181108010431-42b317875d0f
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190426135247-a129542de9ae
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190428024724-550556f78a90
	gopkg.in/alecthomas/kingpin.v2 => github.com/alecthomas/kingpin v2.2.6+incompatible
	gopkg.in/natefinch/lumberjack.v2 => github.com/natefinch/lumberjack v2.0.0+incompatible
	gopkg.in/square/go-jose.v2 => github.com/square/go-jose v2.3.1+incompatible
	gopkg.in/yaml.v2 => github.com/go-yaml/yaml v2.1.0+incompatible

)

require (
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da
	github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc
	github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf
	github.com/bifurcation/mint v0.0.0-20190129141059-83ba9bc2ead9 // indirect
	github.com/cep21/circuit v2.4.1+incompatible
	github.com/codahale/aesnicheck v0.0.0-20140702143623-349fcc471aac
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dustin/go-humanize v0.0.0-20170228161531-259d2a102b87
	github.com/flynn/go-shlex v0.0.0-20150515145356-3f9db97f8568
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/mock v1.2.0
	github.com/golang/protobuf v1.3.1
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v0.0.0-20170718202341-a69d9f6de432
	github.com/hashicorp/go-syslog v0.0.0-20170829120034-326bf4a7f709
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/jimstudt/http-authentication v0.0.0-20140401203705-3eca13d6893a
	github.com/klauspost/cpuid v0.0.0-20180102081000-ae832f27941a
	github.com/lucas-clemente/aes12 v0.0.0-20171027163421-cd47fb39b79f // indirect
	github.com/lucas-clemente/quic-go v0.11.1
	github.com/lucas-clemente/quic-go-certificates v0.0.0-20160823095156-d2f86524cced // indirect
	github.com/miekg/dns v0.0.0-20170721150254-0f3adef2e220
	github.com/naoina/go-stringutil v0.1.0
	github.com/naoina/toml v0.1.1
	github.com/prometheus/client_golang v1.0.0
	github.com/russross/blackfriday v0.0.0-20170610170232-067529f716f4
	github.com/xenolf/lego v0.0.0-20181204200439-4e842a5eb6dc
	golang.org/x/crypto v0.0.0-20190228161510-8dd112bcdc25
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3
	gopkg.in/natefinch/lumberjack.v2 v2.0.0-00010101000000-000000000000
	gopkg.in/square/go-jose.v2 v2.0.0-00010101000000-000000000000 // indirect
	gopkg.in/yaml.v2 v2.2.1
	zvelo.io/ttlru v1.0.9
)
