package grpcserver

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Instrument struct {
	Symbol       string  `json:"symbol"`
	Price        float64 `json:"price"`
	Type         string  `json:"type"`
	BaseCurrency string  `json:"base_currency"`
}

func initInstruments() []*Instrument {
	var baseInstruments []*Instrument
	byteValue, err := os.ReadFile("instruments.json")
	if err != nil {
		log.Fatalf("failed to read instruments file: %v", err)
	}

	json.Unmarshal(byteValue, &baseInstruments)
	fmt.Println("Instruments loaded from JSON file:")
	for _, instrument := range baseInstruments {
		fmt.Printf(
			"Symbol: %s, Price: %.2f, Type: %s, BaseCurrency: %s\n",
			instrument.Symbol,
			instrument.Price,
			instrument.Type,
			instrument.BaseCurrency,
		)
	}
	return baseInstruments
}
