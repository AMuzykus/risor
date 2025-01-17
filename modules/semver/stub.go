//go:build !semver
// +build !semver

package semver

import (
	"github.com/AMuzykus/risor/object"
)

func Module() *object.Module {
	return nil
}
