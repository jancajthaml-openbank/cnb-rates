ifndef GITHUB_RELEASE_TOKEN
$(warning GITHUB_RELEASE_TOKEN is not set)
endif

META := $(shell git rev-parse --abbrev-ref HEAD 2> /dev/null | sed 's:.*/::')
VERSION := $(shell git fetch --tags --force 2> /dev/null; tags=($$(git tag --sort=-v:refname)) && ([ $${\#tags[@]} -eq 0 ] && echo v0.0.0 || echo $${tags[0]}))

.ONESHELL:

.PHONY: all
all: bootstrap sync test package bbtest

.PHONY: package
package:
	@$(MAKE) bundle-binaries
	@$(MAKE) bundle-debian
	@$(MAKE) bundle-docker

.PHONY: bundle-binaries
bundle-binaries:
	@docker-compose run --rm package --arch linux/amd64 --pkg cnb-rates-import
	@docker-compose run --rm package --arch linux/amd64 --pkg cnb-rates-rest
	@docker-compose run --rm package --arch linux/amd64 --pkg cnb-rates-batch

.PHONY: bundle-debian
bundle-debian:
	@docker-compose run --rm debian -v $(VERSION)+$(META) --arch amd64

.PHONY: bundle-docker
bundle-docker:
	@docker build -t openbank/cnb-rates:$(VERSION)-$(META) .

.PHONY: bootstrap
bootstrap:
	@docker-compose build --force-rm go

.PHONY: lint
lint:
	@docker-compose run --rm lint --pkg cnb-rates-import || :
	@docker-compose run --rm lint --pkg cnb-rates-rest || :
	@docker-compose run --rm lint --pkg cnb-rates-batch || :

.PHONY: sec
sec:
	@docker-compose run --rm sec --pkg cnb-rates-import || :
	@docker-compose run --rm sec --pkg cnb-rates-rest || :
	@docker-compose run --rm sec --pkg cnb-rates-batch || :

.PHONY: sync
sync:
	@docker-compose run --rm sync --pkg cnb-rates-rest
	@docker-compose run --rm sync --pkg cnb-rates-import
	@docker-compose run --rm sync --pkg cnb-rates-batch

.PHONY: test
test:
	@docker-compose run --rm test --pkg cnb-rates-import
	@docker-compose run --rm test --pkg cnb-rates-rest
	@docker-compose run --rm test --pkg cnb-rates-batch

.PHONY: release
release:
	@docker-compose run --rm release -v $(VERSION)+$(META) -t ${GITHUB_RELEASE_TOKEN}

.PHONY: bbtest
bbtest:
	@(docker rm -f $$(docker ps -a --filter="name=cnb_rates_bbtest" -q) &> /dev/null || :)
	@docker exec -it $$(\
		docker run -d -ti \
			--name=cnb_rates_bbtest \
			-e UNIT_VERSION="$(VERSION)-$(META)" \
			-e UNIT_ARCH=amd64 \
			-v /sys/fs/cgroup:/sys/fs/cgroup:ro \
			-v /var/run/docker.sock:/var/run/docker.sock \
      -v /var/lib/docker/containers:/var/lib/docker/containers \
			-v $$(pwd)/bbtest:/opt/bbtest \
			-v $$(pwd)/reports:/reports \
			--privileged=true \
			--cap-add=SYS_TIME \
			--security-opt seccomp:unconfined \
		jancajthaml/bbtest:amd64 \
	) rspec --require /opt/bbtest/spec.rb \
		--format documentation \
		--format RspecJunitFormatter \
		--out junit.xml \
		--pattern /opt/bbtest/features/*.feature
	@(docker rm -f $$(docker ps -a --filter="name=cnb_rates_bbtest" -q) &> /dev/null || :)
