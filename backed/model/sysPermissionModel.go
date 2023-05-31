package model

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"time"
)

var (
	_ SysPermissionModel = (*defaultSysPermissionModel)(nil)
)

type (
	SysPermission struct {
		Perms      string     `gorm:"perms"`     // 权限标识
		Name       string     `gorm:"name"`      // 权限名字
		Des        string     `gorm:"des"`       // 描述
		Status     string     `gorm:"status"`    // 状态
		AppPerms   string     `gorm:"app_perms"` // 所属应用权标识
		Group      string     `gorm:"group"`     // 所属应用权标识
		Table      string     `gorm:"table"`     //存在那张表
		CreateBy   int64      `gorm:"create_by"`
		CreateTime *time.Time `gorm:"column:create_time;autoCreateTime"` // 创建时间
	}

	// SysPermissionQuery query po
	SysPermissionQuery struct {
		SearchBase
		*SysPermission
	}

	SysPermissionModel interface {
		// Trans Transaction
		Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error)
		// Builder Custom assembly conditions
		Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB)
		Insert(ctx context.Context, data *SysPermission, db ...*gorm.DB) (err error)
		Update(ctx context.Context, data *SysPermission, db ...*gorm.DB) (err error)
		// Delete When a tombstone field is set, it is tombstone. Otherwise, it is physically deleted
		Delete(ctx context.Context, perms string, appPerms string, db ...*gorm.DB) (err error)
		// ForceDelete Physical deletion
		ForceDelete(ctx context.Context, perms string, appPerms string, db ...*gorm.DB) (err error)
		Count(ctx context.Context, cond *SysPermissionQuery) (total int64, err error)
		FindOne(ctx context.Context, perms string, appPerms string) (data *SysPermission, err error)
		// FindByPage Contains total information
		FindByPage(ctx context.Context, cond *SysPermissionQuery) (total int64, list []*SysPermission, err error)
		// FindListByPage Normal pagination
		FindListByPage(ctx context.Context, cond *SysPermissionQuery) (list []*SysPermission, err error)
		// FindListByCursor Cursor is required based on cursor paging, Only the primary key is of type int, and other types can be expanded by themselves
		FindListByCursor(ctx context.Context, cond *SysPermissionQuery) (list []*SysPermission, err error)
		FindAll(ctx context.Context, cond *SysPermissionQuery) (list []*SysPermission, err error)
		// ---------------Write your other interfaces below---------------

		FindByPerms(ctx context.Context, perms string) (data *SysPermission, err error)
		FindByPermsList(ctx context.Context, perms []string, appPerms string) (data []*SysPermission, err error)
		UpdateByPerms(ctx context.Context, perms, appPerms string, data *SysPermission, db ...*gorm.DB) (err error)
	}

	defaultSysPermissionModel struct {
		*customConn
		table string
	}
)

func NewSysPermissionModel(db *gorm.DB) SysPermissionModel {
	return &defaultSysPermissionModel{
		customConn: newCustomConnNoCache(db),
		table:      "sys_permission",
	}
}

func (m *SysPermission) TableName() string {
	return "`sys_permission`"
}

func (m *defaultSysPermissionModel) conn(ctx context.Context, db ...*gorm.DB) *gorm.DB {
	if len(db) == 0 {
		return m.db.Model(&SysPermission{}).Session(&gorm.Session{Context: ctx})
	}
	return db[0]
}

func (m *defaultSysPermissionModel) Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB) {
	return m.conn(ctx, db...)
}

func (m *defaultSysPermissionModel) Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error) {
	return m.conn(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}

func (m *defaultSysPermissionModel) Insert(ctx context.Context, data *SysPermission, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Create(data).Error
	})
}

func (m *defaultSysPermissionModel) Update(ctx context.Context, data *SysPermission, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Model(&data).Updates(data).Error
	})
}

func (m *defaultSysPermissionModel) Delete(ctx context.Context, perms string, appPerms string, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Where("perms = ? and app_perms", perms, appPerms).Delete(SysPermission{}).Error
	})
}

func (m *defaultSysPermissionModel) ForceDelete(ctx context.Context, perms string, appPerms string, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Unscoped().Where("perms = ? and app_perms", perms, appPerms).Delete(SysPermission{}).Error
	})
}

func (m *defaultSysPermissionModel) Count(ctx context.Context, cond *SysPermissionQuery) (total int64, err error) {
	err = m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.SysPermission).Count(&total).Error
	return total, err
}

func (m *defaultSysPermissionModel) FindOne(ctx context.Context, perms string, appPerms string) (data *SysPermission, err error) {
	err = m.QueryRow(ctx, &data, func(v interface{}) error {
		tx := m.conn(ctx).Where("perms = ? and app_perms", perms, appPerms).Find(v)
		if tx.RowsAffected == 0 {
			return sql.ErrNoRows
		}
		return tx.Error
	})
	switch err {
	case nil:
		return data, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *defaultSysPermissionModel) FindByPage(ctx context.Context, cond *SysPermissionQuery) (total int64, list []*SysPermission, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.SysPermission)

	conn.Count(&total)
	err = conn.Scopes(
		pageScope(cond.PageCurrent, cond.PageSize),
		orderScope(cond.OrderSort...),
	).Find(&list).Error
	return total, list, err
}

func (m *defaultSysPermissionModel) FindListByPage(ctx context.Context, cond *SysPermissionQuery) (list []*SysPermission, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		orderScope(cond.OrderSort...),
		pageScope(cond.PageCurrent, cond.PageSize),
	)
	err = conn.Where(cond.SysPermission).Find(&list).Error
	return list, err
}

func (m *defaultSysPermissionModel) FindListByCursor(ctx context.Context, cond *SysPermissionQuery) (list []*SysPermission, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.SysPermission).Find(&list).Error
	return list, err
}

func (m *defaultSysPermissionModel) FindAll(ctx context.Context, cond *SysPermissionQuery) (list []*SysPermission, err error) {
	conn := m.conn(ctx).Scopes(
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.SysPermission).Find(&list).Error
	return list, err
}

func (m *defaultSysPermissionModel) FindByPerms(ctx context.Context, perms string) (data *SysPermission, err error) {
	tx := m.conn(ctx).Where("perms = ?", perms).Find(&data)
	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}
	return data, tx.Error
}

func (m *defaultSysPermissionModel) FindByPermsList(ctx context.Context, perms []string, appPerms string) (data []*SysPermission, err error) {
	err = m.conn(ctx).Where("perms in ? AND app_perms = ?", perms, appPerms).Find(&data).Error
	return data, err
}

func (m *defaultSysPermissionModel) UpdateByPerms(ctx context.Context, perms, appPerms string, data *SysPermission, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Model(&data).Where("perms = ? AND app_perms = ?", perms, appPerms).Updates(data).Error
	})
}
