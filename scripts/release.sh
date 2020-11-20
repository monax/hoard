#!/usr/bin/env bash

set -e

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [[ -z ${DOCKER_USER} || -z ${DOCKER_PASS} ]]; then
    echo '$DOCKER_USER and $DOCKER_PASS not set.'
    exit 1
fi

# If passed argument try to use that as tag otherwise read from local repo
if [[ $1 ]]; then
    # Override mode, try to release this tag
    export tag=$1
else
    echo "Getting tag from HEAD which is $(git rev-parse HEAD)"
    export tag=$(git tag --points-at HEAD)
fi

if [[ ! ${tag} ]]; then
    echo "No tag so not releasing."
    exit 0
fi


# Only release semantic version syntax tags
version_regex="^v[0-9]+\.[0-9]+\.[0-9]+$"
if [[ ! ${tag} =~ ${version_regex} ]] ; then
    echo "Tag '$tag' does not match version regex '$version_regex' so not releasing."
    exit 0
fi

echo "Building and pushing $tag..."
$script_dir/build_tool.sh ${tag#v}
echo ${DOCKER_PASS} | docker login --username ${DOCKER_USER} ${DOCKER_HUB} --password-stdin
docker build -t quay.io/monax/hoard:${tag#v} -t quay.io/monax/hoard:latest .
docker push quay.io/monax/hoard:${tag#v}

git config --global user.email "billings@monax.io"
yarn-cli-login
yarn version from-git
yarn publish --access public .

[[ -e notes.md ]] && goreleaser --release-notes notes.md || goreleaser
