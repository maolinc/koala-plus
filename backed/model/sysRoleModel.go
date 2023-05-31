package model

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

var (
	_ SysRoleModel = (*defaultSysRoleModel)(nil)
)

type (
	SysRole struct {
		Id         int64                 `gorm:"id;primary_key" primaryKey:"yes"`    // 角色ID
		ParentId   int64                 `gorm:"parent_id" fid:"Id"`                 // 父角色
		RoleName   string                `gorm:"role_name"`                          // 角色名称
		Perms      string                `gorm:"perms"`                              // 角色权限字符串
		RoleSort   int64                 `gorm:"role_sort"`                          // 显示顺序
		DataScope  string                `gorm:"data_scope"`                         // 数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）
		Status     string                `gorm:"status"`                             // 角色状态（0正常 1停用）
		DeleteFlag soft_delete.DeletedAt `gorm:"column:delete_flag;softDelete:flag"` // 删除标志（0代表存在 2代表删除）
		CreateBy   int64                 `gorm:"create_by"`                          // 创建者
		CreateTime *time.Time            `gorm:"column:create_time;autoCreateTime"`  // 创建时间
		UpdateBy   int64                 `gorm:"update_by"`                          // 更新者
		UpdateTime *time.Time            `gorm:"column:update_time;autoUpdateTime"`  // 更新时间
		Remark     string                `gorm:"remark"`                             // 备注
		AppId      int64                 `gorm:"app_id"`                             //应用id
		Children   []SysRole             `gorm:"-" json:"children"`
	}

	// query po
	SysRoleQuery struct {
		SearchBase
		*SysRole
		OpenTree  bool
		ParentIds []int64
	}

	SysRoleModel interface {
		// Trans Transaction
		Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error)
		// Builder Custom assembly conditions
		Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB)
		Insert(ctx context.Context, data *SysRole, db ...*gorm.DB) (err error)
		Update(ctx context.Context, data *SysRole, db ...*gorm.DB) (err error)
		// Delete When a tombstone field is set, it is tombstone. Otherwise, it is physically deleted
		Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		// ForceDelete Physical deletion
		ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		Count(ctx context.Context, cond *SysRoleQuery) (total int64, err error)
		FindOne(ctx context.Context, id int64) (data *SysRole, err error)
		// FindByPage Contains total information
		FindByPage(ctx context.Context, cond *SysRoleQuery) (total int64, list []*SysRole, err error)
		// FindListByPage Normal pagination
		FindListByPage(ctx context.Context, cond *SysRoleQuery) (list []*SysRole, err error)
		// FindListByCursor Cursor is required based on cursor paging, Only the primary key is of type int, and other types can be expanded by themselves
		FindListByCursor(ctx context.Context, cond *SysRoleQuery) (list []*SysRole, err error)
		FindAll(ctx context.Context, cond *SysRoleQuery) (list []SysRole, err error)
		// ---------------Write your other interfaces below---------------

		FindByPerms(ctx context.Context, perms string) (*SysRole, error)
		FindByParentIds(ctx context.Context, pIds []int64) ([]SysRole, error)
		FindByIds(ctx context.Context, ids []int64) ([]SysRole, error)
		FindByPermsList(ctx context.Context, perms []string) ([]SysRole, error)
	}

	defaultSysRoleModel struct {
		*customConn
		table string
	}
)

func NewSysRoleModel(db *gorm.DB) SysRoleModel {
	return &defaultSysRoleModel{
		customConn: newCustomConnNoCache(db),
		table:      "sys_role",
	}
}

func (m *SysRole) TableName() string {
	return "`sys_role`"
}

func (m *defaultSysRoleModel) conn(ctx context.Context, db ...*gorm.DB) *gorm.DB {
	if len(db) == 0 {
		return m.db.Model(&SysRole{}).Session(&gorm.Session{Context: ctx})
	}
	return db[0]
}

func (m *defaultSysRoleModel) Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB) {
	return m.conn(ctx, db...)
}

func (m *defaultSysRoleModel) Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error) {
	return m.conn(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}

func (m *defaultSysRoleModel) Insert(ctx context.Context, data *SysRole, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Create(data).Error
	})
}

func (m *defaultSysRoleModel) Update(ctx context.Context, data *SysRole, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Model(&data).Updates(data).Error
	})
}

func (m *defaultSysRoleModel) Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Where(" id = ?", id).Updates(&SysRole{DeleteFlag: 1}).Error
	})
}

func (m *defaultSysRoleModel) ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Unscoped().Where(" id = ?", id).Delete(SysRole{}).Error
	})
}

func (m *defaultSysRoleModel) Count(ctx context.Context, cond *SysRoleQuery) (total int64, err error) {
	err = m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.SysRole).Count(&total).Error
	return total, err
}

func (m *defaultSysRoleModel) FindOne(ctx context.Context, id int64) (data *SysRole, err error) {
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

func (m *defaultSysRoleModel) FindByPage(ctx context.Context, cond *SysRoleQuery) (total int64, list []*SysRole, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.SysRole).Where("parent_id = ?", cond.ParentId)

	conn.Count(&total)
	err = conn.Scopes(
		pageScope(cond.PageCurrent, cond.PageSize),
		orderScope(cond.OrderSort...),
	).Find(&list).Error
	return total, list, err
}

func (m *defaultSysRoleModel) FindListByPage(ctx context.Context, cond *SysRoleQuery) (list []*SysRole, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		orderScope(cond.OrderSort...),
		pageScope(cond.PageCurrent, cond.PageSize),
	)
	err = conn.Where(cond.SysRole).Find(&list).Error
	return list, err
}

func (m *defaultSysRoleModel) FindListByCursor(ctx context.Context, cond *SysRoleQuery) (list []*SysRole, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		cursorScope(cond.Cursor, cond.CursorAsc, cond.PageSize, "id"),
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.SysRole).Find(&list).Error
	return list, err
}

func (m *defaultSysRoleModel) FindAll(ctx context.Context, cond *SysRoleQuery) (list []SysRole, err error) {
	conn := m.conn(ctx).Scopes(
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.SysRole).Find(&list).Error
	return list, err
}

func (m *defaultSysRoleModel) FindByPerms(ctx context.Context, perms string) (*SysRole, error) {
	var t SysRole
	tx := m.conn(ctx).Where(&SysRole{Perms: perms}).Find(&t)
	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}
	return &t, tx.Error
}

func (m *defaultSysRoleModel) FindByParentIds(ctx context.Context, pIds []int64) ([]SysRole, error) {
	var ts []SysRole
	err := m.conn(ctx).Where("parent_id in ?", pIds).Find(&ts).Error
	return ts, err
}

func (m *defaultSysRoleModel) FindByIds(ctx context.Context, ids []int64) ([]SysRole, error) {
	var ts []SysRole
	err := m.conn(ctx).Where("id in ?", ids).Find(&ts).Error
	return ts, err
}

func (m *defaultSysRoleModel) FindByPermsList(ctx context.Context, perms []string) ([]SysRole, error) {
	var ts []SysRole
	err := m.conn(ctx).Where("perms in ?", perms).Find(&ts).Error
	return ts, err
}
