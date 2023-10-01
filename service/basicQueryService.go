package service

import (
	"fmt"
	"phenikaa/infrastructure"
	"phenikaa/model"
	"phenikaa/utils"
	"reflect"

	"gorm.io/gorm"
)

type BasicQueryService interface {
	Upsert(payload model.BasicQueryPayload) (interface{}, error)
	Delete(payload model.ListModelId) error
}

type basicQueryService struct{}

func (s *basicQueryService) Upsert(payload model.BasicQueryPayload) (interface{}, error) {
	var db = infrastructure.GetDB()
	var modelType = model.MapModelType[payload.ModelType]
	// var modelPayload = payload.Data

	var listModelId []uint
	if err := db.Model(modelType).Pluck("id", &listModelId).Error; err != nil {
		return nil, fmt.Errorf("get list id error: %v", err)
	}

	// Upsert multiple
	var listModelCreate []interface{}
	var listModelUpdate []interface{}
	if reflect.TypeOf(payload.Data).Kind() == reflect.Slice || reflect.TypeOf(payload.Data).Elem().Kind() == reflect.Slice {
		for _, data := range payload.Data.([]interface{}) {
			data := data.(map[string]interface{})
			if ok, _ := utils.InArray(data["id"].(uint), listModelId); ok {
				listModelUpdate = append(listModelUpdate, data)
			} else {
				listModelCreate = append(listModelCreate, data)
			}
		}

		if err := db.Transaction(func(tx *gorm.DB) error {
			if len(listModelCreate) > 0 {
				if err := tx.Model(modelType).Create(listModelCreate).Error; err != nil {
					return fmt.Errorf("create error: %v", err)
				}
			}

			if len(listModelUpdate) > 0 {
				if err := tx.Model(modelType).Updates(listModelUpdate).Error; err != nil {
					return fmt.Errorf("update error: %v", err)
				}
			}
			return nil
		}); err != nil {
			return nil, fmt.Errorf("upsert error: %v", err)
		}
		goto End
	}

	// Upsert single
	if ok, _ := utils.InArray(payload.Data.(map[string]interface{})["id"].(uint), listModelId); ok {
		if err := db.Model(modelType).Updates(payload.Data).Error; err != nil {
			return nil, fmt.Errorf("update error: %v", err)
		}
	} else {
		if err := db.Model(modelType).Create(payload.Data).Error; err != nil {
			return nil, fmt.Errorf("create error: %v", err)
		}
	}

	if len(payload.ModelType) == 0 {
		return nil, nil
	}

End:
	return payload.Data, nil
}

func (s *basicQueryService) Delete(payload model.ListModelId) error {
	var db = infrastructure.GetDB()
	var modelType = model.MapModelType[payload.ModelType]
	if err := db.Where("id IN ?", payload.ID).Delete(modelType).Error; err != nil {
		return fmt.Errorf("Delete error: %v", err)
	}

	return nil
}

func NewBasicQueryService() BasicQueryService {
	return &basicQueryService{}
}
