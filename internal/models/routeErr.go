package models

type DomainErr struct {
	err string
}

func (d DomainErr) Error() string {
	return d.err
}
