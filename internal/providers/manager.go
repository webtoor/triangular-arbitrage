package providers

type manager struct {
	exchange ExchangeList
}

func New() *manager {
	return &manager{
		exchange: ExchangeList{},
	}
}

func (p *manager) Registry(name string, sapi SpotAPI) *manager {
	_, ok := p.exchange[name]
	if !ok {
		p.exchange[name] = sapi
	}

	return p
}

func (p *manager) SetExchange(name string) SpotAPI {
	h, ok := p.exchange[name]
	if ok {
		return h
	}

	return &noop{}
}
