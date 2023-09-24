package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"pdt-phenikaa-htdn/repository"
	"pdt-phenikaa-htdn/utils"
	"strconv"
	"time"

	"gorm.io/gorm/clause"

	"pdt-phenikaa-htdn/infrastructure"
	"pdt-phenikaa-htdn/model"

	"github.com/go-chi/jwtauth"
	"github.com/go-ldap/ldap"
	"github.com/twinj/uuid"
	"gorm.io/gorm"
)

var (
	infoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Llongfile)
	errLog  = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

type accessService struct {
	advanceFilterRepo repository.AdvanceFilterRepo
	basicQueryRepo    repository.BasicQueryRepository
	db                *gorm.DB
}

// AccessService service access
type AccessService interface {
	Login()
	CreateToken(userID uint, role, departmentCode string) (*model.TokenDetail, error)
	CreateTokenWithUsername(db *gorm.DB, username string, role string) (*model.TokenDetail, error)
	CreateLDAPToken(username, password, code, email, role string) (string, error)
	CreateAuth(userID uint, tokenDetail *model.TokenDetail) error
	ExtractTokenMetadata(r *http.Request) (*model.AccessDetail, error)
	FetchAuth(authD *model.AccessDetail) (int, error)
	DeleteAuth(givenUUID string) (int64, error)
	ClearAuth(userID uint) error
	CookieKiller(w http.ResponseWriter, fullDomain, domain string)
	CreateSessionUserRole(username string, role string) error
	GetUserPermission(username string, roleCode string) (*model.UserPermission, error)
	QueryUserPermission(username string, roleCode string) (*model.UserPermission, error)
	GetUserRoles(username string) ([]model.UserPermissionRole, error)
	GetUserRolesString(username string) ([]string, error)
	TryLDAP(username string, password string) (bool, *ldap.SearchResult, error)
	ApplyUserRoles() error
	ResetRoles(codes []string) error
}

func (s *accessService) Login() {

}

func (s *accessService) CreateSessionUserRole(username string, role string) error {
	userType, err := s.basicQueryRepo.GetUserType(role)
	if err != nil {
		return err
	}
	userTableName := userType + "s"
	userTypeTableName := userType + "_roles"
	var userId uint

	if err := s.db.Debug().Table(userTableName).Select("id").Where("username = ?", username).Find(&userId).Error; err != nil {
		return err
	}
	if role == "" {
		// user_managers
		var resultRole string
		queryString := "SELECT role FROM managers WHERE profile_code = (SELECT code FROM profiles WHERE user_id = " + strconv.Itoa(int(userId)) + ")"
		if err := s.db.Debug().Raw(queryString).Find(&resultRole).Error; err != nil {
			return errors.New("unauthorized")
		}
		role = resultRole
	}
	var roleId uint
	if err := s.db.Debug().Model(&model.Role{}).Select("id").Where("code = ?", role).Find(&roleId).Error; err != nil {
		return err
	}

	sessionRecord := struct {
		UserId uint
		RoleId uint
		Active *bool
	}{
		UserId: userId,
		RoleId: roleId,
		Active: &model.TrueValue,
	}

	if err := s.db.Debug().Table(userTypeTableName).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "role_id"}},
		DoUpdates: clause.Set{clause.Assignment{
			Column: clause.Column{Name: "active"},
			Value:  true,
		}}}).Create(&sessionRecord).Error; err != nil {
		return err
	}
	if err := s.db.Debug().Table(userTypeTableName).Where("user_id = ? AND role_id NOT IN (?)", userId, roleId).Update("active", false).Error; err != nil {
		return err
	}
	return nil
}

