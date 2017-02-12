package client

type Address struct {
	Addr       string
	Path       string
	DataCenter string
	ACLToken   string
}

func (a *Address) fixupValues() {
	a.Path = fixPath(a.Path)
}
