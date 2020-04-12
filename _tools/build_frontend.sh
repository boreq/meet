#!/bin/bash
set -e

# Build frontend
echo "Running yarn build"
cd ../meet-frontend
rm -rf dist
yarn build

# Build backend
cd ../meet/ports/http/frontend
echo "Running https://github.com/rakyll/statik"
statik -f -src=../../../../meet-frontend/dist
