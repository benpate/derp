package derp

type Unwrapper interface {
	Error() string
	Unwrap() error
}
