package model

var MapModelType = map[string]interface{}{
	"users": []User{},
}

var MapAssociation = map[string]map[string]interface{}{ // Alown preload association 2 level model
	"users": {},
}
