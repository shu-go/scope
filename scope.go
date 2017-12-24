package scope

type scope struct {
	begin, end func()
}

func New(begin, end func()) scope {
	return scope{
		begin: begin,
		end:   end,
	}
}

func (s scope) Block(body func()) {
	s.begin()
	body()
	s.end()
}

func (s scope) Defer() func() {
	s.begin()
	return s.end
}
