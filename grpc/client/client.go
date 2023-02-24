package client

import (
	"websocket-service/config"
)

type ServiceManagerI interface {
	// // Instrument Service
	// InstrumentService() instrument_service.InstrumentSeviceClient
}

type grpcClients struct {
	// instrumentService   instrument_service.InstrumentSeviceClient
}

func NewGrpcClients(cfg config.Config) (ServiceManagerI, error) {
	// connInstrumentService, err := grpc.Dial(
	// 	cfg.InstrumentServiceHost+cfg.InstrumentGRPCPort,
	// 	grpc.WithInsecure(),
	// )
	// if err != nil {
	// 	return nil, err
	// }

	return &grpcClients{
		// instrumentService:   instrument_service.NewInstrumentSeviceClient(connInstrumentService),
	}, nil
}

// func (g *grpcClients) InstrumentService() instrument_service.InstrumentSeviceClient {
// 	return g.instrumentService
// }
