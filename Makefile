.PHONY: build
build:
	@go build -o swrapper

.PHONY: enc
enc:
	./swrapper -m enc -p plain.sh -c plain.sh.enc

.PHONY: run
run:
	./swrapper -m run -c plain.sh.enc
