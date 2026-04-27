#!/bin/bash

set -u

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

BINARY_NAME="lexpilot"
SERVICE_NAME="lexpilot"
BUILD_TARGET="./cmd/lexpilot"

echo -e "${BLUE}Spoustim deployment LexPilot Demo...${NC}"

export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin:/home/skybedy/go/bin
if ! command -v go >/dev/null 2>&1; then
    echo -e "${RED}Go neni dostupne v PATH.${NC}"
    exit 1
fi

if [ ! -f ".env" ]; then
    echo -e "${YELLOW}Chybi .env. Aplikace po restartu pobezi v mock demo rezimu, pokud OPENAI_API_KEY neni nastaven jinak.${NC}"
fi

if ! command -v pdftotext >/dev/null 2>&1; then
    echo -e "${YELLOW}pdftotext neni dostupny. PDF upload nebude fungovat, dokud nenainstalujes poppler-utils.${NC}"
    echo -e "${YELLOW}   Ubuntu/Debian: sudo apt install poppler-utils${NC}"
fi

echo -e "${BLUE}Stahuji zmeny z Gitu...${NC}"

CURRENT_BRANCH=$(git branch --show-current)
if [ -z "$CURRENT_BRANCH" ]; then
    echo -e "${RED}Nepodarilo se zjistit aktualni git vetev.${NC}"
    exit 1
fi

git pull origin "$CURRENT_BRANCH"
if [ $? -ne 0 ]; then
    echo -e "${RED}Chyba pri stahovani z Gitu.${NC}"
    exit 1
fi

echo -e "${BLUE}Spoustim testy...${NC}"
go test ./...
if [ $? -ne 0 ]; then
    echo -e "${RED}Testy selhaly. Deployment zastaven.${NC}"
    exit 1
fi

echo -e "${BLUE}Sestavuji binarku...${NC}"
go build -o "$BINARY_NAME" "$BUILD_TARGET"
if [ $? -ne 0 ]; then
    echo -e "${RED}Build selhal.${NC}"
    exit 1
fi

echo -e "${GREEN}Build byl uspesny.${NC}"

if command -v systemctl >/dev/null 2>&1 && systemctl list-unit-files "${SERVICE_NAME}.service" >/dev/null 2>&1; then
    echo -e "${BLUE}Restartuji systemd sluzbu ${SERVICE_NAME}...${NC}"
    sudo systemctl restart "${SERVICE_NAME}"
    if [ $? -ne 0 ]; then
        echo -e "${RED}Restart sluzby ${SERVICE_NAME} selhal.${NC}"
        exit 1
    fi
    sudo systemctl status "${SERVICE_NAME}" --no-pager -l
else
    echo -e "${YELLOW}Systemd sluzba ${SERVICE_NAME}.service nebyla nalezena.${NC}"
    echo -e "${YELLOW}   Spust rucne: ./${BINARY_NAME}${NC}"
fi

echo -e "${GREEN}Deployment LexPilot Demo dokoncen.${NC}"
