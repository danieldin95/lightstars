#!/bin/bash


DIR="/etc/lightstar"
FILE=$(find ${DIR} -name "*.json")


for i in ${FILE}; do
    echo "-- Checking ${i}"
    output=$(python -m json.tool $i 2>&1)
    if [ $? -ne 0 ]; then
        echo "-- ... incorrect and output errors:"
        echo "${output}"
    else
        echo "-- ... success"
    fi
    echo ""
done
