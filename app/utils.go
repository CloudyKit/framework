package app

// NewContextBundle creates a component that will bind the passed context at bootstrap
func NewContextBundle(initial_contexts ...Context) func() ComponentFunc {
	return func() ComponentFunc {
		return func(a *App) {
			a.BindContext(initial_contexts...)
		}
	}
}
