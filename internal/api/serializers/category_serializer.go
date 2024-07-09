package serializers

import (
	"github.com/christo-andrew/haven/internal/api/responses"
	"github.com/christo-andrew/haven/internal/models"
	"github.com/christo-andrew/haven/pkg/errors"
)

type CategorySerializer struct {
	Data interface{}
	many bool
}

func (cs CategorySerializer) Serialize() (interface{}, error) {
	switch cs.Data.(type) {
	case []models.Category:
		return cs.serializeMany()
	case models.Category:
		return cs.serializeSingle(cs.Data)
	default:
		return nil, errors.InvalidDataError()
	}
}

func NewCategorySerializer(data interface{}, many bool) CategorySerializer {
	return CategorySerializer{
		Data: data,
		many: many,
	}
}

func (cs CategorySerializer) serializeSingle(obj interface{}) (interface{}, error) {
	category, ok := obj.(models.Category)
	if !ok {
		return nil, errors.InvalidDataError()
	}

	return &responses.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Context:     category.Context,
		ContextType: category.ContextType,
	}, nil

}

func (cs CategorySerializer) serializeMany() (interface{}, error) {
	var response []*responses.CategoryResponse
	categories, ok := cs.Data.([]models.Category)
	if !ok {
		return nil, errors.InvalidDataError()
	}
	for _, category := range categories {
		response = append(response, &responses.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Context:     category.Context,
			ContextType: category.ContextType,
		})
	}
	return response, nil
}
