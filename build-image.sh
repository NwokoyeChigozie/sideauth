#!/bin/sh
set -e

image_name="vesicash/auth"
image_tag="${1:-latest}"

if [[ ! -f "app.env" ]]
then
    echo "Copying environment file for app ⏳"
    cp app-sample.env app.env
else
    echo "Environment file found 👌"
fi

echo "Building docker image ${image_tag} version 🛠️"
docker build -t $image_name:$image_tag .