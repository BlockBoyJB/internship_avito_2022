package postgres

type Option func(postgres *Postgres)

func MaxPoolSize(size int) Option {
	return func(postgres *Postgres) {
		postgres.maxPoolSize = size
	}
}
