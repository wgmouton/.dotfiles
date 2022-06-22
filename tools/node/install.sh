#!/usr/bin/env bash

## Import Util Fuctions
source ./scripts/util.sh

## Install Macos AMD64
install_macos_amd64() {
  brew install node npm
}

install_macos_arm() {
  brew install node npm
}

## Install Docker Ubuntu
install_docker_arch() {
  skip
}

execute_platform_installer