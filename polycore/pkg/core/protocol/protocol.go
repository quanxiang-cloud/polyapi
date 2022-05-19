package protocol

// Evaler is script evaler object
type Evaler interface {
	Eval(string) error
	Free()
}
