package types

type Extender[R any] interface {
	Extend() R
}
