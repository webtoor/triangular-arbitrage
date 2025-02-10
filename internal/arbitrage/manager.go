package arbitrage

type manager struct {
	exchange ArbitrageList
}

func New() *manager {
	return &manager{
		exchange: ArbitrageList{},
	}
}

func (p *manager) Registry(name string, t Triangular) *manager {
	_, ok := p.exchange[name]
	if !ok {
		p.exchange[name] = t
	}

	return p
}

func (p *manager) Resolve(name string) Triangular {
	h, ok := p.exchange[name]
	if ok {
		return h
	}

	return &noop{}
}
