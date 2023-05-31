package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"sync"
)

const (
	CasRole = "role:"
	CasMenu = "menu:"
	CasApi  = "api:"
	CasUser = "user:"
	CasDom  = "dom:"

	SpiltSymbol = ":"
)

type CBS struct {
	CB *casbin.CachedEnforcer
	*customCBS
}

type customCBS struct {
	db   *dbClient
	cb   *casbin.CachedEnforcer
	lock sync.Mutex
}

func newCustomCBS(db *gorm.DB, cd *casbin.CachedEnforcer) *customCBS {
	return &customCBS{db: newDBClient(db), cb: cd}
}

func NewCBS(db *gorm.DB) *CBS {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(err)
	}
	cachedEnforcer, err := casbin.NewCachedEnforcer(genModel(), adapter)
	if err != nil {
		panic(err)
	}
	err = cachedEnforcer.LoadPolicy()
	if err != nil {
		panic(err)
	}
	return &CBS{
		cachedEnforcer,
		newCustomCBS(db, cachedEnforcer),
	}
}

func genModel() model.Model {
	text := `
			[request_definition]
			r = sub, dom, obj, act
			
			[policy_definition]
			p = sub, dom, obj, act
			
			[role_definition]
			g = _, _, _
			g2 = _, _, _
			
			[policy_effect]
			e = some(where (p.eft == allow))
			
			[matchers]
			m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && (keyMatch2(r.obj,p.obj) || g2(r.obj, p.obj, r.dom)) && r.act== p.act
		`
	m, err := model.NewModelFromString(text)
	if err != nil {
		panic("String loading model failed!")
	}
	return m
}

func (c *customCBS) Reload() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cb.LoadPolicy()
}

func ToUser(user string) string {
	return CasUser + user
}

func ToRole(role string) string {
	return CasRole + role
}

func ToMenu(menu string) string {
	return CasMenu + menu
}

func ToApi(api string) string {
	return CasApi + api
}

func ToDom(dom string) string {
	return CasDom + dom
}
