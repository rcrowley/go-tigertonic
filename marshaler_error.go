package tigertonic

type MarshalerError string

func (e MarshalerError) Error() string { return string(e) }
