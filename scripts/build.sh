#!/bin/sh

echo "Building ethix.."

go build -o bin/ethix cmd/ethix/main.go

echo "Built to bin/ethix"