package service

import (
	"fmt"
	"sync"

	"github.com/gen1us2k/log"
	"github.com/maddevsio/cinematica2telegram/conf"
)

// Cinematica is main struct of daemon
// it stores all services that used by
type Cinematica struct {
	config *conf.CinematicaConfig

	services  map[string]Service
	waitGroup sync.WaitGroup

	logger log.Logger
}

// NewCinematica creates and returns new CinematicaInstance
func NewCinematica(config *conf.CinematicaConfig) *Cinematica {
	pb := new(Cinematica)
	pb.config = config
	pb.logger = log.NewLogger("fb2telegram")
	pb.services = make(map[string]Service)
	pb.AddService(&CinematicaService{})
	pb.AddService(&TelegramService{})
	return pb
}

// Start starts all services in separate goroutine
func (pb *Cinematica) Start() error {
	pb.logger.Info("Starting bot service")
	for _, service := range pb.services {
		pb.logger.Infof("Initializing: %s\n", service.Name())
		if err := service.Init(pb); err != nil {
			return fmt.Errorf("initialization of %q finished with error: %v", service.Name(), err)
		}
		pb.waitGroup.Add(1)

		go func(srv Service) {
			defer pb.waitGroup.Done()
			pb.logger.Infof("running %q service\n", srv.Name())
			if err := srv.Run(); err != nil {
				pb.logger.Errorf("error on run %q service, %v", srv.Name(), err)
			}
		}(service)
	}
	return nil
}

// AddService adds service into Cinematica.services map
func (pb *Cinematica) AddService(srv Service) {
	pb.services[srv.Name()] = srv

}

// Config returns current instance of CinematicaConfig
func (pb *Cinematica) Config() conf.CinematicaConfig {
	return *pb.config
}

// Stop stops all services running
func (pb *Cinematica) Stop() {
	pb.logger.Info("Worker is stopping...")
	for _, service := range pb.services {
		service.Stop()
	}
}

// WaitStop blocks main thread and waits when all goroutines will be stopped
func (pb *Cinematica) WaitStop() {
	pb.waitGroup.Wait()
}

func (pb *Cinematica) CinematicaService() *CinematicaService {
	service, ok := pb.services["cinematica_service"]
	if !ok {
		pb.logger.Error("Error getting cinematica_service")
	}
	return service.(*CinematicaService)
}