func (s *accessService) CreateTokenWithUsername(db *gorm.DB, username string, role string) (*model.TokenDetail, error) {
	// Create token details
	tokenDetail := &model.TokenDetail{}
	var err error
	// var userManager *model.UserManager
	// err := db.Model(&model.UserManager{}).Where("username = ?", username).Find(&userManager).Error
	// if err != nil {
	// 	return nil, fmt.Errorf("can't create token. UserManager not found")
	// }

	var profile model.Employee
	if err := s.db.Model(&model.Employee{}).Select("code", "user_id").Where("email = ?", username).Find(&profile).Error; err != nil {
		return nil, err
	}

	userID := profile.UserID
	tokenDetail.Username = username
	tokenDetail.AtExpires = time.Now().Add(time.Hour * time.Duration(model.AccessTokenTime)).Unix()
	tokenDetail.AccessUUID = utils.GetPattern(userID) + uuid.NewV4().String()
	tokenDetail.RtExpires = time.Now().Add(time.Hour * time.Duration(model.RefreshTokenTime)).Unix()
	tokenDetail.RefreshUUID = utils.GetPattern(userID) + uuid.NewV4().String()

	atClaims := make(map[string]interface{})
	atClaims["username"] = tokenDetail.Username
	atClaims["access_uuid"] = tokenDetail.AccessUUID
	atClaims["code"] = profile.Code
	atClaims["user_id"] = userID
	atClaims["role"] = role
	atClaims["exp"] = tokenDetail.AtExpires
	_, tokenDetail.AccessToken, err = infrastructure.GetEncodeAuth().Encode(atClaims)
	if err != nil {
		return nil, err
	}

	// Create Refresh Token
	rtClaims := make(map[string]interface{})
	rtClaims["username"] = tokenDetail.Username
	rtClaims["refresh_uuid"] = tokenDetail.RefreshUUID
	rtClaims["code"] = profile.Code
	rtClaims["user_id"] = userID
	rtClaims["role"] = role
	rtClaims["exp"] = tokenDetail.RtExpires
	_, tokenDetail.RefreshToken, err = infrastructure.GetEncodeAuth().Encode(rtClaims)
	if err != nil {
		return nil, err
	}

	return tokenDetail, nil
}

func (s *accessService) CreateToken(userID uint, role, departmentCode string) (*model.TokenDetail, error) {
	var err error
	// Create token details
	tokenDetail := &model.TokenDetail{}

	if userID != 0 {
		user, _, err := s.advanceFilterRepo.AdvanceFilter("users", "", &model.User{ID: userID}, model.Pagination{Page: 1, PageSize: -1}, false, false, false, []string{}, []string{"*"}, []string{"*"}, []model.Sort{}, []model.TimeFilter{}, model.DevClaims)
		if err != nil || len(user.([]model.User)) <= 0 {
			return nil, fmt.Errorf("can't create token. User not found")
		}
		tokenDetail.Username = user.([]model.User)[0].Username
		tokenDetail.AtExpires = time.Now().Add(time.Hour * time.Duration(infrastructure.GetExtendAccessHour())).Unix()
		tokenDetail.AccessUUID = utils.GetPattern(userID) + uuid.NewV4().String()
		tokenDetail.RtExpires = time.Now().Add(time.Hour * time.Duration(infrastructure.GetExtendRefreshHour())).Unix()
		tokenDetail.RefreshUUID = utils.GetPattern(userID) + uuid.NewV4().String()
	} else {
		tokenDetail.Username = "manager"
		tokenDetail.AtExpires = time.Now().Add(time.Hour * time.Duration(infrastructure.GetExtendAccessHour())).Unix()
		tokenDetail.AccessUUID = utils.GetPattern(userID) + uuid.NewV4().String()
		tokenDetail.RtExpires = time.Now().Add(time.Hour * time.Duration(infrastructure.GetExtendRefreshHour())).Unix()
		tokenDetail.RefreshUUID = utils.GetPattern(userID) + uuid.NewV4().String()
	}

	// Create Access Token
	atClaims := make(map[string]interface{})
	// atClaims["username"] = tokenDetail.Username
	atClaims["access_uuid"] = tokenDetail.AccessUUID
	atClaims["user_id"] = userID
	atClaims["role"] = role
	atClaims["department_code"] = departmentCode
	atClaims["exp"] = tokenDetail.AtExpires

	_, tokenDetail.AccessToken, err = infrastructure.GetEncodeAuth().Encode(atClaims)
	if err != nil {
		return nil, err
	}

	// Create Refresh Token
	rtClaims := make(map[string]interface{})
	rtClaims["username"] = tokenDetail.Username
	rtClaims["refresh_uuid"] = tokenDetail.RefreshUUID
	rtClaims["user_id"] = userID
	rtClaims["role"] = role
	rtClaims["department_code"] = departmentCode
	rtClaims["exp"] = tokenDetail.RtExpires
	_, tokenDetail.RefreshToken, err = infrastructure.GetEncodeAuth().Encode(rtClaims)
	if err != nil {
		return nil, err
	}

	return tokenDetail, nil
}

