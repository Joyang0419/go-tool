package authenticator

type Interface[TSO any] interface {
	Verify(so TSO) error
}
