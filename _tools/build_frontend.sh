#!/bin/bash
set -e

# Build frontend
echo "Running yarn build"
cd ../hydro-frontend
rm -rf dist
yarn build

# Build backend
cd ../hydro/ports/http/frontend
echo "Running https://github.com/rakyll/statik"
statik -f -src=../../../../hydro-frontend/dist
