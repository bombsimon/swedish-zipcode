#!/bin/sh

START=${START:-10000}
STOP=${STOP:-99999}
FILENAME=${FILENAME:-sweden-zipcode.csv}
COUNTER=0
STEP=${STEP:-10}

echo "Removing $FILENAME"
rm "$FILENAME"

echo "Fetching zip codes between $START and $STOP. Will print every $STEP step"
while [ "$START" -lt "$STOP" ]; do
    COUNTER=$(( COUNTER + 1 ))
    if [ "$COUNTER" -ge "$STEP" ]; then
        echo "Fetching $START"
        COUNTER=0
    fi

    id="$START"
    read -r valid result \
        <<<"$(curl -sL "https://api.bring.com/shippingguide/api/postalCode.json?clientUrl=ex&country=SE&pnr=${id}" | jq -r '"\(.valid) \(.result)"')"

    if [ "$valid" = "true" ]; then
        echo "${id},${result}" >> "$FILENAME"
    fi

    START=$(( START + 1 ))
done
