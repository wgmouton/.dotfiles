#!/usr/bin/env bash

skip() {
  echo "Skipping this application for this platform...."
}

execute_platform_installer() {
  case ${OSTYPE} in
    ## MacOS
    'darwin_amd64')
      install_macos_amd64
      configure $(pwd)
    ;;
    'darwin21.0')
      install_macos_arm 
    ;;
    
    ## Linux
    'linux')
      install_macos_amd64
    ;;

    ## Docker
    'docker')
      install_macos_amd64
    ;;

    ## -------
    *)
      echo -n "OS NOT SUPPORTED!!!"
      echo
    ;;
  esac
}
