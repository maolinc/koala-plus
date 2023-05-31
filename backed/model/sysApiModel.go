package model

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

var (
	_ SysApiModel = (*defaultSysApiModel)(nil)
)

type (
	SysApi struct {
		Id         int64      `gorm:"id;primary_key"`                    //菜单ID
		Des        string     `gorm:"des"`                               // 描述
		Path       string     `gorm:"path"`                              // 路径
		Method     string     `gorm:"method"`                            // 路径
		Group      string     `gorm:"group"`                             // 分组
		Status     string     `gorm:"status"`                            // 菜单状态（0正常 1停用）
		CreateBy   int64      `gorm:"create_by"`                         // 创建者
		CreateTime *time.Time `gorm:"column:create_time;autoCreateTime"` // 创建时间
		UpdateBy   int64      `gorm:"update_by"`                         // 更新者
		UpdateTime *time.Time `gorm:"column:update_time;autoUpdateTime"` // 更新时间
		Remark     string     `gorm:"remark"`
		Perms      string     `gorm:"perms"` // 角色权限字符串
		Name       string     `gorm:"name"`
		AppId      int64      `gorm:"app_id"` //应用id
	}

	// SysApiQuery query po
	SysApiQuery struct {
		SearchBase
		*SysApi
	}

	SysApiModel interface {
		// Trans Transaction
		Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error)
		// Builder Custom assembly conditions
		Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB)
		Insert(ctx context.Context, data *SysApi, db ...*gorm.DB) (err error)
		Update(ctx context.Context, data *SysApi, db ...*gorm.DB) (err error)
		// Delete When a tombstone field is set, it is tombstone. Otherwise, it is physically deleted
		Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		// ForceDelete Physical deletion
		ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		Count(ctx context.Context, cond *SysApiQuery) (total int64, err error)
		FindOne(ctx context.Context, id int64) (data *SysApi, err error)
		// FindByPage Contains total information
		FindByPage(ctx context.Context, cond *SysApiQuery) (total int64, list []*SysApi, err error)
		// FindListByPage Normal pagination
		FindListByPage(ctx context.Context, cond *SysApiQuery) (list []*SysApi, err error)
		// FindListByCursor Cursor is required based on cursor paging, Only the primary key is of type int, and other types can be expanded by themselves
		FindListByCursor(ctx context.Context, cond *SysApiQuery) (list []*SysApi, err error)
		FindAll(ctx context.Context, cond *SysApiQuery) (list []*SysApi, err error)
		// ---------------Write your other interfaces below---------------

		FindByPath(ctx context.Context, path string, method string) (*SysApi, error)
		FindByPerms(ctx context.Context, perms string) (data *SysApi, err error)
	}

	defaultSysApiModel struct {
		*customConn
		table string
	}
)

func NewSysApiModel(db *gorm.DB) SysApiModel {
	return &defaultSysApiModel{
		customConn: newCustomConnNoCache(db),
		table:      "sys_api",
	}
}

func (m *SysApi) TableName() string {
	return "`sys_api`"
}

func (m *defaultSysApiModel) conn(ctx context.Context, db ...*gorm.DB) *gorm.DB {
	if len(db) == 0 {
		return m.db.Model(&SysApi{}).Session(&gorm.Session{Context: ctx})
	}
	return db[0]
}

func (m *defaultSysApiModel) Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB) {
	return m.conn(ctx, db...)
}

func (m *defaultSysApiModel) Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error) {
	return m.conn(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}

func (m *defaultSysApiModel) Insert(ctx context.Context, data *SysApi, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Create(data).Error
	})
}

func (m *defaultSysApiModel) Update(ctx context.Context, data *SysApi, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Model(&data).Updates(data).Error
	})
}

func (m *defaultSysApiModel) Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Where(" id = ?", id).Delete(SysApi{}).Error
	})
}

func (m *defaultSysApiModel) ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Unscoped().Where(" id = ?", id).Delete(SysApi{}).Error
	})
}

func (m *defaultSysApiModel) Count(ctx context.Context, cond *SysApiQuery) (total int64, err error) {
	err = m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.SysApi).Count(&total).Error
	return total, err
}

func (m *defaultSysApiModel) FindOne(ctx context.Context, id int64) (data *SysApi, err error) {
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

func (m *defaultSysApiModel) FindByPage(ctx context.Context, cond *SysApiQuery) (total int64, list []*SysApi, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.SysApi)

	conn.Count(&total)
	err = conn.Scopes(
		pageScope(cond.PageCurrent, cond.PageSize),
		orderScope(cond.OrderSort...),
	).Find(&list).Error
	return total, list, err
}

func (m *defaultSysApiModel) FindListByPage(ctx context.Context, cond *SysApiQuery) (list []*SysApi, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		orderScope(cond.OrderSort...),
		pageScope(cond.PageCurrent, cond.PageSize),
	)
	err = conn.Where(cond.SysApi).Find(&list).Error
	return list, err
}

func (m *defaultSysApiModel) FindListByCursor(ctx context.Context, cond *SysApiQuery) (list []*SysApi, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		cursorScope(cond.Cursor, cond.CursorAsc, cond.PageSize, "id"),
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.SysApi).Find(&list).Error
	return list, err
}

func (m *defaultSysApiModel) FindAll(ctx context.Context, cond *SysApiQuery) (list []*SysApi, err error) {
	conn := m.conn(ctx).Scopes(
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.SysApi).Find(&list).Error
	return list, err
}

func (m *defaultSysApiModel) FindByPath(ctx context.Context, path string, method string) (data *SysApi, err error) {
	tx := m.conn(ctx).Where("path = ? method = ?", path, method).Find(&data)
	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}
	return data, tx.Error
}

func (m *defaultSysApiModel) FindByPerms(ctx context.Context, perms string) (data *SysApi, err error) {
	tx := m.conn(ctx).Where("perms = ?", perms).Find(&data)
	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}
	return data, tx.Error
}
