package model

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

var (
	_ SysMenuModel = (*defaultSysMenuModel)(nil)
)

type (
	SysMenu struct {
		Id                int64      `gorm:"id;primary_key" primaryKey:"yes"`   // 菜单ID
		MenuName          string     `gorm:"menu_name"`                         // 菜单名称
		ParentId          int64      `gorm:"parent_id" fid:"Id"`                // 父菜单ID
		OrderNo           int64      `gorm:"order_no"`                          // 显示顺序
		Path              string     `gorm:"path"`                              // 路由地址
		Component         string     `gorm:"component"`                         // 组件路径
		Redirect          string     `gorm:"redirect"`                          // 重定向
		Query             string     `gorm:"query"`                             // 路由参数
		IsFrame           int64      `gorm:"is_frame"`                          // 是否为外链（0是 1否）
		IsCache           int64      `gorm:"is_cache"`                          // 是否缓存（0缓存 1不缓存）
		MenuType          string     `gorm:"menu_type"`                         // 菜单类型（M目录 C菜单 F按钮）
		Visible           string     `gorm:"visible"`                           // 菜单状态（0显示 1隐藏）
		Status            string     `gorm:"status"`                            // 菜单状态（0正常 1停用）
		Perms             string     `gorm:"perms"`                             // 权限标识
		Icon              string     `gorm:"icon"`                              // 菜单图标
		CreateBy          int64      `gorm:"create_by"`                         // 创建者
		CreateTime        *time.Time `gorm:"column:create_time;autoCreateTime"` // 创建时间
		UpdateBy          int64      `gorm:"update_by"`                         // 更新者
		UpdateTime        *time.Time `gorm:"column:update_time;autoUpdateTime"` // 更新时间
		Remark            string     `gorm:"remark"`                            // 备注
		CurrentActiveMenu string     `gorm:"currentActiveMenu"`                 // 当前激活的菜单
		ConfigJson        string     `gorm:"config_json"`
		AppId             int64      `gorm:"app_id"` //应用id
		Children          []SysMenu  `gorm:"-" json:"children"`
	}

	// SysMenuQuery query po
	SysMenuQuery struct {
		SearchBase
		*SysMenu
	}

	SysMenuModel interface {
		// Trans Transaction
		Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error)
		// Builder Custom assembly conditions
		Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB)
		Insert(ctx context.Context, data *SysMenu, db ...*gorm.DB) (err error)
		Update(ctx context.Context, data *SysMenu, db ...*gorm.DB) (err error)
		// Delete When a tombstone field is set, it is tombstone. Otherwise, it is physically deleted
		Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		// ForceDelete Physical deletion
		ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		Count(ctx context.Context, cond *SysMenuQuery) (total int64, err error)
		FindOne(ctx context.Context, id int64) (data *SysMenu, err error)
		// FindByPage Contains total information
		FindByPage(ctx context.Context, cond *SysMenuQuery) (total int64, list []*SysMenu, err error)
		// FindListByPage Normal pagination
		FindListByPage(ctx context.Context, cond *SysMenuQuery) (list []*SysMenu, err error)
		// FindListByCursor Cursor is required based on cursor paging, Only the primary key is of type int, and other types can be expanded by themselves
		FindListByCursor(ctx context.Context, cond *SysMenuQuery) (list []*SysMenu, err error)
		FindAll(ctx context.Context, cond *SysMenuQuery) (list []SysMenu, err error)
		// ---------------Write your other interfaces below---------------

		FindByPerms(ctx context.Context, perms string) (*SysMenu, error)
		FindByIds(ctx context.Context, ids []int64) ([]SysMenu, error)
		FindByPermsList(ctx context.Context, perms []string) ([]SysMenu, error)
	}

	defaultSysMenuModel struct {
		*customConn
		table string
	}
)

func NewSysMenuModel(db *gorm.DB) SysMenuModel {
	return &defaultSysMenuModel{
		customConn: newCustomConnNoCache(db),
		table:      "sys_menu",
	}
}

func (m *SysMenu) TableName() string {
	return "`sys_menu`"
}

func (m *defaultSysMenuModel) conn(ctx context.Context, db ...*gorm.DB) *gorm.DB {
	if len(db) == 0 {
		return m.db.Model(&SysMenu{}).Session(&gorm.Session{Context: ctx})
	}
	return db[0]
}

func (m *defaultSysMenuModel) Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB) {
	return m.conn(ctx, db...)
}

func (m *defaultSysMenuModel) Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error) {
	return m.conn(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}

func (m *defaultSysMenuModel) Insert(ctx context.Context, data *SysMenu, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Create(data).Error
	})
}

func (m *defaultSysMenuModel) Update(ctx context.Context, data *SysMenu, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Model(&data).Updates(data).Error
	})
}

func (m *defaultSysMenuModel) Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Where(" id = ?", id).Delete(SysMenu{}).Error
	})
}

func (m *defaultSysMenuModel) ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Unscoped().Where(" id = ?", id).Delete(SysMenu{}).Error
	})
}

func (m *defaultSysMenuModel) Count(ctx context.Context, cond *SysMenuQuery) (total int64, err error) {
	err = m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.SysMenu).Count(&total).Error
	return total, err
}

func (m *defaultSysMenuModel) FindOne(ctx context.Context, id int64) (data *SysMenu, err error) {
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

func (m *defaultSysMenuModel) FindByPage(ctx context.Context, cond *SysMenuQuery) (total int64, list []*SysMenu, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.SysMenu)

	conn.Count(&total)
	err = conn.Scopes(
		pageScope(cond.PageCurrent, cond.PageSize),
		orderScope(cond.OrderSort...),
	).Find(&list).Error
	return total, list, err
}

func (m *defaultSysMenuModel) FindListByPage(ctx context.Context, cond *SysMenuQuery) (list []*SysMenu, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		orderScope(cond.OrderSort...),
		pageScope(cond.PageCurrent, cond.PageSize),
	)
	err = conn.Where(cond.SysMenu).Find(&list).Error
	return list, err
}

func (m *defaultSysMenuModel) FindListByCursor(ctx context.Context, cond *SysMenuQuery) (list []*SysMenu, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		cursorScope(cond.Cursor, cond.CursorAsc, cond.PageSize, "id"),
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.SysMenu).Find(&list).Error
	return list, err
}

func (m *defaultSysMenuModel) FindAll(ctx context.Context, cond *SysMenuQuery) (list []SysMenu, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.SysMenu).Find(&list).Error
	return list, err
}

func (m *defaultSysMenuModel) FindByPerms(ctx context.Context, perms string) (data *SysMenu, err error) {
	tx := m.conn(ctx).Where(&SysMenu{Perms: perms}).Find(&data)
	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}
	return data, tx.Error
}

func (m *defaultSysMenuModel) FindByIds(ctx context.Context, ids []int64) ([]SysMenu, error) {
	var ts []SysMenu
	err := m.conn(ctx).Where("id in ?", ids).Order("parent_id asc, order_no asc").Find(&ts).Error
	return ts, err
}

func (m *defaultSysMenuModel) FindByPermsList(ctx context.Context, perms []string) ([]SysMenu, error) {
	var ts []SysMenu
	err := m.conn(ctx).Where("perms in ?", perms).Find(&ts).Error
	return ts, err
}
