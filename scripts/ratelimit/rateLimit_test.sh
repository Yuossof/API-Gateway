#!/bin/bash
for i in {1..15}
do
   echo "--- Request $i ---"
   curl -i http://localhost:8080/api/health
   echo ""
done