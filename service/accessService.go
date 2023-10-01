package service

type AccessService interface {
	CreateToken() (*model.TokenDetails, error)
	RefreshToken()
}

type accessService struct{}

func (s *accessService) CreateToken() (*model.TokenDetails, error) {

}



func NewAccessService() AccessService {
	return &accessService{}
}
