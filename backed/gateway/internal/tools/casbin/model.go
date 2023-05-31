package casbin

import "gorm.io/gorm"

type CasbinRule struct {
	Id    int64  `gorm:"id"`
	Ptype string `gorm:"ptype"`
	V0    string `gorm:"v0"`
	V1    string `gorm:"v1"`
	V2    string `gorm:"v2"`
	V3    string `gorm:"v3"`
	V4    string `gorm:"v4"`
	V5    string `gorm:"v5"`
}

func (c *CasbinRule) TableName() string {
	return "`casbin_rule`"
}

type dbClient struct {
	*gorm.DB
}

func newDBClient(db *gorm.DB) *dbClient {
	return &dbClient{db.Model(&CasbinRule{})}
}

func (d *dbClient) conn() *gorm.DB {
	return d.Table("casbin_rule").Session(&gorm.Session{})
}
