package gotest

type DB interface {
	Get(key string) (int, error)
}

func GetFromDB(db DB, key string) int {
	value, err := db.Get(key)
	if err != nil {
		return -1
	}
	return value
}
