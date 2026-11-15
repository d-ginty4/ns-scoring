#!/usr/bin/env bash

set -e

upload(){
    BUCKET="ns-scoring-backend-dev-590184097018-storage-eu-west-1"
    DIRECTORY=$1
    ZIP=$DIRECTORY.zip
    FUNCTION=ns-scoring-backend-dev-$1-controller

    cd "$DIRECTORY"

    # the script will compile and zip the lambda folder and upload it to s3
    GOOS=linux CGO_ENABLED=0 GOARCH=arm64 go build -o bootstrap .

    zip "$ZIP" bootstrap

    # Check if s3 bucket exists before uploading
    if ! aws s3api head-bucket --bucket $BUCKET > /dev/null;then
        echo "Bucket $BUCKET does not exist or you don't have permission to access it"
        exit 1
    fi

    # Upload zip to s3 bucket
    aws s3 cp "$ZIP" s3://$BUCKET

    if [ "$CREATE" = false ]; then
        aws lambda update-function-code \
            --function-name "$FUNCTION" \
            --s3-bucket "$BUCKET" \
            --s3-key "$ZIP" > /dev/null
        echo "Updated lambda $FUNCTION"
    fi

    rm bootstrap "$ZIP"

    exit 0
}

LAMBDAS=($(ls -d */ | sed 's#/##'))

if [ -z "$1" ]; then
    echo "No lambda provided"
    echo "Please use one of the following: ${LAMBDAS[@]}"
    exit 1
fi

CREATE=false
for arg in "$@"; do
    if [[ "$arg" == "-c" || "$arg" == "--create" ]]; then
        CREATE=true
        break
    fi
done

for folder in "${LAMBDAS[@]}"; do
    if [[ "$folder" == "$1" ]]; then
        upload "$1"
    fi
done

echo "$1 lambda not found"
echo "Please use one of the following: ${LAMBDAS[@]}"
