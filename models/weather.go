package models

type Weather struct {
	Temp     float64 `json:"temp"`
	Humidity float64 `json:"humidity"`
	Wind     float64 `json:"wind"`
	RainProb float64 `json:"rainProb"`
}
