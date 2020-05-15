package tui

type Widget interface {
	Build() error
}
