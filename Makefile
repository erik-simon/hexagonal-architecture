builds:
	bash ./scripts/build.bash

run-http: builds
	./builds/http

run-cli: builds
	./builds/cli