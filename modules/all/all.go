package all

import (
	"github.com/AMuzykus/risor/builtins"
	modBase64 "github.com/AMuzykus/risor/modules/base64"
	modBytes "github.com/AMuzykus/risor/modules/bytes"
	modColor "github.com/AMuzykus/risor/modules/color"
	modErrors "github.com/AMuzykus/risor/modules/errors"
	modExec "github.com/AMuzykus/risor/modules/exec"
	modFilepath "github.com/AMuzykus/risor/modules/filepath"
	modFmt "github.com/AMuzykus/risor/modules/fmt"
	modGha "github.com/AMuzykus/risor/modules/gha"
	modHTTP "github.com/AMuzykus/risor/modules/http"
	modIsTTY "github.com/AMuzykus/risor/modules/isatty"
	modJSON "github.com/AMuzykus/risor/modules/json"
	modMath "github.com/AMuzykus/risor/modules/math"
	modNet "github.com/AMuzykus/risor/modules/net"
	modOs "github.com/AMuzykus/risor/modules/os"
	modRand "github.com/AMuzykus/risor/modules/rand"
	modRegexp "github.com/AMuzykus/risor/modules/regexp"
	modStrconv "github.com/AMuzykus/risor/modules/strconv"
	modStrings "github.com/AMuzykus/risor/modules/strings"
	modTablewriter "github.com/AMuzykus/risor/modules/tablewriter"
	modTime "github.com/AMuzykus/risor/modules/time"
	modYAML "github.com/AMuzykus/risor/modules/yaml"
	"github.com/AMuzykus/risor/object"
)

func Builtins() map[string]object.Object {
	result := map[string]object.Object{
		"base64":      modBase64.Module(),
		"bytes":       modBytes.Module(),
		"color":       modColor.Module(),
		"errors":      modErrors.Module(),
		"exec":        modExec.Module(),
		"filepath":    modFilepath.Module(),
		"fmt":         modFmt.Module(),
		"gha":         modGha.Module(),
		"http":        modHTTP.Module(),
		"isatty":      modIsTTY.Module(),
		"json":        modJSON.Module(),
		"math":        modMath.Module(),
		"net":         modNet.Module(),
		"os":          modOs.Module(),
		"rand":        modRand.Module(),
		"regexp":      modRegexp.Module(),
		"strconv":     modStrconv.Module(),
		"strings":     modStrings.Module(),
		"tablewriter": modTablewriter.Module(),
		"time":        modTime.Module(),
		"yaml":        modYAML.Module(),
	}
	for k, v := range modHTTP.Builtins() {
		result[k] = v
	}
	for k, v := range modFmt.Builtins() {
		result[k] = v
	}
	for k, v := range builtins.Builtins() {
		result[k] = v
	}
	for k, v := range modOs.Builtins() {
		result[k] = v
	}
	return result
}
