package model

var (
	AccessTokenTime  int64 = 24
	RefreshTokenTime int64 = 72
)

var MapModelType = map[string]interface{}{
	"users":               []User{},
	"roles":               []Role{},
	"userRole":            []UserRole{},
	"profiles":            []Profile{},
	"companies":           []Company{},
	"internships":         []InternShip{},
	"internshipEvaluates": []InternShipEvaluate{},
	"internJobs":          []InternJob{},
	"recruitments":        []Recruitment{},
}

var MapAssociation = map[string]map[string]interface{}{ // Alown preload association 2 level model
	"users": {
		"UserRoles": "",
	},
	"roles":    {},
	"userRole": {},
	"profiles": {
		"User":        "",
		"Recruitment": "",
		"InternShip":  "",
	},
}
