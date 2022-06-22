#!/usr/bin/env bash

## Import Util Fuctions
source ./scripts/util.sh

## Install Macos AMD64
install_macos_amd64() {
  ln -s ~/.dotfiles/tools/git/git ~/.config
}

install_macos_arm() {
  ln -s ~/.dotfiles/tools/git/git ~/.config
  git config --global core.excludesFile '~/.config/git/.gitignore'
}

## Install Docker Ubuntu
install_docker_arch() {
  skip
}

execute_platform_installer