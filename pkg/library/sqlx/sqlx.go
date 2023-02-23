package sqlx

import (
	"github.com/yenole/chainx/pkg/model"
	"gorm.io/gorm"
)

type wrap struct {
	db *gorm.DB
}

func Wrap(db *gorm.DB) *wrap {
	return &wrap{db: db}
}

func (w *wrap) Chains() ([]*model.Chain, error) {
	var list []*model.Chain
	return list, w.db.Find(&list).Error
}
