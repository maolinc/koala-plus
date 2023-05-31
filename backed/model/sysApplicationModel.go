package model

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"time"
)

var (
	_ SysApplicationModel = (*defaultSysApplicationModel)(nil)
)

type (
	SysApplication struct {
		Id               int64                 `gorm:"id;primary_key"`                     //
		Name             string                `gorm:"name"`                               //应用名字
		Des              string                `gorm:"des"`                                //应用描述
		Perms            string                `gorm:"perms"`                              //权限标识
		OrganizationId   int64                 `gorm:"organization_id"`                    //组织id
		Status           string                `gorm:"status"`                             //
		CreateTime       *time.Time            `gorm:"column:create_time;autoCreateTime"`  // 创建时间
		UpdateTime       *time.Time            `gorm:"column:update_time;autoUpdateTime"`  // 更新时间
		DeleteFlag       soft_delete.DeletedAt `gorm:"column:delete_flag;softDelete:flag"` // 删除标志（0代表存在 2代表删除）
		CreateBy         int64                 `gorm:"create_by"`                          //
		CoverUrl         string                `gorm:"cover_url"`                          //封面
		Href             string                `gorm:"href"`                               //连接地址
		LoginCallbackUrl string                `gorm:"login_callback_url"`                 //连接地址
		AppId            string                `gorm:"app_id"`                             //连接地址
	}

	// SysApplication query cond
	SysApplicationQuery struct {
		SearchBase
		SysApplication
	}

	SysApplicationModel interface {
		// Trans Transaction
		Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error)
		// Builder Custom assembly conditions
		Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB)
		Insert(ctx context.Context, data *SysApplication, db ...*gorm.DB) (err error)
		Update(ctx context.Context, data *SysApplication, db ...*gorm.DB) (err error)
		// Delete When a tombstone field is set, it is tombstone. Otherwise, it is physically deleted
		Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		// ForceDelete Physical deletion
		ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		Count(ctx context.Context, cond *SysApplicationQuery) (total int64, err error)
		FindOne(ctx context.Context, id int64) (data *SysApplication, err error)
		// FindByPage Contains total information
		FindByPage(ctx context.Context, cond *SysApplicationQuery) (total int64, list []*SysApplication, err error)
		// FindListByPage Normal pagination
		FindListByPage(ctx context.Context, cond *SysApplicationQuery) (list []*SysApplication, err error)
		// FindListByCursor Cursor is required based on cursor paging, Only the primary key is of type int, and other types can be expanded by themselves
		FindListByCursor(ctx context.Context, cond *SysApplicationQuery) (list []*SysApplication, err error)
		FindAll(ctx context.Context, cond *SysApplicationQuery) (list []*SysApplication, err error)
		// ---------------Write your other interfaces below---------------

		FindByPerms(ctx context.Context, perms string) (data *SysApplication, err error)
		FindByAppId(ctx context.Context, appId string) (data *SysApplication, err error)
	}

	defaultSysApplicationModel struct {
		*customConn
		table string
	}
)

func NewSysApplicationModel(db *gorm.DB) SysApplicationModel {
	return &defaultSysApplicationModel{
		customConn: newCustomConnNoCache(db),
		table:      "sys_application",
	}
}

func (m *SysApplication) TableName() string {
	return "`sys_application`"
}

func (m *defaultSysApplicationModel) conn(ctx context.Context, db ...*gorm.DB) *gorm.DB {
	if len(db) == 0 {
		return m.db.Model(&SysApplication{}).Session(&gorm.Session{Context: ctx})
	}
	return db[0]
}

func (m *defaultSysApplicationModel) Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB) {
	return m.conn(ctx, db...)
}

func (m *defaultSysApplicationModel) Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error) {
	return m.conn(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}

func (m *defaultSysApplicationModel) Insert(ctx context.Context, data *SysApplication, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Create(data).Error
	})
}

func (m *defaultSysApplicationModel) Update(ctx context.Context, data *SysApplication, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Model(&data).Updates(data).Error
	})
}

func (m *defaultSysApplicationModel) Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Where(" id = ?", id).Updates(&SysApplication{DeleteFlag: 1}).Error
	})
}

func (m *defaultSysApplicationModel) ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Unscoped().Where(" id = ?", id).Delete(SysApplication{}).Error
	})
}

func (m *defaultSysApplicationModel) Count(ctx context.Context, cond *SysApplicationQuery) (total int64, err error) {
	err = m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.SysApplication).Count(&total).Error
	return total, err
}

func (m *defaultSysApplicationModel) FindOne(ctx context.Context, id int64) (data *SysApplication, err error) {
	err = m.QueryRow(ctx, &data, func(v interface{}) error {
		tx := m.conn(ctx).Where(" id = ?", id).Find(v)
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

func (m *defaultSysApplicationModel) FindByPage(ctx context.Context, cond *SysApplicationQuery) (total int64, list []*SysApplication, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.SysApplication)

	total, list, err = pageHandler[*SysApplication](conn, cond.PageCurrent, cond.PageSize)
	return total, list, err
}

func (m *defaultSysApplicationModel) FindListByPage(ctx context.Context, cond *SysApplicationQuery) (list []*SysApplication, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		orderScope(cond.OrderSort...),
		pageScope(cond.PageCurrent, cond.PageSize),
	)
	err = conn.Where(cond.SysApplication).Find(&list).Error
	return list, err
}

func (m *defaultSysApplicationModel) FindListByCursor(ctx context.Context, cond *SysApplicationQuery) (list []*SysApplication, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		cursorScope(cond.Cursor, cond.CursorAsc, cond.PageSize, "id"),
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.SysApplication).Find(&list).Error
	return list, err
}

func (m *defaultSysApplicationModel) FindAll(ctx context.Context, cond *SysApplicationQuery) (list []*SysApplication, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.SysApplication).Find(&list).Error
	return list, err
}

func (m *defaultSysApplicationModel) FindByPerms(ctx context.Context, perms string) (data *SysApplication, err error) {
	tx := m.conn(ctx).Where(&SysApplication{Perms: perms}).Find(&data)
	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}
	return data, tx.Error
}

func (m *defaultSysApplicationModel) FindByAppId(ctx context.Context, appId string) (data *SysApplication, err error) {
	tx := m.conn(ctx).Where("app_id = ?", appId).Find(&data)
	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}
	return data, tx.Error
}
