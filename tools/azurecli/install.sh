#!/usr/bin/env bash

## Import Util Fuctions
source ./scripts/util.sh

## Install Macos AMD64
install_macos_amd64() {
  brew install azure-cli
}

install_macos_arm() {
  brew install azure-cli
}

## Install Docker Ubuntu
install_docker_arch() {
  skip
}

configure() {
  az extension add --name azure-devops
}

execute_platform_installer