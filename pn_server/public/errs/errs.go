package errs

type Option func(err *Err)

func WithMsg(msg string) Option {
	return func(err *Err) {
		err.msg = msg
	}
}

type IErr interface {
	Error() string
	Code() int
}

var _ IErr = &Err{}
var _ error = &Err{}

type Err struct {
	code int
	msg  string
}

func (e *Err) Error() string {
	return e.msg
}

func (e *Err) Code() int {
	return e.code
}

func New(code int, opts ...Option) *Err {
	err := &Err{
		code: code,
		msg:  "",
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}
