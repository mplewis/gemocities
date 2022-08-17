package main

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
)

type Shutdownable interface {
	Shutdown(ctx context.Context) error
}

func gracefullyShutdown(ctx context.Context, name string, srv Shutdownable, wg *sync.WaitGroup) {
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Str("server", name).Msg("Server failed to shutdown")
	} else {
		log.Info().Str("server", name).Msg("Server shutdown complete")
	}
	wg.Done()
}

func gracefullyShutdownAll(servers map[string]Shutdownable) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Info().Msg("Shutting down")
	wg := &sync.WaitGroup{}
	wg.Add(len(servers))
	for name, server := range servers {
		name, server := name, server
		go func() {
			gracefullyShutdown(ctx, name, server, wg)
		}()
	}
	wg.Wait()
	log.Info().Msg("Shutdown complete")
}
