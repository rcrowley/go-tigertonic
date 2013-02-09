package tigertonic

type MarshalerError struct {
	s string
}

func newMarshalerError(s string) error {
	return MarshalerError{s}
}

func (e MarshalerError) Error() string { return e.s }
