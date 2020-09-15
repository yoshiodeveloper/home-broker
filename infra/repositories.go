package infra

// Repositories holds all repositories.
// Repositories know how to interact with the database.
type Repositories struct {
	UserRepo   UserRepoI
	WalletRepo WalletRepoI
}

// RepositoryI is an interface for a repository.
type RepositoryI interface{}
