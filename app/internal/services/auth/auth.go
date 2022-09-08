package auth

type AuthService interface {
	IsActionAllowed(action Action, ip string) (bool, error)
}

type IPDataProvider interface {
	GetRequestLocation(ip string) (string, error)
}

type service struct {
	client IPDataProvider
}

func NewService(ipDataProvider IPDataProvider) *service {
	return &service{
		client: ipDataProvider,
	}
}
