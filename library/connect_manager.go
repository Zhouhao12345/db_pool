package library

import (
	"gorm_demo/library/types"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ConnectManager struct {
	Pool *PoolManager
}

func (p *ConnectManager) New() (err error) {
	p.Pool.init()
	return
}

func (p *ConnectManager) Get() (types.Connect, error) {
	var (
		cur types.Connect
		err error
	)

	if cur, err = p.Pool.Borrow(); err != nil {
		return nil, err
	}
	return cur, err
}

func (p *ConnectManager) Set(c types.Connect) (err error) {
	err = p.Pool.Back(c)
	return
}

func (p *ConnectManager) Context(fun types.HandlerFunc) types.HandlerFunc {

	return func(context *types.Context) {
		con, err := p.Get()
		if err != nil {
			panic(err)
		}
		con.HandlerRequest(fun)(context)
		p.Set(con)
	}
}
