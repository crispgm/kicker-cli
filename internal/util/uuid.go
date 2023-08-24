// Package util .
package util

import "github.com/google/uuid"

// UUID generates a string UUID. Wrap of google/uuid
func UUID() string {
	return uuid.NewString()
}
