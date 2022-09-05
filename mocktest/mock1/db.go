package mock1

import "errors"

type Repository interface {
	Create(key string, value []byte) error
	Retrieve(key string) ([]byte, error)
	Update(key string, value []byte) error
	Delete(key string) error
}

func CreateRepo(repo Repository, key string, value string) (string, error) {
	err := repo.Create(key, []byte(value))
	if err != nil {
		errStr := err.Error() + "Create failed, key= " + key + "value= " + value
		return "", errors.New(errStr)
	}
	val, err := repo.Retrieve(key)
	return string(val), err
}
