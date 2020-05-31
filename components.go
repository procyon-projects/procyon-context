package context

import core "github.com/Rollcomp/procyon-core"

type RepositoryMetadata struct {
	typ *core.Type
}

type Repository interface {
	GetRepositoryMetadata() RepositoryMetadata
}

func NewRepositoryMetadata(typ *core.Type) RepositoryMetadata {
	return RepositoryMetadata{
		typ,
	}
}

type ServiceMetadata struct {
}

func NewServiceMetadata() ServiceMetadata {
	return ServiceMetadata{}
}

type Service interface {
	GetServiceMetadata() ServiceMetadata
}
