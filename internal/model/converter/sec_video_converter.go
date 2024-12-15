package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
)

func SecVideoToResponse(secVideo *entity.SectionVideo) *model.SecVideoResponse {

	return &model.SecVideoResponse{
		UUID:         secVideo.UUID,
		Title:        secVideo.Title,
		Notes:        secVideo.Notes,
		Sequence:     secVideo.Sequence,
		OriginalName: secVideo.OriginalName,
		OriginalSize: secVideo.OriginalSize,
		OriginalMime: secVideo.OriginalMime,
		Bucket:       secVideo.Bucket,
		Dir:          secVideo.Dir,
		CreatedAt:    secVideo.CreatedAt,
		UpdatedAt:    secVideo.UpdatedAt,
	}
}
