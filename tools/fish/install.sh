#!/usr/bin/env bash

## Import Util Fuctions
source ./scripts/util.sh

## Install Macos AMD64
install_macos_amd64() {
  mkdir -p ~/.config
  ln -s ~/.dotfiles/tools/fish/fish ~/.config
  brew install fish fisher
}

## Install Docker Ubuntu
install_docker_arch() {
  skip
}

configure() {
  fisher install ilancosman/tide
  fisher install franciscolourenco/done
  fisher install acomagu/fish-async-prompt
  fisher install jethrokuan/z
}

execute_platform_installer