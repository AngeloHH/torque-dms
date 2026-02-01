package models

import "time"

type Entity struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	Type           EntityType   `json:"type"`
	FirstName      string       `json:"first_name"`
	LastName       string       `json:"last_name"`
	BusinessName   string       `json:"business_name"`
	TaxID          string       `json:"tax_id"`
	Email          string       `json:"email"`
	Address        string       `json:"address"`
	City           string       `json:"city"`
	State          string       `json:"state"`
	Zip            string       `json:"zip"`
	CountryID      *uint        `json:"country_id"`
	Country        *Country     `gorm:"foreignKey:CountryID" json:"country,omitempty"`
	IsSystemUser   bool         `gorm:"default:false" json:"is_system_user"`
	IsInternal     bool         `gorm:"default:false" json:"is_internal"`
	ParentEntityID *uint        `json:"parent_entity_id"`
	ParentEntity   *Entity      `gorm:"foreignKey:ParentEntityID" json:"parent_entity,omitempty"`
	Status         EntityStatus `gorm:"default:'active'" json:"status"`
	CreatedAt      time.Time    `json:"created_at"`
	ModifiedAt     time.Time    `json:"modified_at"`
}

type UserAccount struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	EntityID     uint         `gorm:"unique" json:"entity_id"`
	Entity       Entity       `gorm:"foreignKey:EntityID" json:"entity"`
	Username     string       `gorm:"unique" json:"username"`
	PasswordHash string       `json:"-"`
	LastLogin    time.Time    `json:"last_login"`
	Status       EntityStatus `gorm:"default:'active'" json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
}

type EntityPhone struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EntityID  uint      `json:"entity_id"`
	Entity    Entity    `gorm:"foreignKey:EntityID" json:"-"`
	CountryID uint      `json:"country_id"`
	Country   Country   `gorm:"foreignKey:CountryID" json:"country"`
	Number    string    `json:"number"`
	Extension string    `json:"extension"`
	Type      PhoneType `json:"type"`
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`
	Verified  bool      `gorm:"default:false" json:"verified"`
	CreatedAt time.Time `json:"created_at"`
}

type Resource struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Code           string    `gorm:"unique" json:"code"`
	Name           string    `json:"name"`
	URLPattern     string    `json:"url_pattern"`
	Method         string    `json:"method"`
	Module         string    `json:"module"`
	OwnershipField string    `json:"ownership_field"`
	CreatedAt      time.Time `json:"created_at"`
}

type Role struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"unique" json:"name"`
	Description  string    `json:"description"`
	IsSystemRole bool      `gorm:"default:false" json:"is_system_role"`
	CreatedAt    time.Time `json:"created_at"`
}

type RoleResource struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	RoleID     uint        `json:"role_id"`
	Role       Role        `gorm:"foreignKey:RoleID" json:"-"`
	ResourceID uint        `json:"resource_id"`
	Resource   Resource    `gorm:"foreignKey:ResourceID" json:"resource"`
	Scope      AccessScope `gorm:"default:'none'" json:"scope"`
	CreatedAt  time.Time   `json:"created_at"`
}

type EntityResource struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	EntityID   uint        `json:"entity_id"`
	Entity     Entity      `gorm:"foreignKey:EntityID" json:"-"`
	ResourceID uint        `json:"resource_id"`
	Resource   Resource    `gorm:"foreignKey:ResourceID" json:"resource"`
	Scope      AccessScope `gorm:"default:'none'" json:"scope"`
	AssignedBy uint        `json:"assigned_by"`
	Reason     string      `json:"reason"`
	ExpiresAt  *time.Time  `json:"expires_at"`
	CreatedAt  time.Time   `json:"created_at"`
}

type EntityRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EntityID  uint      `json:"entity_id"`
	Entity    Entity    `gorm:"foreignKey:EntityID" json:"-"`
	RoleID    uint      `json:"role_id"`
	Role      Role      `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt time.Time `json:"created_at"`
}