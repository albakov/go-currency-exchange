package storage

import (
	"fmt"
)

var (
	EntityAlreadyExistsError = fmt.Errorf("entity already exists")
	EntitiesNotFoundError    = fmt.Errorf("entities not found")
)
