#!/usr/bin/env bash

## Import Util Fuctions
source ./scripts/util.sh

## Install Macos AMD64
install_macos_amd64() {
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  brew tap homebrew/autoupdate
  brew autoupdate start
}

install_macos_arm() {
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  brew tap homebrew/autoupdate
  brew autoupdate start
}

## Install Docker Ubuntu
install_docker_arch() {
  skip
}

configure() {
  
}

execute_platform