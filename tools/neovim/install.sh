#!/usr/bin/env bash

## Import Util Fuctions
source ./scripts/util.sh

## Install Macos AMD64
install_macos_amd64() {
  brew install neovim
}

## Install Docker Ubuntu
install_docker_arch() {
  skip
}

## ARCH sudo pacman -S neovim

execute_platform_installer