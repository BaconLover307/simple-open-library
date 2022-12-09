package service

type PickupService interface {
	Schedule()
	Update()
	Delete()
	FindById()
	FindAll()
}