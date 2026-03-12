# ============================================================
# Config
# ============================================================
APP_NAME    := event-budaya-ticketing-bcc
BINARY      := bin/app
MAIN        := ./cmd/main.go
BRANCH      := $(shell git rev-parse --abbrev-ref HEAD)
SERVICE     := event-budaya   # sesuaikan dengan nama unit systemd kamu

# ============================================================
# Help
# ============================================================
.PHONY: help
help:
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "  redeploy   - Full redeploy: pull → deps → build → restart"
	@echo "  pull       - Git pull branch saat ini"
	@echo "  deps       - Download & tidy Go dependencies"
	@echo "  build      - Build binary ke $(BINARY)"
	@echo "  restart    - Restart systemd service $(SERVICE)"
	@echo "  status     - Cek status systemd service"
	@echo "  logs       - Lihat log service secara live"
	@echo ""

# ============================================================
# Main target
# ============================================================
.PHONY: redeploy
redeploy: pull deps build restart
	@echo ""
	@echo "Redeploy selesai. Service $(SERVICE) sudah berjalan."

# ============================================================
# Steps
# ============================================================
.PHONY: pull
pull:
	@echo ">>> [1/4] Git pull branch: $(BRANCH)"
	git pull origin $(BRANCH)

.PHONY: deps
deps:
	@echo ">>> [2/4] Install & tidy dependencies"
	go mod download
	go mod tidy

.PHONY: build
build:
	@echo ">>> [3/4] Build binary -> $(BINARY)"
	@mkdir -p bin
	go build -o $(BINARY) $(MAIN)

.PHONY: restart
restart:
	@echo ">>> [4/4] Restart service: $(SERVICE)"
	systemctl restart $(SERVICE)

# ============================================================
# Utilities
# ============================================================
.PHONY: status
status:
	systemctl status $(SERVICE)

.PHONY: logs
logs:
	journalctl -u $(SERVICE) -f
