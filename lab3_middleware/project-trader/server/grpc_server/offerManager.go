package grpcserver

import (
	"math"
	"math/rand/v2"
	"time"

	"example.com/trading-app/trader"
)

var CurrentOffer []*trader.Offer

func RunSimulation() {
	if len(CurrentOffer) == 0 {
		instruments := initInstruments()
		for _, instrument := range instruments {
			CurrentOffer = append(CurrentOffer, &trader.Offer{
				Symbol:       instrument.Symbol,
				CurrentPrice: instrument.Price,
				Type:         trader.InstrumentType(trader.InstrumentType_value[instrument.Type]),
				BaseCurrency: trader.Currency(trader.Currency_value[instrument.BaseCurrency]),
			})
		}
	}
	for {
		for _, offer := range CurrentOffer {
			// Simulate price changes
			offer.CurrentPrice += (rand.Float64() - 0.5) * offer.CurrentPrice * 0.05 // Random price change
			offer.CurrentPrice = math.Max(
				0,
				offer.CurrentPrice,
			)
		}
		time.Sleep(
			time.Duration(rand.IntN(5)) * time.Second,
		) // Simulate a delay between price updates
		// fmt.Println("Current offers:")
		// for _, offer := range CurrentOffer {
		// 	fmt.Printf("Symbol: %s, Current Price: %.2f\n", offer.Symbol, offer.CurrentPrice)
		// }
		// println("------------------------------------------------")
	}
}
func FetCurrentOffer(symbols ...string) []*trader.Offer {
	if len(symbols) == 0 {
		return CurrentOffer
	}
	var filteredOffers []*trader.Offer
	for _, symbol := range symbols {
		for _, offer := range CurrentOffer {
			if offer.Symbol == symbol {
				filteredOffers = append(filteredOffers, offer)
				break
			}
		}
	}
	return filteredOffers
}
