package contract

const IDKey = "goweb:id"

type IDService interface {
	NewID() string
}