package idx

import (
	"github.com/google/uuid"
)

// GenUUID 生成UUID
func GenUUID() string {
	id := uuid.New().String()
	return id
}
