package api

import (
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yenole/chainx/pkg/library/ether"
	"github.com/yenole/chainx/pkg/library/sqlx"
	"github.com/yenole/chainx/pkg/model"
)

const (
	cs_ok   uint32 = iota // OK
	cs_fail               // not ok
)

type chain struct {
	ch     *model.Chain
	height uint64
	state  uint32
}

func (c *chain) isValid() bool {
	return atomic.LoadUint32(&c.state) == cs_ok
}

func (c *chain) Valid(logger *logrus.Logger) {
	height, err := ether.BlockNumber(c.ch)
	if err != nil {
		atomic.StoreUint32(&c.state, cs_fail)
		logger.Errorf("chain %v -> %v valid fail:%v", c.ch.CID, c.ch.URL, err)
		return
	}
	c.height = height
	atomic.StoreUint32(&c.state, cs_ok)
}

// 1.定时获取新的chain配置
// 2.定时获取区块新的高度，用于判断是否可用
func (s *Service) runLoop() {
	s.loadChain()

	ticker := time.NewTicker(time.Minute * 10)
	for {
		select {
		case <-ticker.C:
			s.loopTicker()

		}
	}
}

func (s *Service) loadChain() {
	list, err := sqlx.Wrap(s.db).Chains()
	if err != nil {
		s.logger.Errorf("load chain fail:%v", err)
		time.AfterFunc(time.Minute, s.loadChain)
		return
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	for _, c := range list {
		s.logger.Infof("load chain %v -> %v", c.CID, c.URL)
		ch := &chain{ch: c}
		s.chs[c.CID] = append(s.chs[c.CID], ch)
		go ch.Valid(s.logger)
	}
}

func (s *Service) loopTicker() {
	s.mux.RLock()
	defer s.mux.RUnlock()
	for id, list := range s.chs {
		s.logger.Infof("ticker valid %v-%v rpcs", id, len(list))
		for _, c := range list {
			go c.Valid(s.logger)
		}
	}
}

func (s *Service) safeChain(id uint) *model.Chain {
	s.mux.RLock()
	defer s.mux.RUnlock()
	if list, ok := s.chs[id]; ok && len(list) > 0 {
		valids := make([]*chain, 0, len(list))
		for _, c := range list {
			if c.isValid() {
				valids = append(valids, c)
			}
		}
		return valids[rand.Intn(len(valids))].ch
	}
	return nil
}
