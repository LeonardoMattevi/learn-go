BINARY   := druid2
WEB_DIR  := web
WASM_OUT := $(WEB_DIR)/main.wasm
WASM_JS  := $(WEB_DIR)/wasm_exec.js

.PHONY: build wasm test serve deploy clean

build:
	go build -o $(BINARY) .

wasm:
	GOOS=js GOARCH=wasm go build -o $(WASM_OUT) .
	find "$$(go env GOROOT)" -name "wasm_exec.js" -exec cp {} $(WASM_JS) \;

test:
	go test ./... -cover

serve: wasm
	python3 -m http.server --directory $(WEB_DIR) 8080

deploy: wasm
	npx wrangler pages deploy $(WEB_DIR)/ --project-name go-game

clean:
	rm -f $(BINARY) $(WASM_OUT) $(WASM_JS)
