package casbin

import "github.com/zeromicro/go-zero/core/logx"

func (c *customCBS) BindRoleRole(role1, role2, dom string) (bool, error) {
	//sub := ToRole(role1)
	//obj := ToApi(role2)
	//domain := ToDom(dom)
	return c.cb.AddGroupingPolicy(role1, role2, dom)
}

func (c *customCBS) RemoveBindRoleRole(role1, role2, dom string) (bool, error) {
	//sub := ToRole(role1)
	//obj := ToApi(role2)
	//domain := ToDom(dom)
	return c.cb.RemoveGroupingPolicy(role1, role2, dom)
}

func (c *customCBS) BindRoleUser(role, user, dom string) (bool, error) {
	//sub := ToRole(user)
	//obj := ToApi(role)
	//domain := ToDom(dom)
	return c.cb.AddGroupingPolicy(user, role, dom)
}

func (c *customCBS) RemoveBindRoleUser(role, user, dom string) (bool, error) {
	//sub := ToRole(user)
	//obj := ToApi(role)
	//domain := ToDom(dom)
	return c.cb.RemoveGroupingPolicy(user, role, dom)
}

func (c *customCBS) BindRoleApi(role, api, dom, act string) (bool, error) {
	//sub := ToRole(role)
	//obj := ToApi(api)
	//domain := ToDom(dom)
	return c.cb.AddPolicy(role, api, dom, act)
}

func (c *customCBS) RemoveBindRoleApi(role, api, dom, act string) (bool, error) {
	//sub := ToRole(role)
	//obj := ToApi(api)
	//domain := ToDom(dom)
	return c.cb.RemovePolicy(role, dom, api, act)
}

func (c *customCBS) BindRoleMenu(role, menu, dom, act string) (bool, error) {
	//sub := ToRole(role)
	//obj := ToMenu(menu)
	//domain := ToDom(dom)
	return c.cb.AddPolicy(role, dom, menu, act)
}

func (c *customCBS) RemoveBindRoleMenu(role, menu, dom, act string) (bool, error) {
	//sub := ToRole(role)
	//obj := ToMenu(menu)
	//domain := ToDom(dom)
	return c.cb.RemovePolicy(role, dom, menu, act)
}

func (c *customCBS) DeletePolicyWithPerms(perms, dom string) (bool, error) {
	return c.cb.RemoveFilteredPolicy(1, dom, perms)
}

func (c *customCBS) DeletePolicyWithAct(perms, act, dom string) (bool, error) {
	return c.cb.RemoveFilteredPolicy(1, dom, perms, act)
}

func (c *customCBS) UpdatePolicyWithPerms(oldPerms, newPerms, dom string) bool {
	err := c.db.Where("ptype = p").Where("v1 = ? and v2 = ?", dom, oldPerms).Update("v2", newPerms).Error
	if err != nil {
		logx.Errorf("UpdatePolicyWithPerms update error: %+v", err)
		return false
	}
	err = c.Reload()
	if err != nil {
		logx.Errorf("UpdatePolicyWithPerms reload error: %+v", err)
		return false
	}
	return true
}

func (c *customCBS) UpdatePolicy(old *Permission, new Permission, dom string) bool {
	if old.Dom == "" || old.Perms == "" {
		return false
	}
	conn := c.db.conn()
	if old.Role != "" {
		conn = conn.Where("v0 = ?", old.Role)
	}
	if old.Act != "" {
		conn = conn.Where("v3 = ?", old.Act)
	}
	conn = conn.Where("ptype = p AND v1 = ? AND v2 = ?", dom, old.Perms)
	err := conn.Updates(&CasbinRule{
		V0: new.Role,
		V2: new.Perms,
		V3: new.Act,
	}).Error
	if err != nil {
		logx.Errorf("UpdatePolicy update error: %+v", err)
		return false
	}
	err = c.Reload()
	if err != nil {
		logx.Errorf("UpdatePolicy reload error: %+v", err)
		return false
	}
	return true
}

func (c *customCBS) ClearApp(dom string) (res bool) {
	res = true
	// delete all user and role
	_, err := c.cb.DeleteDomains(dom)
	if err != nil {
		logx.Errorf("ClearApp clear user and role error,appPerms: %s, error: %+v", dom, err)
		res = false
	}
	// delete all policy
	_, err = c.cb.RemoveFilteredPolicy(1, dom)
	if err != nil {
		logx.Errorf("ClearApp clear policy error,appPerms: %s, error: %+v", dom, err)
		res = false
	}

	return res
}
