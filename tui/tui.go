package tui

type Widget interface {
	Build() error
}

// Run start the tui app and run the loop
func Run() error {
	return nil
}
