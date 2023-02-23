package chainx

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/yenole/chainx/pkg/api"
	"github.com/yenole/chainx/pkg/config"
	"github.com/yenole/chainx/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Service struct {
	db  *gorm.DB
	api *api.Service

	logger *logrus.Logger
}

func New(logger *logrus.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) setup() (err error) {
	s.db, err = gorm.Open(sqlite.Open(config.JoinPath("data.db")))
	if err != nil {
		return err
	}
	s.db.AutoMigrate(new(model.Chain))

	s.api = api.New(s.db, s.logger)
	return nil
}

func (s *Service) Run() {
	s.logger.Infof("chainx is running!")
	err := s.setup()
	if err != nil {
		s.logger.Fatalf("setup fail:%v", err)
	}

	err = s.api.Run()
	if err != nil {
		s.logger.Fatalf("setup api fail:%v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	<-ctx.Done()

	s.logger.Infof("chainx is stoping")
}
