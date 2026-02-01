package models

import "database/sql/driver"

type StringEnum string

func (e *StringEnum) Scan(value interface{}) error {
	*e = StringEnum(value.(string))
	return nil
}

func (e StringEnum) Value() (driver.Value, error) {
	return string(e), nil
}

// Entity Enums
type EntityType string

const (
	EntityPerson       EntityType = "person"
	EntityCompany      EntityType = "company"
	EntityDealer       EntityType = "dealer"
	EntityOrganization EntityType = "organization"
)

type EntityStatus string

const (
	StatusActive    EntityStatus = "active"
	StatusInactive  EntityStatus = "inactive"
	StatusSuspended EntityStatus = "suspended"
)

type PhoneType string

const (
	PhoneMobile PhoneType = "mobile"
	PhoneHome   PhoneType = "home"
	PhoneWork   PhoneType = "work"
	PhoneFax    PhoneType = "fax"
	PhoneOther  PhoneType = "other"
)

// Access Enums
type AccessScope string

const (
	ScopeAll  AccessScope = "all"
	ScopeOwn  AccessScope = "own"
	ScopeTeam AccessScope = "team"
	ScopeNone AccessScope = "none"
)

// Location & Inventory Enums
type LocationType string

const (
	LocSalesLot  LocationType = "sales_lot"
	LocStorage   LocationType = "storage"
	LocService   LocationType = "service"
	LocOffsite   LocationType = "offsite"
	LocInTransit LocationType = "in_transit"
)

type BodyType string

const (
	BodySedan       BodyType = "sedan"
	BodySUV         BodyType = "suv"
	BodyTruck       BodyType = "truck"
	BodyCoupe       BodyType = "coupe"
	BodyVan         BodyType = "van"
	BodyHatchback   BodyType = "hatchback"
	BodyConvertible BodyType = "convertible"
	BodyWagon       BodyType = "wagon"
)

type VehicleCondition string

const (
	CondNew       VehicleCondition = "new"
	CondUsed      VehicleCondition = "used"
	CondCertified VehicleCondition = "certified"
)

type VehicleStatus string

const (
	VStatusInTransit    VehicleStatus = "in_transit"
	VStatusInRecon      VehicleStatus = "in_recon"
	VStatusReadyForSale VehicleStatus = "ready_for_sale"
	VStatusSold         VehicleStatus = "sold"
	VStatusWholesale    VehicleStatus = "wholesale"
)

type LotType string

const (
	LotNew       LotType = "new"
	LotUsed      LotType = "used"
	LotCPO       LotType = "cpo"
	LotWholesale LotType = "wholesale"
)

type AcquisitionSource string

const (
	SourceFactory        AcquisitionSource = "factory"
	SourceTradeIn        AcquisitionSource = "trade_in"
	SourceAuction        AcquisitionSource = "auction"
	SourceDealerTransfer AcquisitionSource = "dealer_transfer"
	SourceConsignment    AcquisitionSource = "consignment"
)

type PhotoPerspective string

const (
	PerspFront         PhotoPerspective = "front"
	PerspRear          PhotoPerspective = "rear"
	PerspLeftSide      PhotoPerspective = "left_side"
	PerspRightSide     PhotoPerspective = "right_side"
	PerspInteriorFront PhotoPerspective = "interior_front"
	PerspInteriorRear  PhotoPerspective = "interior_rear"
	PerspDashboard     PhotoPerspective = "dashboard"
	PerspEngine        PhotoPerspective = "engine"
	PerspTrunk         PhotoPerspective = "trunk"
	PerspDetail        PhotoPerspective = "detail"
	PerspDamage        PhotoPerspective = "damage"
	PerspOther         PhotoPerspective = "other"
)

type PhotoPurpose string

const (
	PurposeListing    PhotoPurpose = "listing"
	PurposeService    PhotoPurpose = "service"
	PurposeInspection PhotoPurpose = "inspection"
	PurposeDamage     PhotoPurpose = "damage"
	PurposeInternal   PhotoPurpose = "internal"
)

// Lead Enums
type StepStatus string

const (
	StepPending   StepStatus = "pending"
	StepCompleted StepStatus = "completed"
	StepSkipped   StepStatus = "skipped"
	StepFailed    StepStatus = "failed"
)

type ActivityType string

const (
	ActCallOutbound         ActivityType = "call_outbound"
	ActCallInbound          ActivityType = "call_inbound"
	ActEmailSent            ActivityType = "email_sent"
	ActEmailReceived        ActivityType = "email_received"
	ActSMSSent              ActivityType = "sms_sent"
	ActSMSReceived          ActivityType = "sms_received"
	ActAppointmentScheduled ActivityType = "appointment_scheduled"
	ActAppointmentCompleted ActivityType = "appointment_completed"
	ActAppointmentCancelled ActivityType = "appointment_cancelled"
	ActDemo                 ActivityType = "demo"
	ActQuoteSent            ActivityType = "quote_sent"
	ActOther                ActivityType = "other"
)