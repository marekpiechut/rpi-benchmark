#!/usr/bin/env bash

url=$1
connections=(10 20 50 250 500 1200 2000 6000 10000)
token="eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ0ZXN0dXNlciJ9.gfWlnHqLf-Uv1pTYBIvN0nlYwl1tXrxrGBnSn-dwhuc"
upload="./bench-img.png"
runtime="10s"

for c in "${connections[@]}"
do
	echo "Connections: ${c}"
	out="bombardier-${c}.out"
	bombardier -p i,r -t 45s --fasthttp -c ${c} -d ${runtime} -m POST -H "Authorization:Bearer ${token}" -f ${upload} ${url} | tee "${out}"
done