func (s *accessService) CreateAuth(userID uint, tokenDetail *model.TokenDetail) error {
	// converting Unix to UTC(to Time Object)
	accessToken := time.Unix(tokenDetail.AtExpires, 0)
	refreshToken := time.Unix(tokenDetail.RtExpires, 0)
	now := time.Now()

	if errAccess := infrastructure.
		GetRedisClient().
		Set(tokenDetail.AccessUUID, strconv.Itoa(int(userID)), accessToken.Sub(now)).
		Err(); errAccess != nil {
		return errAccess
	}

	if errRefresh := infrastructure.
		GetRedisClient().
		Set(tokenDetail.RefreshUUID, strconv.Itoa(int(userID)), refreshToken.Sub(now)).
		Err(); errRefresh != nil {
		return errRefresh
	}

	return nil
}

func (s *accessService) ExtractTokenMetadata(r *http.Request) (*model.AccessDetail, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return nil, err
	}

	accessUUID, ok := claims["access_uuid"].(string)
	if !ok {
		errLog.Println("can't parse access uuid from token")
		return nil, errors.New("can't parse access uuid from token")
	}

	userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
	if err != nil {
		errLog.Println(err)
		return nil, err
	}

	return &model.AccessDetail{
		AccessUUID: accessUUID,
		UserID:     uint(userID),
	}, nil
}

func (s *accessService) CreateLDAPToken(username, password, code, email, role string) (string, error) {
	var err error
	// Create Access Token
	encryptPassword, err := infrastructure.RsaEncrypt(password)
	if err != nil {
		return "", err
	}
	ldapTokenTime := time.Now().Add(time.Hour * time.Duration(model.LDAPTokenTime)).Unix()
	ldapUUID := uuid.NewV4().String()
	atClaims := make(map[string]interface{})
	atClaims["role"] = role
	atClaims["ldap_username"] = username
	atClaims["ldap_password"] = encryptPassword
	atClaims["code"] = code
	atClaims["email"] = email
	atClaims["ldap_uuid"] = ldapUUID
	atClaims["exp"] = ldapTokenTime

	_, tokenString, err := infrastructure.GetEncodeAuth().Encode(atClaims)
	if err != nil {
		return "", err
	}
	ldapToken := time.Unix(ldapTokenTime, 0)
	now := time.Now()

	if errLdap := infrastructure.
		GetRedisClient().
		Set(ldapUUID, code, ldapToken.Sub(now)).
		Err(); errLdap != nil {
		return "", errLdap
	}

	return tokenString, nil
}

// FetchAuth Looking up the token metadata in redis
func (s *accessService) FetchAuth(authD *model.AccessDetail) (int, error) {
	userIDStr, err := infrastructure.GetRedisClient().Get(authD.AccessUUID).Result()
	if err != nil {
		errLog.Println(err)
		return 0, err
	}

	userID, _ := strconv.ParseUint(userIDStr, 10, 64)
	return int(userID), nil
}

