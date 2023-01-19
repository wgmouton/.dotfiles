#!/usr/bin/env bash

## Import Util Fuctions
source ./scripts/util.sh

## Install Macos AMD64
install_macos_amd64() {
  brew install 1password --cask
  brew install 1password-cli
}

install_macos_arm() {
  brew install 1password --cask
  brew install 1password-cli
}

## Install Docker Ubuntu
install_docker_arch() {
  skip
}

execute_platform_installer