#!/usr/bin/env bash

version_regex="^v[0-9]+\.[0-9]+\.[0-9]+$"
export deploy_dir="./bin"

function release {
    echo "Releasing $tag..."
    if [[ $(ls ${deploy_dir} 2> /dev/null) ]]; then
        echo "Deploying artifacts from '$deploy_dir'..."
        ghr ${tag} ${deploy_dir}
    else
        echo "No artifacts in '$deploy_dir' aborting deployment."
    fi
}

# If passed argument then use as tag otherwise read from local repo
if [[ $1 ]]; then
    export tag=$1
else
    export tag=$(git tag --points-at HEAD)
fi

if [[ ${tag} =~ ${version_regex} ]]; then
    release
elif [[ ${tag} ]]; then
    echo "Version tag '$tag' did not match version regex '$version_regex'"
else
    echo "No tag so not releasing."
fi
