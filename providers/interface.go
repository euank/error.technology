package providers

import (
	"github.com/euank/api.error.technology/errortech"
	"github.com/euank/api.error.technology/providers/artisanal"
	"github.com/euank/api.error.technology/providers/ondisk"
)

type ErrorProvider interface {
	GetError(lang string, tags []string) errortech.Error
	Name() string
}

type Providers struct {
	FS        ErrorProvider
	Artisinal ErrorProvider
}

func NewDefaultProviders() Providers {
	return Providers{
		FS:        ondisk.New(),
		Artisinal: artisanal.New(),
	}
}

func (ps *Providers) Random() ErrorProvider {
	return ps.FS
}

func (ps *Providers) All() []ErrorProvider {
	return []ErrorProvider{ps.FS, ps.Artisinal}
}
