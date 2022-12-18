package util

type Symlink interface {
	Link(path, link string) error
}

type Windows struct{}
type Macos struct{}
type Arch struct{}
