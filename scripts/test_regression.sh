#!/usr/bin/env bash

[[ -z "$REGRESSION_DIR" ]] && echo "Please provide REGRESSION_DIR path" && exit 1
[[ -z "$REGRESSION_SNAPSHOT" ]] && echo "Please provide REGRESSION_SNAPSHOT name" && exit 1

# We only check stability of the plaintexts
REGRESSION_FIXTURES="$REGRESSION_DIR/fixtures"
REGRESSION_OUTPUT="$REGRESSION_DIR/snapshots/$REGRESSION_SNAPSHOT"
REGRESSION_PLAINTEXTS="$REGRESSION_OUTPUT/plaintexts"

go run test/regression/main.go "$REGRESSION_FIXTURES" "$REGRESSION_OUTPUT"

changes=$(git ls-files --deleted --modified --other -- "$REGRESSION_PLAINTEXTS")

if [[ -z "$changes" ]]; then
  echo
  echo "Regression test passed!"
  echo
  # Clear out any changes we don't care about
  git clean --quiet -fdx -- "$REGRESSION_OUTPUT" && git checkout --quiet --no-overlay HEAD -- "$REGRESSION_OUTPUT"
else
  echo
  echo "Regression test failed, plaintext changes detected:"
  echo
  echo "$changes"
  echo
  echo "Commit these files if you would like to create or update the snapshot '$REGRESSION_SNAPSHOT'"
fi

