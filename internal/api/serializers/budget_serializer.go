package serializers

import (
	"github.com/christo-andrew/haven/internal/api/responses"
	"github.com/christo-andrew/haven/internal/models"
	"github.com/christo-andrew/haven/pkg/errors"
)

type BudgetSerializer struct {
	Data interface{}
	many bool
}

func NewBudgetSerializer(data interface{}, many bool) *BudgetSerializer {
	return &BudgetSerializer{
		Data: data,
		many: many,
	}
}

func (bs BudgetSerializer) Serialize() (interface{}, error) {
	switch bs.Data.(type) {
	case []models.Budget:
		return bs.serializeMany(bs.Data)
	case models.Budget:
		return bs.serializeSingle(bs.Data)
	default:
		return nil, errors.InvalidDataError()
	}
}

func (bs BudgetSerializer) serializeSingle(obj interface{}) (*responses.BudgetResponse, error) {
	if obj == nil {
		return nil, errors.InvalidDataError()
	}
	budget, ok := obj.(models.Budget)
	if !ok {
		return nil, errors.InvalidDataError()
	}

	return responses.BudgetResponse{}.FromBudget(budget), nil
}

func (bs BudgetSerializer) serializeMany(obj interface{}) (interface{}, error) {
	budgets, ok := obj.([]models.Budget)
	if !ok {
		return nil, errors.InvalidDataError()
	}

	var response []*responses.BudgetResponse
	for _, budget := range budgets {
		data, err := bs.serializeSingle(budget)
		if err != nil {
			return nil, err
		}
		response = append(response, data)
	}
	return response, nil
}
