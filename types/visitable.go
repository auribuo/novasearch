package types

type Visitable interface {
	Visit()
	Reset()
}
