package serializers

import (
	"github.com/christo-andrew/haven/internal/api/responses"
	"github.com/christo-andrew/haven/internal/models"
	"github.com/christo-andrew/haven/pkg"
)

type UserSerializer struct {
	Data interface{}
	many bool
}

func NewUserSerializer(data interface{}, many bool) UserSerializer {
	return UserSerializer{
		Data: data,
		many: many,
	}
}

func (us UserSerializer) Serialize() (interface{}, error) {
	switch us.Data.(type) {
	case []models.User:
		data, ok := us.Data.([]responses.CreateUserResponse)
		if !ok {
			return nil, pkg.InvalidDataError()
		}
		return data, nil
	case models.User:
		data, err := us.serializeSingle(us.Data)
		if err != nil {
			return nil, err
		}

		return data, nil
	default:
		return nil, pkg.InvalidDataError()
	}
}

func (us UserSerializer) serializeSingle(obj interface{}) (interface{}, error) {
	user, ok := obj.(models.User)
	if !ok {
		return nil, pkg.InvalidDataError()
	}

	return &responses.CreateUserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.GetFullName(),
	}, nil
}
