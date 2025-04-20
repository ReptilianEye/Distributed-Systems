package grpcserver

import (
	"context"
	"fmt"
	"time"

	"example.com/trading-app/trader"
	"google.golang.org/grpc"
)

type server struct {
	trader.UnimplementedTraderServer
}

func New() *grpc.Server {
	s := grpc.NewServer()
	trader.RegisterTraderServer(s, &server{})
	return s
}

func (s *server) GetOffers(ctx context.Context, filter *trader.Filter) (*trader.OfferList, error) {
	// Implement the logic to get offers based on the filter
	return &trader.OfferList{
		Offers: FetCurrentOffer(),
	}, nil
}

func (s *server) Subscribe(
	in *trader.Subscription, stream trader.Trader_SubscribeServer,
) error {
	symbols := in.Symbols
	previousRates := make(map[string]float64)

	for {
		offers := FetCurrentOffer(symbols...)
		updatedOffers := make([]*trader.Offer, 0)
		for _, offer := range offers {
			if previousRate, ok := previousRates[offer.Symbol]; !ok ||
				offer.CurrentPrice != previousRate {
				updatedOffers = append(updatedOffers, offer)
				previousRates[offer.Symbol] = offer.CurrentPrice
			}
		}
		fmt.Println(offers, previousRates)
		stream.Send(&trader.OfferList{Offers: updatedOffers})
		time.Sleep(time.Second)
	}
}
