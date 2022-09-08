package auth

type Action int

const (
	ActionCompanyCreate = iota + 1
	ActionCompanyDelete
)

func (s *service) IsActionAllowed(action Action, ip string) (bool, error) {
	switch action {
	case ActionCompanyCreate, ActionCompanyDelete:
		country, err := s.client.GetRequestLocation(ip)
		if err != nil {
			return false, err
		}
		return isCountryAllowed(action, country), nil
	default:
		return false, nil
	}
}

func isCountryAllowed(action Action, country string) bool {
	switch action {
	case ActionCompanyCreate, ActionCompanyDelete:
		return country == "CY"
	}
	return false
}
