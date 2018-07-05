package services

type Service interface {
	Start(string, *[]string)
}
