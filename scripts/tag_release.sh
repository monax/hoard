#!/usr/bin/env bash
# uncomment to use the local changes
# vcmd="go run ./cmd/hoarctl/main.go version"

set -e

# Wait a second so we don't see ephemeral file changes
sleep 1

# Don't tag if there is a dirty working dir
if ! git diff-index --quiet HEAD  ; then
    echo "Warning there appears to be uncommitted changes in the working directory:"
    git diff-index HEAD
    echo
    echo "Please commit them or stash them before tagging a release."
fi

echo "This command will tag the current commit $(git rev-parse --short HEAD)"
echo "based on the latest version/release info defined programmatically in"
echo "./release/release.go. It will then push the version tag."
echo "In order for this tag to be released the commit must already be merged"
echo "to master."
echo
read -p "Do you want to continue? [Y\n]: " -r
# Just hitting return defaults to continuing
[[ $REPLY ]] && [[ ! $REPLY =~ ^[Yy]$ ]] && echo && exit 0
echo

# We expect this to be built from HEAD (which is ensured from make tag_release)
vcmd="bin/hoarctl version"
version=v$(${vcmd})
notes=$(${vcmd} notes)
echo "Tagging version $version with message:"
echo ""
echo "$notes"
echo ""
echo "$notes" | git tag -a ${version} -F-

git push origin ${version}

