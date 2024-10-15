package main

type PowerInfo struct {
	CircuitID           int     `json:"CircuitID"`
	PowerProduction     float64 `json:"PowerProduction"`
	PowerConsumed       float64 `json:"PowerConsumed"`
	PowerCapacity       float64 `json:"PowerCapacity"`
	PowerMaxConsumed    float64 `json:"PowerMaxConsumed"`
	BatteryInput        float64 `json:"BatteryInput"`
	BatteryOutput       float64 `json:"BatteryOutput"`
	BatteryDifferential float64 `json:"BatteryDifferential"`
	BatteryPercent      float64 `json:"BatteryPercent"`
	BatteryCapacity     float64 `json:"BatteryCapacity"`
	BatteryTimeEmpty    string  `json:"BatteryTimeEmpty"`
	BatteryTimeFull     string  `json:"BatteryTimeFull"`
	FuseTriggered       bool    `json:"FuseTriggered"`
}
