#!/usr/bin/env bash

## Import Util Fuctions
source ./scripts/util.sh

## Install Macos AMD64
install_macos_amd64() {
  brew install romkatv/powerlevel10k/powerlevel10k
  echo "source $(brew --prefix)/opt/powerlevel10k/powerlevel10k.zsh-theme" >>~/.zshrc
}

install_macos_arm() {
  brew install romkatv/powerlevel10k/powerlevel10k
  echo "source $(brew --prefix)/opt/powerlevel10k/powerlevel10k.zsh-theme" >>~/.zshrc
}

## Install Docker Ubuntu
install_docker_arch() {
  skip
}

execute_platform_installer