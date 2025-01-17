//go:build !vault
// +build !vault

package vault

import (
	"github.com/AMuzykus/risor/object"
)

func Module() *object.Module {
	return nil
}
