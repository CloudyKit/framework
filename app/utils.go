package app

func NewContextBundle(initial_contexts ...Context) func() ComponentFunc {
	return func() ComponentFunc {
		return func(a *App) {
			a.BindContext(initial_contexts...)
		}
	}
}
