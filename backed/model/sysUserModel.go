package model

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

var (
	_ SysUserModel = (*defaultSysUserModel)(nil)
)

type (
	SysUser struct {
		Id             int64                 `gorm:"id;primary_key"`                     // 用户ID
		DeptId         int64                 `gorm:"dept_id"`                            // 部门ID
		UserName       string                `gorm:"user_name"`                          // 用户账号
		NickName       string                `gorm:"nick_name"`                          // 用户昵称
		UserType       string                `gorm:"user_type"`                          // 用户类型（00系统用户）
		Email          string                `gorm:"email"`                              // 用户邮箱
		Phone          string                `gorm:"Phone"`                              // 手机号码
		Sex            string                `gorm:"sex"`                                // 用户性别（0男 1女 2未知）
		Avatar         string                `gorm:"avatar"`                             // 头像地址
		Password       string                `gorm:"password"`                           // 密码
		Salt           string                `gorm:"salt"`                               // 加密盐
		Status         string                `gorm:"status"`                             // 帐号状态（0正常 1停用）
		DeleteFlag     soft_delete.DeletedAt `gorm:"column:delete_flag;softDelete:flag"` // 删除标志（0代表存在 2代表删除）
		LoginIp        string                `gorm:"login_ip"`                           // 最后登录IP
		LoginDate      *time.Time            `gorm:"login_date"`                         // 最后登录时间
		CreateBy       int64                 `gorm:"create_by"`                          // 创建者
		CreateTime     *time.Time            `gorm:"column:create_time;autoCreateTime"`  // 创建时间
		UpdateBy       int64                 `gorm:"update_by"`                          // 更新者
		UpdateTime     *time.Time            `gorm:"column:update_time;autoUpdateTime"`  // 更新时间
		Remark         string                `gorm:"remark"`                             // 备注
		Dept           SysDept               `gorm:"foreignKey:Id;references:DeptId"`    // 部门
		Perms          string                `gorm:"perms"`                              // 角色权限字符串
		OrganizationId int64                 `gorm:"organization_id"`                    // 角色权限字符串
	}

	// SysUserQuery query po
	SysUserQuery struct {
		SearchBase
		*SysUser
		DeptIds []int64 `json:"deptIds"`
	}

	SysUserModel interface {
		// Trans Transaction
		Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error)
		// Builder Custom assembly conditions
		Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB)
		Insert(ctx context.Context, data *SysUser, db ...*gorm.DB) (err error)
		Update(ctx context.Context, data *SysUser, db ...*gorm.DB) (err error)
		// Delete When a tombstone field is set, it is tombstone. Otherwise, it is physically deleted
		Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		// ForceDelete Physical deletion
		ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		Count(ctx context.Context, cond *SysUserQuery) (total int64, err error)
		FindOne(ctx context.Context, id int64) (data *SysUser, err error)
		// FindByPage Contains total information
		FindByPage(ctx context.Context, cond *SysUserQuery) (total int64, list []*SysUser, err error)
		// FindListByPage Normal pagination
		FindListByPage(ctx context.Context, cond *SysUserQuery) (list []*SysUser, err error)
		// FindListByCursor Cursor is required based on cursor paging, Only the primary key is of type int, and other types can be expanded by themselves
		FindListByCursor(ctx context.Context, cond *SysUserQuery) (list []*SysUser, err error)
		FindAll(ctx context.Context, cond *SysUserQuery) (list []*SysUser, err error)
		// ---------------Write your other interfaces below---------------

		FindByPerms(ctx context.Context, perms string) (data *SysUser, err error)
		FindByPhone(ctx context.Context, phone string) (*SysUser, error)
		FindByUserName(ctx context.Context, name string) (*SysUser, error)
	}

	defaultSysUserModel struct {
		*customConn
		table string
	}
)

func NewSysUserModel(db *gorm.DB) SysUserModel {
	return &defaultSysUserModel{
		customConn: newCustomConnNoCache(db),
		table:      "sys_user",
	}
}

func (m *SysUser) TableName() string {
	return "`sys_user`"
}

func (m *defaultSysUserModel) conn(ctx context.Context, db ...*gorm.DB) *gorm.DB {
	if len(db) == 0 {
		return m.db.Model(&SysUser{}).Session(&gorm.Session{Context: ctx})
	}
	return db[0]
}

func (m *defaultSysUserModel) Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB) {
	return m.conn(ctx, db...)
}

func (m *defaultSysUserModel) Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error) {
	return m.conn(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}

func (m *defaultSysUserModel) Insert(ctx context.Context, data *SysUser, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Create(data).Error
	})
}

func (m *defaultSysUserModel) Update(ctx context.Context, data *SysUser, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Model(&data).Updates(data).Error
	})
}

func (m *defaultSysUserModel) Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Where(" id = ?", id).Updates(SysUser{DeleteFlag: 1}).Error
	})
}

func (m *defaultSysUserModel) ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Unscoped().Where(" id = ?", id).Delete(SysUser{}).Error
	})
}

func (m *defaultSysUserModel) Count(ctx context.Context, cond *SysUserQuery) (total int64, err error) {
	err = m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.SysUser).Count(&total).Error
	return total, err
}

func (m *defaultSysUserModel) FindOne(ctx context.Context, id int64) (data *SysUser, err error) {
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

func (m *defaultSysUserModel) FindByPage(ctx context.Context, cond *SysUserQuery) (total int64, list []*SysUser, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.SysUser)

	conn.Count(&total)
	err = conn.Scopes(
		pageScope(cond.PageCurrent, cond.PageSize),
		orderScope(cond.OrderSort...),
	).Find(&list).Error
	return total, list, err
}

func (m *defaultSysUserModel) FindListByPage(ctx context.Context, cond *SysUserQuery) (list []*SysUser, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		orderScope(cond.OrderSort...),
		pageScope(cond.PageCurrent, cond.PageSize),
	)
	err = conn.Where(cond.SysUser).Find(&list).Error
	return list, err
}

func (m *defaultSysUserModel) FindListByCursor(ctx context.Context, cond *SysUserQuery) (list []*SysUser, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		cursorScope(cond.Cursor, cond.CursorAsc, cond.PageSize, "id"),
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.SysUser).Find(&list).Error
	return list, err
}

func (m *defaultSysUserModel) FindAll(ctx context.Context, cond *SysUserQuery) (list []*SysUser, err error) {
	conn := m.conn(ctx).Scopes(
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.SysUser).Find(&list).Error
	return list, err
}

func (m *defaultSysUserModel) FindByPhone(ctx context.Context, phone string) (data *SysUser, err error) {
	tx := m.conn(ctx).Model(&SysUser{}).Where("phone = ?", phone).Find(&data)
	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}
	return data, tx.Error
}

func (m *defaultSysUserModel) FindByUserName(ctx context.Context, name string) (data *SysUser, err error) {
	tx := m.conn(ctx).Model(&SysUser{}).Where("user_name = ?", name).Find(&data)
	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}
	return data, tx.Error
}

func (m *defaultSysUserModel) FindByPerms(ctx context.Context, perms string) (data *SysUser, err error) {
	tx := m.conn(ctx).Where(&SysUser{Perms: perms}).Find(&data)
	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}
	return data, tx.Error
}
