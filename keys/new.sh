#!/bin/bash
FILE=$1
KEY="enc.pem"
COUNTER=0
while read LINE; do
    COUNTER=$[COUNTER + 1]
    echo -ne "\\033[KPassword count [$COUNTER] Trying password [$LINE]\\r"
	openssl ec -in $KEY -passin pass:"$LINE" -outform DER | xxd -ps 
done < $FILE