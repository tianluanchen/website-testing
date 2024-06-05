#!/usr/bin/env bash
if ! [ -d "node_modules" ]; then
    pnpm install
fi
pnpm run format
pnpm run build
