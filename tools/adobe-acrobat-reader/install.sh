#!/usr/bin/env bash

## Import Util Fuctions
source ./scripts/util.sh

## Install Macos AMD64
install_macos_amd64() {
  brew install --cask adobe-acrobat-reader
}

install_macos_arm() {
  brew install --cask adobe-acrobat-reader
}

## Install Docker Ubuntu
install_docker_arch() {
  skip
}

execute_platform_installer