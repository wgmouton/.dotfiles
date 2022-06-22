#!/usr/bin/env bash

## Import Util Fuctions
source ./scripts/util.sh

## Install Macos AMD64
install_macos_amd64() {
  skip
}

install_macos_arm() {
  /usr/sbin/softwareupdate --install-rosetta --agree-to-license
  defaults write com.apple.finder AppleShowAllFiles -boolean true; killall Finder;
}

## Install Docker Ubuntu
install_docker_arch() {
  skip
}

execute_platform_installer