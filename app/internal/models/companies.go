package models

type GetCompanyRequest struct {
	Name    string `schema:"name"`
	Code    string `schema:"code"`
	Country string `schema:"country"`
	Website string `schema:"website"`
	Phone   string `schema:"phone"`
}

type Company struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Country string `json:"country"`
	Website string `json:"website"`
	Phone   string `json:"phone"`
}
