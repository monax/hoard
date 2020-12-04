#!/bin/bash

# This pushes the plaintext 'foo' through various methods and back out the other end (note that hoarctl insert does
# insert to the same address as the first hoarctl put because encrypt drops the header, but seal adds the header back)

input=foo

# first ref is header ref
secret_key=$(echo "$input" | hoarctl ref | jq -r '.[1].SecretKey')

echo ">>"
echo ">> Testing hoarctl..."
echo ">>"
set -x
output=$(echo "$input" | hoarctl put | hoarctl get | hoarctl putseal | hoarctl unsealget | hoarctl encrypt | hoarctl insert | hoarctl stat | hoarctl cat | hoarctl decrypt -k "$secret_key" | hoarctl ref | hoarctl seal | hoarctl reseal | hoarctl unseal | hoarctl get)
set +x
echo ">>"
[[ "$input" = "$output" ]] && echo ">> hoarctl test succeeded!" || ( echo ">> expected output '$output' to equal input '$input'" ; exit 1 )
echo ">>"
