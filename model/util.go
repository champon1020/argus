package model

import (
	"strconv"
	"time"
)

// IDType is the type for argument of GenNewID().
type IDType int

// Enum variables for GenNewID().
const (
	TypeArticle IDType = iota
	TypeDraft
	TypeCategory
)

// GetNewID generates the random id string.
func GetNewID(typ IDType) string {
	unix := time.Now().UnixNano()
	id := strconv.FormatInt(unix, 10)

	switch typ {
	case TypeArticle:
		id = "A" + id
	case TypeDraft:
		id = "D" + id
	case TypeCategory:
		id = "C" + id
	}

	return id
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