// DeleteAuth deleting the jwt metadata from redis store
func (s *accessService) DeleteAuth(givenUUID string) (int64, error) {
	deleted, err := infrastructure.GetRedisClient().Del(givenUUID).Result()
	if err != nil {
		return 0, err
	}

	return deleted, nil
}

func (s *accessService) ClearAuth(userID uint) error {
	list, _ := infrastructure.GetRedisClient().Keys(utils.PatternGet(userID)).Result()
	for _, key := range list {
		infrastructure.GetRedisClient().Del(key)
	}

	return nil
}

func (s *accessService) CookieKiller(w http.ResponseWriter, fullDomain, domain string) {
	accessCookie := http.Cookie{
		Name:   "AccessToken",
		Domain: fullDomain,
		Path:   "/",
		Value:  "",
		// HttpOnly: true,
		// Secure:   true,
		MaxAge: -1,
	}
	refreshCookie := http.Cookie{
		Name:   "RefreshToken",
		Domain: fullDomain,
		Path:   "/",
		Value:  "",
		// HttpOnly: true,
		// Secure:   true,
		MaxAge: -1,
	}
	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)
}

func (s *accessService) GetActivatedUserRole(username string, roleCode string) (*model.ActivatedUserRole, error) {
	// 	key := "user_role::" + username
	var res *model.ActivatedUserRole = nil

	// 	client := infrastructure.GetRedisClient()
	// 	resStr, err := client.Get(key).Result()
	// 	if err != nil {
	// 		goto QUERY
	// 	}
	// 	err = json.Unmarshal([]byte(resStr), &res)
	// 	if err != nil {
	// 		goto QUERY
	// 	}
	// QUERY:
	var user *model.User
	if err := s.db.Debug().Model(&model.User{}).Where("username = ?", username).
		Find(&user).Error; err != nil {
		return nil, err
	}
	res = &model.ActivatedUserRole{
		UserActive: user.Active,
	}

	userRoles := []string{}
	roles := []model.Role{}
	if err := s.db.Debug().Model(&model.Role{}).
		Where("id IN (SELECT role_id FROM user_roles WHERE active = true AND user_id = (SELECT id FROM users WHERE lower(username) = lower(?)))", username).
		Find(&roles).Error; err != nil {
		return nil, err
	}
	for _, item := range roles {
		userRoles = append(userRoles, item.Code)
		res.Roles = append(res.Roles, model.UserPermissionRole{
			ID:   item.ID,
			Code: item.Code,
			Name: item.Name,
		})
	}
	if !utils.Contains(userRoles, roleCode) {
		return nil, errors.New("no role")
	}
	res.ActiveRole = roleCode
	// save to redis
	// jsonString, err := json.Marshal(res)
	// if err != nil {
	// 	return nil, err
	// }
	// if errAccess := client.
	// 	Set(key, jsonString, 72*time.Hour).
	// 	Err(); errAccess != nil {
	// 	return nil, errAccess
	// }

	return res, nil
}

