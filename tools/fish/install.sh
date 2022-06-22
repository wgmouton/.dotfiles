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
  
  brew tap homebrew/cask-fonts
  brew install --cask font-hack-nerd-font

  brew install fish fisher

  fish -c "fisher update ilancosman/tide"
  fish -c "fisher update franciscolourenco/done"
  # fish -c "fisher intall acomagu/fish-async-prompt" ##NOTE: Breaks on m1 mac
  fish -c "fisher update jethrokuan/z"
  fish -c "fisher update PatrickF1/fzf.fish"
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