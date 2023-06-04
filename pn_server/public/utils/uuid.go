package utils

import "github.com/google/uuid"

func UUID() string {
	return uuid.NewString()
}
