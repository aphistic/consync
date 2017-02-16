package client

type Address struct {
	Addr       string
	Scheme     string
	Path       string
	DataCenter string
	ACLToken   string
}

func (a *Address) fixupValues() {
	a.Path = fixPath(a.Path)
}
