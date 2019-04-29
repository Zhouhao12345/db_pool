package library

import (
	"gorm_demo/library/types"
	"gorm_demo/util"
	"sync"
)

type PoolManager struct {
	Pool       *util.Queue
	All        []types.Connect
	lock       *sync.Cond
	MinConnect int
	MaxConnect int
	New        func() (c types.Connect, err error)
}

func (p *PoolManager) init() {
	for i := 0; i < p.MinConnect; i++ {
		if cur, err := p.New(); err == nil {
			p.Append(cur)
		} else {
			panic(err)
		}
	}
}

func (p *PoolManager) Append(c types.Connect) (err error) {
	p.All = append(p.All, c)
	p.Pool.Append(c)
	return
}

func (p *PoolManager) Borrow() (cur types.Connect, err error) {
	p.lock.L.Lock()
	for {
		if p.Pool.End == nil {
			if len(p.All) >= p.MaxConnect {
				p.lock.Wait()
			} else {
				cur, err = p.New()
				if err != nil {
					panic(err)
				}
				cur.SetUsed(true)
				p.All = append(p.All, cur)
				p.lock.L.Unlock()
				return
			}
		} else {
			break
		}
	}
	var qn *util.QNode
	qn, err = p.Pool.Pop()
	if err != nil {
		return nil, err
	}
	cur = qn.Value.(types.Connect)

	cur.SetUsed(true)
	p.lock.L.Unlock()
	return
}

func (p *PoolManager) Back(c types.Connect) (err error) {
	p.lock.L.Lock()
	c.SetUsed(false)
	p.Pool.Append(c)
	p.lock.Broadcast()
	p.lock.L.Unlock()
	return
}
