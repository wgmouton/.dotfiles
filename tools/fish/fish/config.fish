if status is-interactive
    # Commands to run in interactive sessions can go here
end

test -e {$HOME}/.iterm2_shell_integration.fish ; and source {$HOME}/.iterm2_shell_integration.fish
eval "$(/opt/homebrew/bin/brew shellenv)"

set -x GOPATH (go env GOPATH)
set -x PATH $PATH (go env GOPATH)/bin