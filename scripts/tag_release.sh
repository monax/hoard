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

echo "This command will create a tag based on the latest release defined"
echo "programmatically in ./release/release.go. It will then push that version tag."
echo "In order to release merge the tagged commit to master."
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

