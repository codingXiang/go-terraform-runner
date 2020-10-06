package command

type Meta struct {
	Label    string
	Path     string
	Commands []string
}

func NewMeta() *Meta {
	return &Meta{
		Commands: make([]string, 0),
	}
}
