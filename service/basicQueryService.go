package service

import (
	"fmt"
	"phenikaa/infrastructure"
	"phenikaa/model"
)

type BasicQueryService interface {
	Upsert(payload model.BasicQueryPayload) (interface{}, error)
	Delete(payload model.BasicQueryPayload) error
}

type basicQueryService struct{}

func (s *basicQueryService) Upsert(payload model.BasicQueryPayload) (interface{}, error) {
	var db = infrastructure.GetDB()
	var modelType = model.MapModelType[payload.ModelType]
	// var modelPayload = payload.Data

	var listModelId []uint
	if err := db.Model(modelType).Pluck("id", &listModelId).Error; err != nil {
		return nil, fmt.Errorf("Get list id error: %v", err)
	}

	// for index, data := range modelPayload {
	// }

	if len(payload.ModelType) == 0 {
		return nil, nil
	}

	return nil, nil
}

func (s *basicQueryService) Delete(payload model.BasicQueryPayload) error {
	return nil
}
