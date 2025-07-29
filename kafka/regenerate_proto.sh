#!/bin/bash

set -e

echo "Regenerating Protocol Buffers..."

cd "$(dirname "$0")/kafka-model/resources/proto"

# Generate Go code
protoc --go_out=../generated \
       --go_opt=paths=source_relative \
       *.proto

echo "Protocol Buffers regenerated successfully!"
echo "Generated files:"
ls -la ../generated/ 