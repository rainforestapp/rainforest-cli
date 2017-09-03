#!/usr/bin/env bash

# Parse the CHANGELOG.md to extract the first empty-line-deliminated block and
# update the github releases page with it
# e.g. given a CHANGELOG.md of:
# <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
# # TITLE OF CHANGELOG
#
# ## Latest changelog entry
# - list of
# - what got
# - updated
#
# ## Next changelog entry
# ...
# >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
# Insert into the github release:
# - list of
# - what got
# - updated
#
# The GITHUB_AUTH_TOKEN environment variable must be set and must contain a
# Github personal access token with at least the repo:public_repo scope

set -eou pipefail

function create_release() {
  TOKEN=$1
  API_PATH=$2
  TAG=$3
  RELEASE_NOTES=$4

  curl \
    -H "Authorization: token $TOKEN" \
    "https://api.github.com/$API_PATH/releases" \
    --data-binary @- <<-EOF
  {
    "tag_name": "$TAG",
    "target_commitish": "master",
    "name": "$TAG",
    "body": "$RELEASE_NOTES",
    "draft": false,
    "prerelease": false
  }
EOF

}

TAG=${1:=CIRCLE_TAG}
ORG=${2:=CIRCLE_PROJECT_USERNAME}
REPO=${3:=CIRCLE_PROJECT_REPONAME}
CHANGELOG=${4:=CHANGELOG.md}

RELEASE_NOTES=$(awk -v RS='' '/##/ { print; exit }' "$CHANGELOG" | sed '1d; s/$/\\n/' | tr -d '\n')

create_release "$GITHUB_AUTH_TOKEN" "repos/$ORG/$REPO" "$TAG" "$RELEASE_NOTES"

# How to upload the binary
# ASSET=$4
# RELEASE_JSON=$(create_release "$GITHUB_AUTH_TOKEN" "$API_PATH" "$TAG" "$RELEASE_NOTES")
# UPLOAD_URL=$(jq -r .upload_url <<< "$RELEASE_JSON" | sed 's/{.*$//')
# curl -H "Authorization: token $GITHUB_AUTH_TOKEN" \
#   -H "Content-Type: application/zip" \
#   --data-binary "@$ASSET" \
#   "$UPLOAD_URL?name=$REPO-$TAG.zip"
