package api

type Service interface {
	Start() error
	Stop() error
}
