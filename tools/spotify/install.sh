#!/usr/bin/env bash

## Import Util Fuctions
source ./scripts/util.sh

## Install Macos AMD64
install_macos_amd64() {
  brew install spotify --cask
}

install_macos_arm() {
  brew install spotify --cask
}

## Install Docker Ubuntu
install_docker_arch() {
  skip
}

execute_platform_installer