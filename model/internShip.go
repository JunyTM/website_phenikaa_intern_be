package model

// Bảng quản lý thông tin báo cáo thực tập của sinh viên sinh
type InternShip struct {
	ID                   uint   `json:"id" gorm:"autoIncrement"`
	ProfileId            uint   `json:"profile_id"`
	CompanyId            uint   `json:"company_id"`
	InternShipEvaluateId uint   `json:"internship_evaluate_id"`
	Code                 string `json:"code"`

	Profile            *Profile            `json:"profile" gorm:"foreignKey:ProfileId"`
	Company            *Company            `json:"company" gorm:"foreignKey:CompanyId"`
	InternShipEvaluate *InternShipEvaluate `json:"internship_evaluate" gorm:"foreignKey:InternShipEvaluateId"`
}
