package domain

import "time"

// Processor - структура для представления информации о процессоре
type Processor struct {
	ID                  int       `json:"id" db:"id"`
	Link                string    `json:"link" db:"link"`
	Line                string    `json:"line" db:"line"`
	Model               string    `json:"model" db:"model"`
	Core                string    `json:"core" db:"core"`
	NumCores            int       `json:"num_cores" db:"num_cores"`
	Manufacturer        string    `json:"manufacturer" db:"manufacturer"`
	ManufacturerCode    string    `json:"manufacturer_code" db:"manufacturer_code"`
	Socket              string    `json:"socket" db:"socket"`
	ThermalPower        int       `json:"thermal_power" db:"thermal_power"`
	ManufacturerWebsite string    `json:"manufacturer_website" db:"manufacturer_website"`
	TechnologyProcess   int       `json:"technology_process" db:"technology_process"`
	Technologies        string    `json:"technologies" db:"technologies"`
	DeliveryType        string    `json:"delivery_type" db:"delivery_type"`
	L1CacheVolume       float64   `json:"l1_cache_volume" db:"l1_cache_volume"`
	L2CacheVolume       float64   `json:"l2_cache_volume" db:"l2_cache_volume"`
	L3CacheVolume       float64   `json:"l3_cache_volume" db:"l3_cache_volume"`
	IntegratedGraphics  bool      `json:"integrated_graphics" db:"integrated_graphics"`
	VideoProcessor      string    `json:"video_processor" db:"video_processor"`
	NewPrice            float64   `json:"new_price" db:"new_price"`
	OldPrice            float64   `json:"old_price" db:"old_price"`
	AvailabilityDate    string    `json:"availability_date" db:"availability_date"`
	City                string    `json:"city" db:"city"`
	Street              string    `json:"street" db:"street"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}
