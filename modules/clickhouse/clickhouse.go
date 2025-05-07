package clickhouse

import (
	"context"

	"github.com/AMuzykus/risor/object"

	ch "github.com/ClickHouse/clickhouse-go/v2"
)

func Open(ctx context.Context, args ...object.Object) object.Object {

	if len(args) != 1 {
		return object.TypeErrorf("type error: clickhouse.open() takes exactly one argument (%d given)", len(args))
	}
	options, ok := args[0].(*object.Map)
	if !ok {
		return object.TypeErrorf("type error: clickhouse.open() expected a map argument (got %s)", args[0].Type())
	}

	chOptions, err := NewClickHouseOptions(options)
	if err != nil {
		return object.NewError(err)
	}

	conn, err := ch.Open(chOptions)
	if err != nil {
		return object.NewError(err)
	}
	return New(ctx, conn)
}

// Module returns the `ch` module object
func Module() *object.Module {
	return object.NewBuiltinsModule("clickhouse", map[string]object.Object{
		"open": object.NewBuiltin("open", Open),
	})
}
