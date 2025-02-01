include vendor.proto.mk

# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN := $(CURDIR)/bin

BUF_BIN := $(LOCAL_BIN)/buf

# устанавливаем необходимые плагины
.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps:
	$(info Installing binary dependencies...)

	go install github.com/bufbuild/buf/cmd/buf@v1.41.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Форматировние .proto файлов с помощью buf
.buf-format:
	$(info run buf format...)
	PATH="$(LOCAL_BIN):$(PATH)" $(BUF_BIN) format -w

# Линтинг .proto файлов с помощью buf
.buf-lint:
	$(info run buf lint...)
	PATH="$(LOCAL_BIN):$(PATH)" $(BUF_BIN) lint

# Генерация .pb файлов с помощью buf
.buf-generate:
	$(info run buf generate...)
	PATH="$(LOCAL_BIN):$(PATH)" $(BUF_BIN) generate

.tidy:
	go mod tidy

# Генерация .pb файлов
generate: .buf-format .buf-lint .buf-generate .tidy

# Объявляем, что текущие команды не являются файлами и
# интсрументируем Makefile не искать изменения в файловой системе
.PHONY: \
	.bin-deps \
	.buf-format \
	.buf-lint \
	.buf-generate \
	.tidy \
	generate