func (s *accessService) QueryUserPermission(username string, roleCode string) (*model.UserPermission, error) {
	// RES
	res := &model.UserPermission{
		Roles:      []model.UserPermissionRole{}, // xxx
		Routes:     []string{},
		APIs:       []model.UserPermissionAPI{},
		ModelTypes: make(map[string]model.UserPemissionModelType),
	}

	// ROLES
	userRoles := []string{}
	MAP_ROLE := make(map[string]uint)

	dat, err := s.GetActivatedUserRole(username, roleCode)
	if err != nil {
		return nil, err
	}
	res.ActiveRole = dat.ActiveRole
	res.UserActive = dat.UserActive
	// ROLE _ ID
	for _, item := range dat.Roles {
		userRoles = append(userRoles, item.Code)
		res.Roles = append(res.Roles, model.UserPermissionRole{
			ID:   item.ID,
			Code: item.Code,
			Name: item.Name,
		})
		MAP_ROLE[item.Code] = item.ID
	}

	if roleCode == "" && len(userRoles) > 0 {
		roleCode = userRoles[0]
	}
	if !utils.Contains(userRoles, roleCode) {
		return nil, errors.New("Unauthorized")
	}

	roleID := MAP_ROLE[roleCode]

	// ROUTE
	routes := []string{}
	if err := s.db.Debug().Model(&model.Route{}).Select("url").
		Where("id IN (SELECT route_id FROM role_routes WHERE active = true AND role_id = ?)", roleID).
		Find(&routes).Error; err != nil {
		fmt.Println("dm1", err)
		return nil, err
	}
	// API
	apis := []model.API{}
	if err := s.db.Debug().Model(&model.API{}).Where("id IN (SELECT api_id FROM role_apis WHERE role_id = ? AND active = true)", roleID).Find(&apis).Error; err != nil {
		fmt.Println("dm2", err)
		return nil, err
	}
	resApis := []model.UserPermissionAPI{}
	for _, item := range apis {
		resApis = append(resApis, model.UserPermissionAPI{
			Url:    item.Url,
			Method: item.Method,
		})
	}
	// MODEL TYPE
	modelTypes := []model.ModelType{}
	if err := s.db.Debug().Model(&model.ModelType{}).Where("id IN (SELECT model_type_id FROM role_model_type_permissions WHERE role_id = ?)", roleID).Find(&modelTypes).Error; err != nil {
		fmt.Println("dm3", err)
		return nil, err
	}
	roleModelTypes := []model.RoleModelTypePermission{}
	if err := s.db.Model(&model.RoleModelTypePermission{}).
		Preload("Fields").
		Where("role_id = ?", roleID).
		Find(&roleModelTypes).Error; err != nil {
		fmt.Println("dm4", err)
		return nil, err
	}
	modelFields := []model.ModelTypeField{}
	if err := s.db.Model(&model.ModelTypeField{}).
		Where("model_type_id IN (SELECT model_type_id FROM role_model_type_permissions WHERE role_id = ?)", roleID).
		Find(&modelFields).Error; err != nil {
		fmt.Println("dm5", err)
		return nil, err
	}

	var MAP_FIELD_ID_CODE = make(map[uint]string)
	for _, item := range modelFields {
		MAP_FIELD_ID_CODE[item.ID] = item.Code
	}

	var MAP_ROLE_MODEL = make(map[uint]model.UserPermissonBoolean)
	var MAP_ROLE_FIELD = make(map[uint]map[string]model.UserPermissonBoolean, 0)

	for _, item := range roleModelTypes {
		MAP_ROLE_MODEL[item.ModelTypeId] = model.UserPermissonBoolean{
			Create: *item.CreatePermission,
			Read:   *item.ReadPermission,
			Update: *item.UpdatePermission,
			Delete: *item.DeletePermission,
		}
		MAP_ROLE_FIELD[item.ModelTypeId] = make(map[string]model.UserPermissonBoolean, 0)
		for _, field := range item.Fields {
			_code := MAP_FIELD_ID_CODE[field.FieldId]
			MAP_ROLE_FIELD[item.ModelTypeId][_code] = model.UserPermissonBoolean{
				Create: *field.CreatePermission,
				Read:   *field.ReadPermission,
				Update: *field.UpdatePermission,
				Delete: *field.DeletePermission,
			}
		}
	}

	for _, item := range modelTypes {
		res.ModelTypes[item.Code] = model.UserPemissionModelType{
			Create: MAP_ROLE_MODEL[item.ID].Create,
			Read:   MAP_ROLE_MODEL[item.ID].Read,
			Update: MAP_ROLE_MODEL[item.ID].Update,
			Delete: MAP_ROLE_MODEL[item.ID].Delete,
			Fields: MAP_ROLE_FIELD[item.ID],
		}
	}

	// RETURN
	res.Routes = routes
	res.APIs = resApis
	return res, nil
}

