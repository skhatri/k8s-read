#!/usr/bin/env bash

set -eo pipefail

current_hash=$(git log --pretty=format:'%h' --max-count=1)
current_branch=$(git branch --show-current|sed 's#/#_#')


create_tag() {
    if [[ ${current_branch} == "main" ]]; 
    then
      source project.properties
      git config --global user.email "${email}"
      git config --global user.name "${name}"
      bumpversion build
    fi;
}
create_tag
version=$(cat VERSION)

if [[ ! -z ${version} ]];
then
  source project.properties
  image_version_tag="${owner}/${project}:${version}"
  image_latest_tag="${owner}/${project}:latest"
  echo building ${image_version_tag}
  docker build --no-cache -t ${image_version_tag} .
  docker push ${image_version_tag}
  docker tag ${image_version_tag} ${image_latest_tag}
  docker push ${image_latest_tag}
  git push origin main --tags
fi;