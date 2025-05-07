package ws

import (
	"context"

	"github.com/AMuzykus/risor/object"

	"github.com/coder/websocket"
)

func Dial(ctx context.Context, args ...object.Object) object.Object {

	if len(args) != 1 {
		return object.TypeErrorf("type error: ws.dial() takes exactly one argument (%d given)", len(args))
	}
	url, ok := args[0].(*object.String)
	if !ok {
		return object.TypeErrorf("type error: ws.open() expected a string argument (got %s)", args[0].Type())
	}

	conn, _, err := websocket.Dial(ctx, url.String(), nil)
	if err != nil {
		return object.NewError(err)
	}
	return New(ctx, conn)
}

// Module returns the `ch` module object
func Module() *object.Module {
	return object.NewBuiltinsModule("ws", map[string]object.Object{
		"dial": object.NewBuiltin("dial", Dial),
	})
}
