package contracts

type (
	Registry interface {
	}

	Application interface {
		Bootstrap(bundles ...Bundle) error
	}

	Bundle interface {
		Bootstrap(registry Registry) error
	}
)
