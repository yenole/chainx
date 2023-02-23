package api

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yenole/chainx/pkg/config"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type Service struct {
	db *gorm.DB
	rt *gin.Engine

	mux    sync.RWMutex
	chs    map[uint][]*chain
	logger *logrus.Logger
}

func New(db *gorm.DB, logger *logrus.Logger) *Service {
	rt := gin.New()
	rt.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(params gin.LogFormatterParams) string {
			if params.Latency > time.Minute {
				params.Latency = params.Latency - params.Latency%time.Second
			}
			logger.Infof("[GIN]:| %03d | %13v | %15s | %-7s  %v",
				params.StatusCode,
				params.Latency,
				params.ClientIP,
				params.Method,
				params.Path,
			)
			if params.ErrorMessage != "" {
				logger.Errorln(params.ErrorMessage)
			}
			return ""
		},
		Output: logger.Out,
	}), gin.Recovery())
	s := &Service{db: db, rt: rt, chs: make(map[uint][]*chain), logger: logger}
	s.setupRouting()
	go s.runLoop()
	return s
}

func (s *Service) Run() error {
	s.logger.Infof("api listen:%v", *config.Listen)
	go func() {
		err := s.rt.Run(*config.Listen)
		if err != nil {
			s.logger.Errorf("api listen:%v fail:%v", *config.Listen, err)
		}
	}()
	return nil
}
