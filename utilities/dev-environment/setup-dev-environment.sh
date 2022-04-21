#!/usr/bin/env bash

CONTEXT=$1
DIR=$2

tmux
vim

if [$CONTEXT == "docker"] then
  docker run .....
fi