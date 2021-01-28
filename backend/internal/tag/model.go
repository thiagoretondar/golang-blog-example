package tag

import "github.com/google/uuid"

type Tag struct {
	ID        *uuid.UUID
	Name      *string
	CreatedAt *string
	UpdatedAt *string
}
