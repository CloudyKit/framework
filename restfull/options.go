package restfull

import "github.com/CloudyKit/framework/request"

type Option[Type Resource] func(resource *Controller[Type])

func WithPerPageLimit[Type Resource](perPage int) Option[Type] {
	return func(resource *Controller[Type]) {
		resource.perPage = perPage
	}
}

func WithFindOneFilters[Type Resource](filters ...request.Handler) Option[Type] {
	return func(resource *Controller[Type]) {
		resource.findOneFilters = append(resource.findOneFilters, filters...)
	}
}

func WithFindAllFilters[Type Resource](filters ...request.Handler) Option[Type] {
	return func(resource *Controller[Type]) {
		resource.findAllFilters = append(resource.findAllFilters, filters...)
	}
}

func WithCreateOneFilters[Type Resource](filters ...request.Handler) Option[Type] {
	return func(resource *Controller[Type]) {
		resource.createOneFilters = append(resource.createOneFilters, filters...)
	}
}

func WithUpdateOneFilters[Type Resource](filters ...request.Handler) Option[Type] {
	return func(resource *Controller[Type]) {
		resource.updateOneFilters = append(resource.updateOneFilters, filters...)
	}
}

func WithDeleteOneFilters[Type Resource](filters ...request.Handler) Option[Type] {
	return func(resource *Controller[Type]) {
		resource.deleteOneFilters = append(resource.deleteOneFilters, filters...)
	}
}

func WithReplaceOneFilters[Type Resource](filters ...request.Handler) Option[Type] {
	return func(resource *Controller[Type]) {
		resource.replaceOneFilters = append(resource.replaceOneFilters, filters...)
	}
}
