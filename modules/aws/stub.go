//go:build !aws
// +build !aws

package aws

import (
	"github.com/AMuzykus/risor/object"
)

func Module() *object.Module {
	return nil
}
