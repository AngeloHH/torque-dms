package domain

import (
	"testing"
	"time"
)

func validVehicle() *Vehicle {
	v, _ := NewVehicle("STK001", "12345678901234567", "Toyota", "Corolla", 2024)
	return v
}

func TestNewVehicle_Valid(t *testing.T) {
	v, err := NewVehicle("STK001", "12345678901234567", "Toyota", "Corolla", 2024)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.StockNumber != "STK001" {
		t.Errorf("StockNumber = %v, want STK001", v.StockNumber)
	}
	if v.Status != VehicleStatusInTransit {
		t.Errorf("Status = %v, want in_transit", v.Status)
	}
	if v.Condition != VehicleConditionUsed {
		t.Errorf("Condition = %v, want used", v.Condition)
	}
	if v.LotType != LotTypeUsed {
		t.Errorf("LotType = %v, want used", v.LotType)
	}
}

func TestNewVehicle_Invalid(t *testing.T) {
	tests := []struct {
		name   string
		stock  string
		vin    string
		make   string
		model  string
		year   int
	}{
		{"empty stock", "", "12345678901234567", "Toyota", "Corolla", 2024},
		{"empty VIN", "STK001", "", "Toyota", "Corolla", 2024},
		{"short VIN", "STK001", "12345", "Toyota", "Corolla", 2024},
		{"empty make", "STK001", "12345678901234567", "", "Corolla", 2024},
		{"empty model", "STK001", "12345678901234567", "Toyota", "", 2024},
		{"year too old", "STK001", "12345678901234567", "Toyota", "Corolla", 1800},
		{"year too future", "STK001", "12345678901234567", "Toyota", "Corolla", time.Now().Year() + 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVehicle(tt.stock, tt.vin, tt.make, tt.model, tt.year)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestVehicle_SetPricing(t *testing.T) {
	v := validVehicle()

	if err := v.SetPricing(35000, 30000, 33000); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.MSRP != 35000 {
		t.Errorf("MSRP = %v, want 35000", v.MSRP)
	}
}

func TestVehicle_SetPricing_Negative(t *testing.T) {
	v := validVehicle()
	if err := v.SetPricing(-1, 0, 0); err == nil {
		t.Error("expected error for negative MSRP")
	}
	if err := v.SetPricing(0, -1, 0); err == nil {
		t.Error("expected error for negative invoice")
	}
	if err := v.SetPricing(0, 0, -1); err == nil {
		t.Error("expected error for negative asking")
	}
}

func TestVehicle_SetAcquisition(t *testing.T) {
	v := validVehicle()
	err := v.SetAcquisition(AcquisitionSourceAuction, 25000, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.AcquisitionSource != AcquisitionSourceAuction {
		t.Errorf("source = %v, want auction", v.AcquisitionSource)
	}
}

func TestVehicle_SetAcquisition_NegativeCost(t *testing.T) {
	v := validVehicle()
	err := v.SetAcquisition(AcquisitionSourceAuction, -1, time.Now())
	if err == nil {
		t.Error("expected error for negative cost")
	}
}

func TestVehicle_SetCondition_New(t *testing.T) {
	v := validVehicle()
	v.SetCondition(VehicleConditionNew)

	if v.Condition != VehicleConditionNew {
		t.Errorf("Condition = %v, want new", v.Condition)
	}
	if v.LotType != LotTypeNew {
		t.Error("LotType should change to new when condition is new")
	}
}

func TestVehicle_SetCondition_Used(t *testing.T) {
	v := validVehicle()
	v.SetCondition(VehicleConditionUsed)
	if v.LotType != LotTypeUsed {
		t.Error("LotType should remain used")
	}
}

func TestVehicle_StatusTransitions(t *testing.T) {
	v := validVehicle()

	// in_transit -> ready_for_sale
	if err := v.MarkAsReadyForSale(); err != nil {
		t.Fatalf("MarkAsReadyForSale error: %v", err)
	}
	if !v.IsAvailable() {
		t.Error("should be available after ready_for_sale")
	}

	// ready_for_sale -> in_recon
	if err := v.SendToRecon(); err != nil {
		t.Fatalf("SendToRecon error: %v", err)
	}
	if v.Status != VehicleStatusInRecon {
		t.Errorf("Status = %v, want in_recon", v.Status)
	}

	// in_recon -> sold
	if err := v.MarkAsSold(); err != nil {
		t.Fatalf("MarkAsSold error: %v", err)
	}
	if !v.IsSold() {
		t.Error("should be sold")
	}

	// sold -> can't change
	if err := v.MarkAsSold(); err == nil {
		t.Error("expected error marking already sold vehicle as sold")
	}
	if err := v.MarkAsReadyForSale(); err == nil {
		t.Error("expected error changing status of sold vehicle")
	}
	if err := v.SendToRecon(); err == nil {
		t.Error("expected error sending sold vehicle to recon")
	}
}

func TestVehicle_ProfitAndMargin(t *testing.T) {
	v := validVehicle()
	v.AskingPrice = 30000
	v.AcquisitionCost = 25000

	if v.Profit() != 5000 {
		t.Errorf("Profit = %v, want 5000", v.Profit())
	}

	margin := v.Margin()
	if margin < 16.66 || margin > 16.67 {
		t.Errorf("Margin = %v, want ~16.67", margin)
	}
}

func TestVehicle_Margin_ZeroAskingPrice(t *testing.T) {
	v := validVehicle()
	if v.Margin() != 0 {
		t.Errorf("Margin should be 0 when asking price is 0, got %v", v.Margin())
	}
}

func TestVehicle_SetLocation(t *testing.T) {
	v := validVehicle()
	v.SetLocation(42)
	if v.LocationID != 42 {
		t.Errorf("LocationID = %v, want 42", v.LocationID)
	}
}