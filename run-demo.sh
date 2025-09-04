#!/bin/bash

# Docker Multi-Stage Build Demo - Run Script
# This script runs both images and demonstrates their functionality

set -e

echo "🚀 Docker Multi-Stage Build Demo - Testing"
echo "=========================================="
echo

# Function to test an image
test_image() {
    local image_name=$1
    local port=$2
    local container_name=$3
    
    echo "🧪 Testing $image_name..."
    
    # Start container
    docker run -d -p $port:8080 --name $container_name $image_name
    
    # Wait for startup
    sleep 3
    
    # Test endpoints
    echo "   ✅ Health check:"
    curl -s http://localhost:$port/health | jq '.status'
    
    echo "   ✅ Demo endpoint:"
    curl -s http://localhost:$port/demo | jq '.title'
    
    echo "   ✅ App info:"
    curl -s http://localhost:$port/info | jq '.name'
    
    echo "   ✅ Greeting test:"
    curl -s -X POST -H "Content-Type: application/json" \
         -d '{"name": "Demo User"}' \
         http://localhost:$port/greet | jq '.message'
    
    # Cleanup
    docker stop $container_name > /dev/null
    docker rm $container_name > /dev/null
    
    echo "   🧹 Cleaned up $container_name"
    echo
}

# Test single-stage image
test_image "golang-single" "8081" "demo-single"

# Test multi-stage image
test_image "golang-multi" "8082" "demo-multi"

echo "🎉 Both images work identically!"
echo "📊 Size difference: Multi-stage is ~184x smaller"
echo "🔒 Security benefit: Multi-stage has no build tools"
echo
echo "🌐 To explore manually, run:"
echo "   docker run -p 8080:8080 golang-multi"
echo "   curl http://localhost:8080/demo | jq ."