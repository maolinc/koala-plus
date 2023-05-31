package model

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

var (
	_ SysDeptModel = (*defaultSysDeptModel)(nil)
)

type (
	SysDept struct {
		Id         int64                 `gorm:"id;primary_key" primaryKey:"yes"`    //部门id
		ParentId   int64                 `gorm:"parent_id" fid:"Id"`                 // 父部门id
		Ancestors  string                `gorm:"ancestors"`                          // 祖级列表
		DeptName   string                `gorm:"dept_name"`                          // 部门名称
		OrderNum   int64                 `gorm:"order_num"`                          // 显示顺序
		Leader     string                `gorm:"leader"`                             // 负责人
		Phone      string                `gorm:"phone"`                              // 联系电话
		Email      string                `gorm:"email"`                              // 邮箱
		Status     string                `gorm:"status"`                             // 部门状态（0正常 1停用）
		DeleteFlag soft_delete.DeletedAt `gorm:"column:delete_flag;softDelete:flag"` // 删除标志（0代表存在 2代表删除）
		CreateBy   int64                 `gorm:"create_by"`                          // 创建者
		CreateTime *time.Time            `gorm:"column:create_time;autoCreateTime"`  // 创建时间
		UpdateBy   int64                 `gorm:"update_by"`                          // 更新者
		UpdateTime *time.Time            `gorm:"column:update_time;autoUpdateTime"`  // 更新时间
		Perms      string                `gorm:"perms"`                              // 角色权限字符串
		AppId      int64                 `gorm:"app_id"`                             //应用id
		Children   []SysDept             `gorm:"-" json:"children"`                  // 备注
	}

	// SysDeptQuery query po
	SysDeptQuery struct {
		SearchBase
		*SysDept
		OpenTree  bool
		ParentIds []int64
	}

	SysDeptModel interface {
		// Trans Transaction
		Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error)
		// Builder Custom assembly conditions
		Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB)
		Insert(ctx context.Context, data *SysDept, db ...*gorm.DB) (err error)
		Update(ctx context.Context, data *SysDept, db ...*gorm.DB) (err error)
		// Delete When a tombstone field is set, it is tombstone. Otherwise, it is physically deleted
		Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		// ForceDelete Physical deletion
		ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		Count(ctx context.Context, cond *SysDeptQuery) (total int64, err error)
		FindOne(ctx context.Context, id int64) (data *SysDept, err error)
		// FindByPage Contains total information
		FindByPage(ctx context.Context, cond *SysDeptQuery) (total int64, list []*SysDept, err error)
		// FindListByPage Normal pagination
		FindListByPage(ctx context.Context, cond *SysDeptQuery) (list []*SysDept, err error)
		// FindListByCursor Cursor is required based on cursor paging, Only the primary key is of type int, and other types can be expanded by themselves
		FindListByCursor(ctx context.Context, cond *SysDeptQuery) (list []*SysDept, err error)
		FindAll(ctx context.Context, cond *SysDeptQuery) (list []SysDept, err error)
		// ---------------Write your other interfaces below---------------

		FindByParentIds(ctx context.Context, pIds []int64) ([]SysDept, error)
		FindByPerms(ctx context.Context, perms string) (data *SysDept, err error)
	}

	defaultSysDeptModel struct {
		*customConn
		table string
	}
)

func NewSysDeptModel(db *gorm.DB) SysDeptModel {
	return &defaultSysDeptModel{
		customConn: newCustomConnNoCache(db),
		table:      "sys_dept",
	}
}

func (m *SysDept) TableName() string {
	return "`sys_dept`"
}

func (m *defaultSysDeptModel) conn(ctx context.Context, db ...*gorm.DB) *gorm.DB {
	if len(db) == 0 {
		return m.db.Model(&SysDept{}).Session(&gorm.Session{Context: ctx})
	}
	return db[0]
}

func (m *defaultSysDeptModel) Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB) {
	return m.conn(ctx, db...)
}

func (m *defaultSysDeptModel) Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error) {
	return m.conn(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}

func (m *defaultSysDeptModel) Insert(ctx context.Context, data *SysDept, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Create(data).Error
	})
}

func (m *defaultSysDeptModel) Update(ctx context.Context, data *SysDept, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Model(&data).Updates(data).Error
	})
}

func (m *defaultSysDeptModel) Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Where(" id = ?", id).Updates(&SysDept{DeleteFlag: 1}).Error
	})
}

func (m *defaultSysDeptModel) ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Unscoped().Where(" id = ?", id).Delete(SysDept{}).Error
	})
}

func (m *defaultSysDeptModel) Count(ctx context.Context, cond *SysDeptQuery) (total int64, err error) {
	err = m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.SysDept).Count(&total).Error
	return total, err
}

func (m *defaultSysDeptModel) FindOne(ctx context.Context, id int64) (data *SysDept, err error) {
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

func (m *defaultSysDeptModel) FindByPage(ctx context.Context, cond *SysDeptQuery) (total int64, list []*SysDept, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.SysDept)

	conn.Count(&total)
	err = conn.Scopes(
		pageScope(cond.PageCurrent, cond.PageSize),
		orderScope(cond.OrderSort...),
	).Find(&list).Error
	return total, list, err
}

func (m *defaultSysDeptModel) FindListByPage(ctx context.Context, cond *SysDeptQuery) (list []*SysDept, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		orderScope(cond.OrderSort...),
		pageScope(cond.PageCurrent, cond.PageSize),
	)
	err = conn.Where(cond.SysDept).Find(&list).Error
	return list, err
}

func (m *defaultSysDeptModel) FindListByCursor(ctx context.Context, cond *SysDeptQuery) (list []*SysDept, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		cursorScope(cond.Cursor, cond.CursorAsc, cond.PageSize, "id"),
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.SysDept).Find(&list).Error
	return list, err
}

func (m *defaultSysDeptModel) FindAll(ctx context.Context, cond *SysDeptQuery) (list []SysDept, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.SysDept).Find(&list).Error
	return list, err
}

func (m *defaultSysDeptModel) FindByParentIds(ctx context.Context, pIds []int64) ([]SysDept, error) {
	var ts []SysDept
	err := m.conn(ctx).Where("parent_id in ?", pIds).Find(&ts).Error
	return ts, err
}

func (m *defaultSysDeptModel) FindByPerms(ctx context.Context, perms string) (data *SysDept, err error) {
	tx := m.conn(ctx).Where(&SysDept{Perms: perms}).Find(&data)
	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}
	return data, tx.Error
}
