#!/usr/bin/env bash
# uncomment to use the local changes
# vcmd="go run ./cmd/hoarctl/main.go version"

set -e

# Don't tag if there is a dirty working dir
if ! git diff-index --quiet HEAD  ; then
    echo "There are uncommitted changes in the working directory."
    echo "Please commit them or stash them before tagging a release."
    exit 1
fi

# We expect this to be built from HEAD (which is ensured from make tag_release)
vcmd="bin/hoarctl version"
version=v$(${vcmd})
changes=$(${vcmd} changes)
echo "Tagging version $version with message:"
echo ""
echo "$changes"
echo ""
echo "$changes" | git tag -a ${version} -F-

git push origin ${version}

