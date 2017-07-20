#!/usr/bin/env bash

version_regex="^v[0-9]+\.[0-9]+\.[0-9]+$"

function release {
    echo "Building and releasing $tag..."
    goreleaser
}

# If passed argument try to use that as tag otherwise read from local repo
if [[ $1 ]]; then
    export tag=$1
else
    export tag=$(git tag --points-at HEAD)
fi

# Only release semantic version syntax tags
if [[ ${tag} =~ ${version_regex} ]]; then
    release
elif [[ ${tag} ]]; then
    echo "Tag '$tag' does not match version regex '$version_regex' so not releasing."
else
    echo "No tag so not releasing."
fi
