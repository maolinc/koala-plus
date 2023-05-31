package casbin

import (
	"github.com/zeromicro/go-zero/core/fx"
	"strings"
)

type Permission struct {
	Role  string
	Dom   string
	Perms string
	Act   string
}

// FindPermissionAllByUser get user all permission
func (c *customCBS) FindPermissionAllByUser(user string, dom ...string) ([]*Permission, error) {
	ps, err := c.cb.GetImplicitPermissionsForUser(user, dom...)
	if err != nil {
		return nil, err
	}
	permsList := make([]*Permission, 0, len(ps))
	for _, p := range ps {
		permsList = append(permsList, &Permission{
			Role:  p[0],
			Dom:   p[1],
			Perms: p[2],
			Act:   p[3],
		})
	}
	return permsList, nil
}

// FindPermissionAllByRole get user all permission
func (c *customCBS) FindPermissionAllByRole(role string, dom ...string) ([]*Permission, error) {
	ps, err := c.cb.GetImplicitPermissionsForUser(role, dom...)
	if err != nil {
		return nil, err
	}
	permsList := make([]*Permission, 0, len(ps))
	for _, p := range ps {
		permsList = append(permsList, &Permission{
			Role:  p[0],
			Dom:   p[1],
			Perms: p[2],
			Act:   p[3],
		})
	}
	return permsList, nil
}

// FindMenusForRole 获取角色拥有的菜单
func (c *customCBS) FindMenusForRole(role, dom string) ([]string, error) {
	users, err := c.cb.GetImplicitPermissionsForUser(role, dom)
	if err != nil {
		return nil, err
	}
	result := filter(users, 2, CasMenu)
	return result, nil
}

func (c *customCBS) FindApisForRole(role, dom string) ([]string, error) {
	users, err := c.cb.GetImplicitPermissionsForUser(role, dom)
	if err != nil {
		return nil, err
	}
	result := filter(users, 2, CasApi)
	return result, nil
}

func (c *customCBS) FindMenusForUser(userId, dom string) ([]string, error) {
	users, err := c.cb.GetImplicitPermissionsForUser(userId, dom)
	if err != nil {
		return nil, err
	}
	result := filter(users, 2, CasMenu)
	return result, nil
}

func (c *customCBS) FindRolesForUser(userId, dom string) []string {
	roles := c.cb.GetRolesForUserInDomain(userId, dom)
	return filter2(roles, CasRole)
}

func (c *customCBS) FindApisForUser(userId, dom string) ([]string, error) {
	users, err := c.cb.GetImplicitPermissionsForUser(userId, dom)
	if err != nil {
		return nil, err
	}
	result := filter(users, 2, CasApi)
	return result, nil
}

func filter(sources [][]string, index int, filterPre string) []string {
	res := make([]string, 0)

	fx.From(func(source chan<- interface{}) {
		for _, u := range sources {
			source <- u
		}
	}).Filter(func(item interface{}) bool {
		it := item.([]string)[index]
		return strings.HasPrefix(it, filterPre)
	}).ForAll(func(pipe <-chan interface{}) {
		for item := range pipe {
			res = append(res, item.(string))
		}
	})

	return res
}

func filter2(sources []string, filterPre string) []string {
	res := make([]string, 0)

	fx.From(func(source chan<- interface{}) {
		for _, u := range sources {
			source <- u
		}
	}).Filter(func(item interface{}) bool {
		it := item.(string)
		return strings.HasPrefix(it, filterPre)
	}).ForAll(func(pipe <-chan interface{}) {
		for item := range pipe {
			res = append(res, item.(string))
		}
	})

	return res
}
