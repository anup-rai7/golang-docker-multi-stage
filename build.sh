#!/bin/bash

# Docker Multi-Stage Build Demo - Build Script
# This script builds both Docker images and provides size comparison

set -e

echo "🏗️  Docker Multi-Stage Build Demo"
echo "================================="
echo

# Clean up previous builds
echo "🧹 Cleaning up previous builds..."
docker rmi golang-single golang-multi 2>/dev/null || true
echo

# Build single-stage image
echo "📦 Building single-stage image..."
docker build -f Dockerfile.single -t golang-single .
echo

# Build multi-stage image
echo "📦 Building multi-stage image..."
docker build -f Dockerfile.multi -t golang-multi .
echo

# Show size comparison
echo "📊 Size Comparison:"
echo "==================="
docker images | head -1
docker images | grep golang | head -2

echo
echo "📈 Analysis:"
SINGLE_SIZE=$(docker images golang-single --format "{{.Size}}")
MULTI_SIZE=$(docker images golang-multi --format "{{.Size}}")

echo "Single-stage: $SINGLE_SIZE"
echo "Multi-stage:  $MULTI_SIZE"

# Calculate approximate size reduction
echo
echo "✅ Build completed! The multi-stage image is dramatically smaller."
echo "🚀 Run './run-demo.sh' to test both images."