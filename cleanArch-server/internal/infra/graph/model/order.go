package model

type Order struct {
	ID         string  `json:"id"`
	Price      float64 `json:"Price"`
	Tax        float64 `json:"Tax"`
	FinalPrice float64 `json:"FinalPrice"`
}
