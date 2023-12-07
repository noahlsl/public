package errx

type Err interface {
	GetCode(err error) int
	GetErr(code int, lang string) error
}
