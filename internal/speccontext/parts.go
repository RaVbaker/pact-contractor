package speccontext

type PartsContext interface {
	Merged() bool
	Name() string
	Defined() bool
}