func (s *accessService) GetUserPermission(username string, roleCode string) (*model.UserPermission, error) {
	// check if exist in redis
	// 	key := "per::" + roleCode
	var res *model.UserPermission = nil

	// 	client := infrastructure.GetRedisClient()
	// 	resStr, err := client.Get(key).Result()
	// 	dat, err := s.GetActivatedUserRole(username, roleCode)
	// 	if err != nil {
	// 		goto QUERY
	// 	}
	// 	err = json.Unmarshal([]byte(resStr), &res)
	// 	if err != nil {
	// 		goto QUERY
	// 	}
	// 	res.UserActive = dat.UserActive
	// 	res.ActiveRole = roleCode
	// 	res.Roles = dat.Roles
	// 	return res, nil
	// 	// otherwise, query and save to redis for next time get permission
	// QUERY:
	res, err := s.QueryUserPermission(username, roleCode)
	if err != nil {
		return nil, err
	}
	// jsonString, err := json.Marshal(res)
	// if err != nil {
	// 	return nil, err
	// }
	// if errAccess := client.
	// 	Set(key, jsonString, 72*time.Hour).
	// 	Err(); errAccess != nil {
	// 	return nil, errAccess
	// }
	return res, nil
}

func (s *accessService) GetUserRoles(username string) ([]model.UserPermissionRole, error) {
	userType := "user"
	userTableName := userType + "s"
	userTypeTableName := userType + "_roles"

	roles := []model.Role{}
	if err := s.db.Debug().Model(&model.Role{}).
		Where("id IN (SELECT role_id FROM "+userTypeTableName+" WHERE active = true AND user_id = (SELECT id FROM "+userTableName+" WHERE username = ?))", username).
		Find(&roles).Error; err != nil {
		return nil, err
	}

	res := []model.UserPermissionRole{}
	// ROLE _ ID
	for _, item := range roles {
		res = append(res, model.UserPermissionRole{
			Code: item.Code,
			Name: item.Name,
		})
	}

	return res, nil
}

func (s *accessService) GetUserRolesString(username string) ([]string, error) {
	userType := "user"
	userTableName := userType + "s"
	userTypeTableName := userType + "_roles"

	roles := []string{}
	if err := s.db.Debug().Model(&model.Role{}).Select("code").
		Where("id IN (SELECT role_id FROM "+userTypeTableName+" WHERE active = true AND user_id = (SELECT id FROM "+userTableName+" WHERE username = ?))", username).
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *accessService) TryLDAP(username string, password string) (bool, *ldap.SearchResult, error) {
	l, err := utils.ConnectTLS()
	if err != nil {
		return false, nil, err
	}
	defer l.Close()
	sRes, err := utils.BindAndSearch(l, username)
	if err != nil || len(sRes.Entries) == 0 || len(sRes.Entries) > 1 {
		return false, nil, err
	}

	userDN := sRes.Entries[0].DN

	// Bind as the user to verify their password
	err = l.Bind(userDN, string(password))
	if err != nil {
		return false, nil, err
	}
	return true, sRes, nil
	// return true, nil, nil
}

