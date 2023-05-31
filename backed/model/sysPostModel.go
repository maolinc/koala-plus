package model

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

var (
	_ SysPostModel = (*defaultSysPostModel)(nil)
)

type (
	SysPost struct {
		Id         int64      `gorm:"id;primary_key"`                    //岗位ID
		PostCode   string     `gorm:"post_code"`                         // 岗位编码
		PostName   string     `gorm:"post_name"`                         // 岗位名称
		PostSort   int64      `gorm:"post_sort"`                         // 显示顺序
		Status     string     `gorm:"status"`                            // 状态（0正常 1停用）
		CreateBy   int64      `gorm:"create_by"`                         // 创建者
		CreateTime *time.Time `gorm:"column:create_time;autoCreateTime"` // 创建时间
		UpdateBy   int64      `gorm:"update_by"`                         // 更新者
		UpdateTime *time.Time `gorm:"column:update_time;autoUpdateTime"` // 更新时间
		Remark     string     `gorm:"remark"`                            // 备注
		AppId      int64      `gorm:"app_id"`                            //应用id
	}

	// query po
	SysPostQuery struct {
		SearchBase
		*SysPost
	}

	SysPostModel interface {
		// Trans Transaction
		Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error)
		// Builder Custom assembly conditions
		Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB)
		Insert(ctx context.Context, data *SysPost, db ...*gorm.DB) (err error)
		Update(ctx context.Context, data *SysPost, db ...*gorm.DB) (err error)
		// Delete When a tombstone field is set, it is tombstone. Otherwise, it is physically deleted
		Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		// ForceDelete Physical deletion
		ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		Count(ctx context.Context, cond *SysPostQuery) (total int64, err error)
		FindOne(ctx context.Context, id int64) (data *SysPost, err error)
		// FindByPage Contains total information
		FindByPage(ctx context.Context, cond *SysPostQuery) (total int64, list []*SysPost, err error)
		// FindListByPage Normal pagination
		FindListByPage(ctx context.Context, cond *SysPostQuery) (list []*SysPost, err error)
		// FindListByCursor Cursor is required based on cursor paging, Only the primary key is of type int, and other types can be expanded by themselves
		FindListByCursor(ctx context.Context, cond *SysPostQuery) (list []*SysPost, err error)
		FindAll(ctx context.Context, cond *SysPostQuery) (list []*SysPost, err error)
		// ---------------Write your other interfaces below---------------

		//FindByUserId(ctx context.Context, userId) (list []*SysPost, err error)
	}

	defaultSysPostModel struct {
		*customConn
		table string
	}
)

func NewSysPostModel(db *gorm.DB) SysPostModel {
	return &defaultSysPostModel{
		customConn: newCustomConnNoCache(db),
		table:      "sys_post",
	}
}

func (m *SysPost) TableName() string {
	return "`sys_post`"
}

func (m *defaultSysPostModel) conn(ctx context.Context, db ...*gorm.DB) *gorm.DB {
	if len(db) == 0 {
		return m.db.Model(&SysPost{}).Session(&gorm.Session{Context: ctx})
	}
	return db[0]
}

func (m *defaultSysPostModel) Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB) {
	return m.conn(ctx, db...)
}

func (m *defaultSysPostModel) Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error) {
	return m.conn(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}

func (m *defaultSysPostModel) Insert(ctx context.Context, data *SysPost, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Create(data).Error
	})
}

func (m *defaultSysPostModel) Update(ctx context.Context, data *SysPost, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Model(&data).Updates(data).Error
	})
}

func (m *defaultSysPostModel) Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Where(" id = ?", id).Delete(SysPost{}).Error
	})
}

func (m *defaultSysPostModel) ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Unscoped().Where(" id = ?", id).Delete(SysPost{}).Error
	})
}

func (m *defaultSysPostModel) Count(ctx context.Context, cond *SysPostQuery) (total int64, err error) {
	err = m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.SysPost).Count(&total).Error
	return total, err
}

func (m *defaultSysPostModel) FindOne(ctx context.Context, id int64) (data *SysPost, err error) {
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

func (m *defaultSysPostModel) FindByPage(ctx context.Context, cond *SysPostQuery) (total int64, list []*SysPost, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.SysPost)

	conn.Count(&total)
	err = conn.Scopes(
		pageScope(cond.PageCurrent, cond.PageSize),
		orderScope(cond.OrderSort...),
	).Find(&list).Error
	return total, list, err
}

func (m *defaultSysPostModel) FindListByPage(ctx context.Context, cond *SysPostQuery) (list []*SysPost, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		orderScope(cond.OrderSort...),
		pageScope(cond.PageCurrent, cond.PageSize),
	)
	err = conn.Where(cond.SysPost).Find(&list).Error
	return list, err
}

func (m *defaultSysPostModel) FindListByCursor(ctx context.Context, cond *SysPostQuery) (list []*SysPost, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		cursorScope(cond.Cursor, cond.CursorAsc, cond.PageSize, "id"),
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.SysPost).Find(&list).Error
	return list, err
}

func (m *defaultSysPostModel) FindAll(ctx context.Context, cond *SysPostQuery) (list []*SysPost, err error) {
	conn := m.conn(ctx).Scopes(
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.SysPost).Find(&list).Error
	return list, err
}
