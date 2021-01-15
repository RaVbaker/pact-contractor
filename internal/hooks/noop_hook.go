package hooks

type NoopHook struct{}

func (n NoopHook) Run(_ string) error {
	return nil
}
