package output

// import "torque-dms/core/identity/domain"

type Phone struct {
	ID        uint
	EntityID  uint
	CountryID uint
	Number    string
	Type      string
	IsPrimary bool
	Verified  bool
}

type PhoneRepository interface {
	Save(phone *Phone) error
	Update(phone *Phone) error
	FindByID(id uint) (*Phone, error)
	FindByEntityID(entityID uint) ([]*Phone, error)
	FindPrimaryByEntityID(entityID uint) (*Phone, error)
	Delete(id uint) error
}