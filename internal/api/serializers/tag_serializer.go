package serializers

import (
	"github.com/christo-andrew/haven/internal/api/responses"
	"github.com/christo-andrew/haven/internal/models"
)

type TagSerializer struct {
	Data interface{}
	many bool
}

func NewTagSerializer(data interface{}, many bool) *TagSerializer {
	return &TagSerializer{
		Data: data,
		many: many,
	}
}

func (ts TagSerializer) Serialize() interface{} {
	switch ts.Data.(type) {
	case []models.Tag:
		return ts.serializeTags()
	case models.Tag:
		return ts.serializeSingleTag()
	default:
		return nil
	}
}

func (ts TagSerializer) serializeTags() interface{} {
	response := make([]*responses.TagResponse, 0)
	for _, tag := range ts.Data.([]models.Tag) {
		response = append(response, serializeTag(tag))
	}
	return response
}

func serializeTag(tag models.Tag) *responses.TagResponse {
	return &responses.TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

func (ts TagSerializer) serializeSingleTag() *responses.TagResponse {
	return serializeTag(ts.Data.(models.Tag))
}
