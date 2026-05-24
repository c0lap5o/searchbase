#!/bin/bash

echo "Testing the Search API..."
curl -X POST http://localhost:8080/api/v1/search \
     -H "Content-Type: application/json" \
     -d '{
           "query": "what is kubernetes?",
           "limit": 3,
           "js_render": false
         }' | jq
