
clean:
	rm -rf pkg/gen

gen: clean
	openapi-generator generate -c config.yaml -g go-server -i forge_api_v3.json -o pkg/gen/v3
	# for f in pkg/gen/openapi/*.go; do goimports -w "$$f"; done || true
