#!/usr/bin/env bash

## Import Util Fuctions
source ./scripts/util.sh

## Install Macos AMD64
install_macos_amd64() {
  brew install iterm2 --cask
}

install_macos_arm() {
  brew install iterm2 --cask
}

## Install Docker Ubuntu
install_docker_arch() {
  skip
}

execute_platform_installer