func (s *accessService) ApplyUserRoles() error {
	// GET new data
	// dat := []model.User{}
	// if err := s.db.Debug().Model(&model.User{}).Find(&dat).Error; err != nil {
	// 	return err
	// }
	// stdIDs := []uint{}
	// for _, r := range dat {
	// 	stdIDs = append(stdIDs, *r.StudentId)
	// }
	// // GET MAP STUDENT_ID - USER_ID
	// students := []model.Student{}
	// if err := s.db.Debug().Model(&model.Student{}).Select("id", "user_id").Where("id IN (?)", stdIDs).Find(&students).Error; err != nil {
	// 	return err
	// }
	// mapStudent := map[uint]uint{}
	// for _, r := range students {
	// 	mapStudent[r.ID] = r.UserId
	// }
	// // GET MAP ROLE ID-CODE
	// roles := []model.Role{}
	// if err := s.db.Debug().Model(&model.Role{}).Where("code IN ('monitor', 'vice-monitor')").Find(&roles).Error; err != nil {
	// 	return err
	// }
	// mapRole := map[string]uint{}
	// for _, r := range roles {
	// 	mapRole[r.Code] = r.ID
	// }
	// // Deactive all current records
	// if err := s.db.Debug().Model(&model.UserRole{}).
	// 	Where("role_id IN (SELECT id FROM roles WHERE code IN ('monitor', 'vice-monitor') AND deleted_at IS NULL)").
	// 	Update("active", false).Error; err != nil {
	// 	return err
	// }

	// // Update/Create new data user-roles
	// newDat := []model.UserRole{}
	// mapAdded := []string{}
	// for _, item := range dat {
	// 	userId := mapStudent[*item.StudentId]
	// 	roleId := mapRole[item.ClassRole]
	// 	newKey := strconv.Itoa(int(userId)) + "-" + strconv.Itoa(int(roleId))
	// 	if utils.Contains(mapAdded, newKey) {
	// 		mapAdded = append(mapAdded, newKey)
	// 		newDat = append(newDat, model.UserRole{
	// 			ID:     0,
	// 			UserId: userId,
	// 			RoleId: roleId,
	// 			Active: &model.TrueValue,
	// 		})
	// 	}
	// }

	// if len(newDat) > 0 {
	// 	if err := s.db.Debug().Model(&model.UserRole{}).
	// 		Clauses(clause.OnConflict{
	// 			Columns:   []clause.Column{{Name: "user_id"}, {Name: "role_id"}},
	// 			DoUpdates: clause.AssignmentColumns([]string{"active"}),
	// 		}).
	// 		CreateInBatches(newDat, 500).Error; err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func (s *accessService) ResetRoles(codes []string) error {
	if len(codes) == 0 {
		list, _ := infrastructure.GetRedisClient().Keys("per::*").Result()
		for _, key := range list {
			_, err := infrastructure.GetRedisClient().Del(key).Result()
			if err != nil {
				return err
			}
		}
		// ======================
		list2, _ := infrastructure.GetRedisClient().Keys("map_query::*").Result()
		for _, key := range list2 {
			_, err := infrastructure.GetRedisClient().Del(key).Result()
			if err != nil {
				return err
			}
		}
		// ======================
		list3, _ := infrastructure.GetRedisClient().Keys("map_filter::*").Result()
		for _, key := range list3 {
			_, err := infrastructure.GetRedisClient().Del(key).Result()
			if err != nil {
				return err
			}
		}
	} else {
		for _, code := range codes {
			_, err := infrastructure.GetRedisClient().Del("per::" + code).Result()
			if err != nil {
				return err
			}
			_, err = infrastructure.GetRedisClient().Del("map_query::" + code).Result()
			if err != nil {
				return err
			}
			_, err = infrastructure.GetRedisClient().Del("map_filter::" + code).Result()
			if err != nil {
				return err
			}
		}
	}

	// ======================
	list, _ := infrastructure.GetRedisClient().Keys("user_role::*").Result()
	for _, key := range list {
		_, err := infrastructure.GetRedisClient().Del(key).Result()
		if err != nil {
			return err
		}
	}

	return nil
}

// NewAccessService export access service
func NewAccessService() AccessService {
	advanceFilterRepo := repository.NewAdvanceFilterRepo()
	basicQueryRepo := repository.NewBasicQueryRepo()
	db := infrastructure.GetDB()
	return &accessService{
		db:                db,
		basicQueryRepo:    basicQueryRepo,
		advanceFilterRepo: advanceFilterRepo}
}
