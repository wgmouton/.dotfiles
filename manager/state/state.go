package state

type Tool struct {
	Version string
}

type state struct {
	Installed map[string]Tool
}

func (s *state) SetState() {

}
func (s *state) GetState() *state {
	return s
}
