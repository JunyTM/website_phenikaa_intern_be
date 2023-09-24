package model

import "github.com/lib/pq"

type AdvanceFilterPayload struct {
	ModelType         string         `json:"modelType"`
	IgnoreAssociation bool           `json:"ignoreAssociation"`
	Page              int            `json:"page"`
	PageSize          int            `json:"pageSize"`
	IsPaginateDB      bool           `json:"isPaginateDB"`
	QuerySerch        string         `json:"querySearch"`
	SelectColumn      pq.StringArray `json:"selectColumn"`
}

type BasicQueryPayload struct {
	ModelType string      `json:"modelType"`
	Data      interface{} `json:"data"`
}

