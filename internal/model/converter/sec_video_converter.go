package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
)

func SecVideoToResponse(courseSec *entity.SectionVideo) *model.SecVideoResponse {

	return &model.SecVideoResponse{
		UUID:      courseSec.UUID,
		Title:     courseSec.Title,
		Notes:     courseSec.Notes,
		Sequence:  courseSec.Sequence,
		CreatedAt: courseSec.CreatedAt,
		UpdatedAt: courseSec.UpdatedAt,
	}
}
