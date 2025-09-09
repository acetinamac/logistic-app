# Variables
AIR_CONF := .air.toml
BIN := ./tmp/main.exe
MAIN := cmd/api/main.go

# Target por defecto
.DEFAULT_GOAL := dev

.PHONY: dev swag build run clean deps tools help

help:
	@echo "Targets disponibles:"
	@echo "  dev    - Inicia el servidor con Air usando $(AIR_CONF)"
	@echo "  swag   - Genera la documentación Swagger (swag init)"
	@echo "  build  - Compila el binario en $(BIN)"
	@echo "  run    - Ejecuta el binario compilado"
	@echo "  clean  - Limpia artefactos temporales"
	@echo "  deps   - Descarga dependencias de Go"
	@echo "  tools  - Instala herramientas (air, swag)"

dev:
	air -c $(AIR_CONF)

swag:
	swag init -g $(MAIN) -d . --parseDependency --parseInternal -o docs

build:
	mkdir -p tmp
	go build -o $(BIN) $(MAIN)

run: build
	$(BIN)

clean:
	@echo "Limpiando..."
	go clean
	@# Elimina la carpeta tmp; si no existe, ignora errores
	@if [ -d tmp ]; then rm -rf tmp; fi

deps:
	go mod download

tools:
	@echo "Instalando air y swag..."
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest