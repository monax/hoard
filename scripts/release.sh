#!/usr/bin/env bash

version_regex="^v[0-9]+\.[0-9]+\.[0-9]+$"

function release {
    echo "Building and releasing $tag..."
    [[ -e notes.md ]] && goreleaser --release-notes notes.md || goreleaser
}


# If passed argument try to use that as tag otherwise read from local repo
if [[ $1 ]]; then
    # Override mode, try to release this tag
    export tag=$1
else
    echo "Getting tag from HEAD which is $(git rev-parse HEAD)"
    export tag=$(git tag --points-at HEAD)
    # Only release from master unless being run as override
    if [[ $(git symbolic-ref HEAD) != "refs/heads/master" ]]; then
        echo "Branch is not master so not releasing."
        exit 0
    fi
fi

if [[ ! ${tag} ]]; then
    echo "No tag so not releasing."
    exit 0
fi

# Only release semantic version syntax tags
if [[ ! ${tag} =~ ${version_regex} ]] ; then
    echo "Tag '$tag' does not match version regex '$version_regex' so not releasing."
    exit 0
fi

release
