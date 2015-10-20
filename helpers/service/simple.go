package service_helpers

import (
	"os"
	"errors"
	"syscall"
	"os/signal"
	service "github.com/ayufan/golang-kardianos-service"
)

var (
	ErrNotSupported = errors.New("Not supported.")
)

type SimpleService struct {
	i service.Interface
	c *service.Config
}

func (s *SimpleService) Run() (err error) {
	err = s.i.Start(s)
	if err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 3)
	signal.Notify(sigChan, syscall.SIGTERM, os.Interrupt)

	<-sigChan

	return s.i.Stop(s)
}

func (s *SimpleService) Start() error {
	return service.ErrNoServiceSystemDetected
}

func (s *SimpleService) Stop() error {
	return ErrNotSupported
}

func (s *SimpleService) Restart() error {
	return ErrNotSupported
}

func (s *SimpleService) Install() error {
	return ErrNotSupported
}

func (s *SimpleService) Uninstall() error {
	return ErrNotSupported
}

func (s *SimpleService) Logger(errs chan<- error) (service.Logger, error) {
	return service.ConsoleLogger, nil
}

func (s *SimpleService) SystemLogger(errs chan<- error) (service.Logger, error) {
	return nil, ErrNotSupported
}

func (s *SimpleService) String() string {
	return "SimpleService"
}
