package importer

type storage interface {
	add(foreignID string, doc interface{}) error
	get(foreignID string) (interface{}, error)
	getPaginated(start, limit int) (map[string]map[string]interface{}, error)
	remove(foreignID string) error
	removeAll() error
	total() int
}
