package context

type RepositoryMetadata struct {
}

type Repository interface {
	GetRepositoryMetadata() RepositoryMetadata
}

type ServiceMetadata struct {
}

type Service interface {
	GetServiceMetadata() ServiceMetadata
}
