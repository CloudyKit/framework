package restfull

import (
	"encoding/json"
	"fmt"
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/framework/validation"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type SortField struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

type SearchContext struct {
	Page    int
	PerPage int
	SortBy  []SortField `json:"sortBy"`
}

type findAllResponse struct {
	Records any `json:"records"`
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
	Total   int `json:"total"`
}

type Resource interface {
	FindAllModel() any
	FindOneModel() any
	FindUpdateModel() any

	Model() any
	UpdateOneModel() any

	FindAll(pageData SearchContext) (any, int, validation.Result, error)
	FindOne(resourceId string) (any, validation.Result, error)

	CreateOne() (validation.Result, error)
	UpdateOne(resourceId string) (validation.Result, error)
	DeleteOne(resourceId string) (validation.Result, error)
	ReplaceOne(resourceId string) (validation.Result, error)
}

type Controller[Type Resource] struct {
	name       string
	Controller Type
	*request.Context
	filters           []request.Handler
	findOneFilters    []request.Handler
	findAllFilters    []request.Handler
	createOneFilters  []request.Handler
	updateOneFilters  []request.Handler
	deleteOneFilters  []request.Handler
	replaceOneFilters []request.Handler
	perPage           int
}

func NewController[Type Resource](name string, opts ...Option[Type]) app.Controller {
	var t Type
	controller := &Controller[Type]{name: name, Controller: t}
	for _, opt := range opts {
		opt(controller)
	}
	return controller
}

func (resource *Controller[Type]) BasePath() string {
	name := strings.TrimRight(resource.name, "/")
	return name
}

func (resource *Controller[Type]) Mx(mapper *app.Mapper) {

	mapper.BindFilterHandlers(resource.filters...)

	var resourceAPI = reflect.ValueOf(resource.Controller)
	var resourceType = resourceAPI.Type()

	if resourceType.Kind() != reflect.Ptr && resourceType.Kind() != reflect.Interface {
		panic(fmt.Errorf("resource value should be a pointer to struct"))
	}
	resourceType = resourceType.Elem()

	mapper.BindFilterFuncHandlers(func(context *request.Context) {
		context.Response.Header().Set("Content-Type", "application/json")
		newVAL := reflect.New(resourceType)
		newVALStruct := newVAL.Elem()
		if !resourceAPI.IsZero() {
			newVALStruct.Set(resourceAPI.Elem())
		}
		context.Registry.InjectValue(newVALStruct)
		context.Registry.WithTypeAndValue(reflect.TypeOf(&resource.Controller).Elem(), newVAL.Interface())
		_ = context.Next()
	})

	mapper.BindAction("GET", fmt.Sprintf("/%s", resource.name), "FindAll", resource.findAllFilters...)
	mapper.BindAction("GET", fmt.Sprintf("/%s/search/:page", resource.name), "FindAll", resource.findAllFilters...)
	mapper.BindAction("POST", fmt.Sprintf("/%s/search/:page", resource.name), "FindAll", resource.findAllFilters...)
	mapper.BindAction("GET", fmt.Sprintf("/%s/:resourceId", resource.name), "FindOne", resource.findOneFilters...)
	mapper.BindAction("POST", fmt.Sprintf("/%s", resource.name), "CreateOne", resource.createOneFilters...)
	mapper.BindAction("PATCH", fmt.Sprintf("/%s/:resourceId", resource.name), "UpdateOne", resource.updateOneFilters...)
	mapper.BindAction("DELETE", fmt.Sprintf("/%s/:resourceId", resource.name), "DeleteOne", resource.deleteOneFilters...)
	mapper.BindAction("PUT", fmt.Sprintf("/%s/:resourceId", resource.name), "ReplaceOne", resource.replaceOneFilters...)
}

type apiError struct {
	Status        int    `json:"status,omitempty"`
	StatusMessage string `json:"status_message,omitempty"`
	Errors        any    `json:"errors,omitempty"`
}

func (resource *Controller[Type]) sendError(statusCode int, errorMessage string, errors validation.Result) {
	err := &apiError{
		Status:        statusCode,
		StatusMessage: errorMessage,
		Errors:        errors,
	}
	resource.Response.WriteHeader(statusCode)
	var jsonEncoder = json.NewEncoder(resource.Response)
	_ = jsonEncoder.Encode(err)
}

func (resource *Controller[Type]) CreateOne() {
	createModel := resource.Controller.Model()
	if createModel != nil {
		bindErr := resource.BindJSON(createModel)
		if bindErr != nil {
			resource.sendError(http.StatusBadRequest, "error parsing the parameters", nil)
			return
		}
	}

	result, err := resource.Controller.CreateOne()
	if err != nil {
		resource.sendError(http.StatusInternalServerError, "error processing the request", nil)
		return
	}
	if result.HasErrors() {
		resource.sendError(http.StatusBadRequest, "error parsing the parameters", result)
	}

	err = json.NewEncoder(resource.Response).Encode(createModel)
	if err != nil {
		resource.sendError(http.StatusInternalServerError, "error processing the request", nil)
		return
	}
}
func (resource *Controller[Type]) UpdateOne()  {}
func (resource *Controller[Type]) DeleteOne()  {}
func (resource *Controller[Type]) ReplaceOne() {}

func (resource *Controller[Type]) FindAll() {

	findModel := resource.Controller.FindAllModel()

	if findModel != nil && resource.Request.Body != nil {
		bindErr := resource.BindJSON(findModel)
		if bindErr != nil {
			resource.sendError(http.StatusBadRequest, "error parsing the parameters", nil)
			return
		}
	}

	page, _ := strconv.Atoi(resource.Parameters.ByName("page"))

	perPage := resource.perPage
	if perPage == 0 {
		perPage = 20
	}

	context := SearchContext{}
	if resource.Request.Body != nil {
		bindErr := resource.BindJSON(&context)
		if bindErr != nil {
			resource.sendError(http.StatusBadRequest, "error parsing the parameters", nil)
			return
		}
	}
	context.Page = page
	context.PerPage = perPage

	models, total, result, err := resource.Controller.FindAll(context)
	if err != nil {
		resource.sendError(http.StatusInternalServerError, "error processing the request", nil)
		return
	}
	if result.HasErrors() {
		resource.sendError(http.StatusBadRequest, "error parsing the parameters", result)
	}

	err = json.NewEncoder(resource.Response).Encode(&findAllResponse{
		Records: models,
		Page:    page,
		PerPage: perPage,
		Total:   total,
	})
	if err != nil {
		resource.sendError(http.StatusInternalServerError, "error processing the request", nil)
		return
	}
}

func (resource *Controller[Type]) FindOne() {
	findModel := resource.Controller.FindOneModel()
	if findModel != nil {
		bindErr := resource.BindGetForm(findModel)
		if bindErr != nil {
			resource.sendError(http.StatusBadRequest, "error parsing the parameters", nil)
			return
		}
	}
	models, result, err := resource.Controller.FindOne(resource.Parameters.ByName("resourceId"))
	if err != nil {
		resource.sendError(http.StatusInternalServerError, "error processing the request", nil)
		return
	}
	if result.HasErrors() {
		resource.sendError(http.StatusBadRequest, "error parsing the parameters", result)
	}

	err = json.NewEncoder(resource.Response).Encode(models)
	if err != nil {
		resource.sendError(http.StatusInternalServerError, "error processing the request", nil)
		return
	}
}
