.PHONY: build

build:
	cd caddy && go run build.go && cd ..

pkg:
	docker run --rm -v ${PWD}:/work -w /work journeymidnight/caddy bash -c 'bash package/rpmbuild.sh'

image:
	docker build -t  journeymidnight/caddy . -f Dockerfile
