#!/bin/bash -e

OWNER="askasoft"
REPO="pangox-xdemo"
NAME="xdemo"

TAG=$1
ARCH=$2
if [ -z $ARCH ]; then
  ARCH=arm64
fi
ASSET_NAME=${NAME}-linux-$ARCH


: "${GITHUB_TOKEN:?GITHUB_TOKEN is not set}"

API="https://api.github.com/repos/${OWNER}/${REPO}"

if [ -z $TAG]; then
  TAG=snapshot
fi

if [ "$TAG" = "latest" ]; then
  TAG=$(
    curl -fsSL \
        -H "Authorization: Bearer ${GITHUB_TOKEN}" \
        "${API}/releases/latest" |
    jq -r '.tag_name'
  )
fi

ASSET_ID=$(
    curl -fsSL \
        -H "Authorization: Bearer ${GITHUB_TOKEN}" \
        "${API}/releases/tags/${TAG}" |
    jq -r \
        --arg NAME "${ASSET_NAME}" \
        '.assets[] | select(.name==$NAME) | .id'
)

test -n "${ASSET_ID}"

echo "  > Downloading ${TAG}/${ASSET_NAME} ... "

curl -fL \
    -H "Authorization: Bearer ${GITHUB_TOKEN}" \
    -H "Accept: application/octet-stream" \
    "${API}/releases/assets/${ASSET_ID}" \
    -o ${NAME}

chmod a+x ${NAME}
