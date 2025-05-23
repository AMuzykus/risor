package vm

import (
	"context"

	"github.com/AMuzykus/risor/builtins"
	"github.com/AMuzykus/risor/compiler"
	"github.com/AMuzykus/risor/importer"
	modBytes "github.com/AMuzykus/risor/modules/bytes"
	modErrors "github.com/AMuzykus/risor/modules/errors"
	modExec "github.com/AMuzykus/risor/modules/exec"
	modFmt "github.com/AMuzykus/risor/modules/fmt"
	modJSON "github.com/AMuzykus/risor/modules/json"
	modMath "github.com/AMuzykus/risor/modules/math"
	modOS "github.com/AMuzykus/risor/modules/os"
	modRand "github.com/AMuzykus/risor/modules/rand"
	modStrconv "github.com/AMuzykus/risor/modules/strconv"
	modStrings "github.com/AMuzykus/risor/modules/strings"
	modTime "github.com/AMuzykus/risor/modules/time"
	"github.com/AMuzykus/risor/object"
	"github.com/AMuzykus/risor/parser"
)

// Run the given code in a new Virtual Machine and return the result.
func Run(ctx context.Context, main *compiler.Code, options ...Option) (object.Object, error) {
	machine := New(main, options...)
	if err := machine.Run(ctx); err != nil {
		return nil, err
	}
	if result, exists := machine.TOS(); exists {
		return result, nil
	}
	return object.Nil, nil
}

type runOpts struct {
	Globals map[string]interface{}
}

// Run the given source code in a new VM. Used for testing.
func run(ctx context.Context, source string, opts ...runOpts) (object.Object, error) {
	vm, err := newVM(ctx, source, opts...)
	if err != nil {
		return nil, err
	}
	if err := vm.Run(ctx); err != nil {
		return nil, err
	}
	if result, exists := vm.TOS(); exists {
		return result, nil
	}
	return object.Nil, nil
}

// Return a new VM that's ready to run the given source code. Used for testing.
func newVM(ctx context.Context, source string, opts ...runOpts) (*VirtualMachine, error) {
	ast, err := parser.Parse(ctx, source)
	if err != nil {
		return nil, err
	}
	globals := basicBuiltins()
	if len(opts) > 0 {
		for k, v := range opts[0].Globals {
			globals[k] = v
		}
	}
	var globalNames []string
	for k := range globals {
		globalNames = append(globalNames, k)
	}
	main, err := compiler.Compile(ast, compiler.WithGlobalNames(globalNames))
	if err != nil {
		return nil, err
	}
	im := importer.NewLocalImporter(importer.LocalImporterOptions{
		SourceDirs:   []string{"./fixtures"},
		Extensions:  []string{".risor", ".rsr"},
		GlobalNames: globalNames,
	})
	return New(main, WithImporter(im), WithGlobals(globals), WithConcurrency()), nil
}

// Builtins to be used in VM tests.
func basicBuiltins() map[string]any {
	globals := map[string]any{
		"bytes":   modBytes.Module(),
		"exec":    modExec.Module(),
		"json":    modJSON.Module(),
		"errors":  modErrors.Module(),
		"math":    modMath.Module(),
		"rand":    modRand.Module(),
		"strconv": modStrconv.Module(),
		"strings": modStrings.Module(),
		"time":    modTime.Module(),
		"os":      modOS.Module(),
	}
	for k, v := range builtins.Builtins() {
		globals[k] = v
	}
	for k, v := range modFmt.Builtins() {
		globals[k] = v
	}
	return globals
}
