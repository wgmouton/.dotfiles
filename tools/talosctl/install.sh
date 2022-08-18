#!/usr/bin/env bash

## Import Util Fuctions
source ./scripts/util.sh

## Install Macos AMD64
install_macos_amd64() {
  sudo curl -Lo /usr/local/bin/talosctl https://github.com/siderolabs/talos/releases/download/v1.1.1/talosctl-$(uname -s | tr "[:upper:]" "[:lower:]")-amd64
  sudo chmod +x /usr/local/bin/talosctl
}

install_macos_arm() {
  sudo curl -Lo /usr/local/bin/talosctl https://github.com/siderolabs/talos/releases/download/v1.1.1/talosctl-$(uname -s | tr "[:upper:]" "[:lower:]")-arm64
  sudo chmod +x /usr/local/bin/talosctl
}

## Install Docker Ubuntu
install_docker_arch() {
  skip
}

execute_platform_installer