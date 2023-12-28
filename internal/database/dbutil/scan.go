package dbutil

type Scanner interface {
	Scan(dest ...any) error
}
