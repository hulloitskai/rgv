#!/usr/bin/env bash

echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
if [ $? -ne 0 ]; then exit 1; fi

echo "Branches: \$TRAVIS_BRANCH=$TRAVIS_BRANCH, \
  \$RELEASE_BRANCH=$RELEASE_BRANCH"
echo "Image name: $IMAGE_NAME"

if [ "$TRAVIS_BRANCH" == "$RELEASE_BRANCH" ]; then
  ## Push with :latest tag.
  docker push ${IMAGE_NAME}:latest
fi

## Push with branch-specific tag.
docker tag ${IMAGE_NAME}:latest \
           ${IMAGE_NAME}:${TRAVIS_BRANCH}
docker push ${IMAGE_NAME}:${TRAVIS_BRANCH}
