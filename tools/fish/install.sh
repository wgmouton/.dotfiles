#!/usr/bin/env bash

## Import Util Fuctions
source ./scripts/util.sh

## Install Macos AMD64
install_macos_amd64() {
  mkdir -p ~/.config
  ln -s ~/.dotfiles/tools/fish/fish ~/.config
  brew install fish fisher
}

install_macos_arm() {
  mkdir -p ~/.config
  ln -s ~/.dotfiles/tools/fish/fish ~/.config
  brew install fish fisher

  fish -c "fisher install ilancosman/tide"
  fish -c "fisher install franciscolourenco/done"
  fish -c "fisher install acomagu/fish-async-prompt"
  fish -c "fisher install jethrokuan/z"
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