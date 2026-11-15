#!/bin/bash

set -e

ProdVersion="../../../infra/tofu/components/dynamodb.tf"
LocalVersion="./tofu/dynamo.tf"

hash1=$(sha256sum $ProdVersion | awk '{print $1}')
hash2=$(sha256sum $LocalVersion | awk '{print $1}')

if [ "$hash1" != "$hash2" ]; then
	rm $LocalVersion
	cp $ProdVersion $LocalVersion
fi

cd ./tofu
tofu init
tofu apply

