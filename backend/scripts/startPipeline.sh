#!/bin/bash

set -eu

BRANCH=$(git rev-parse --abbrev-ref HEAD)
COMMIT_HASH=$(git rev-parse HEAD)

if ! git ls-remote --exit-code --heads origin "$BRANCH" >/dev/null; then
  echo "Branch $BRANCH does not exist in origin" >&2
  exit 1
fi

git fetch origin "$BRANCH"

if ! git merge-base --is-ancestor "$COMMIT_HASH" "origin/$BRANCH"; then
  echo "Commit $COMMIT_HASH is not on origin/$BRANCH (push first?)" >&2
  exit 1
fi

aws codepipeline start-pipeline-execution \
  --name "ns-scoring-backend-dev-deploy" \
  --source-revisions actionName=Source,revisionType=COMMIT_ID,revisionValue=$COMMIT_HASH
