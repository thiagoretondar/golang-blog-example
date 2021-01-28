package post

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID       *uuid.UUID
	Title    *string
	Content  *string
	AuthorID *int
	TagsID   *[]int
	CreateAt *time.Time
	UpdateAt *time.Time
}
