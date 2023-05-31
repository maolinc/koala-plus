package svc

import (
	"koala/gateway/internal/config"
	"koala/gateway/internal/tools/casbin"
	"koala/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ServiceContext struct {
	Config          config.Config
	CBS             *casbin.CBS
	UserModel       model.SysUserModel
	SysMenuModel    model.SysMenuModel
	SysApiModel     model.SysApiModel
	SysRoleModel    model.SysRoleModel
	SysPostModel    model.SysPostModel
	SysDeptModel    model.SysDeptModel
	AppModel        model.SysApplicationModel
	PermissionModel model.SysPermissionModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, _ := gorm.Open(mysql.Open(c.DB.DataSource), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	return &ServiceContext{
		Config:          c,
		CBS:             casbin.NewCBS(db),
		UserModel:       model.NewSysUserModel(db),
		SysMenuModel:    model.NewSysMenuModel(db),
		SysApiModel:     model.NewSysApiModel(db),
		SysRoleModel:    model.NewSysRoleModel(db),
		SysPostModel:    model.NewSysPostModel(db),
		SysDeptModel:    model.NewSysDeptModel(db),
		AppModel:        model.NewSysApplicationModel(db),
		PermissionModel: model.NewSysPermissionModel(db),
	}
}
