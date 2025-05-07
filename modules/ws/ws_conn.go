package ws

import (
	"context"
	"sync"
	"time"

	// "github.com/AMuzykus/risor/arg"
	"github.com/AMuzykus/risor/errz"
	"github.com/AMuzykus/risor/object"
	"github.com/AMuzykus/risor/op"

	"github.com/coder/websocket"
)

const WS_CONN = object.Type("ws.conn")

type WsConn struct {
	ctx    context.Context
	conn   *websocket.Conn
	once   sync.Once
	closed chan bool
	reader chan object.Object
}

func (c *WsConn) Type() object.Type {
	return WS_CONN
}

func (c *WsConn) Inspect() string {
	return "ws.conn()"
}

func (c *WsConn) Interface() interface{} {
	return c.conn
}

func (c *WsConn) Value() *websocket.Conn {
	return c.conn
}

func (c *WsConn) Equals(other object.Object) object.Object {
	return object.NewBool(c == other)
}

func (c *WsConn) IsTruthy() bool {
	return true
}

func (c *WsConn) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "read":
		return object.NewBuiltin("ws.conn.read", c.Read), true
	case "write":
		return object.NewBuiltin("ws.conn.write", c.Write), true
	case "close":
		return object.NewBuiltin("ws.conn.close", c.Close), true
	}	
	return nil, false
}

func (c *WsConn) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("type error: ws.conn object has no attribute %q", name)
}

func (c *WsConn) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for ws.conn: %v", opType)
}

func (c *WsConn) Cost() int {
	return 8
}

func (c *WsConn) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal ws.conn")
}

func New(ctx context.Context, conn *websocket.Conn) *WsConn {
	obj := &WsConn{
		ctx:    ctx,
		conn:   conn,
		closed: make(chan bool),
		reader: make(chan object.Object),
	}
	obj.waitToClose()
	obj.read()
	return obj
}

// TODO: Корректно обработать: возврат ошибки вызывает исключение в defer
func (c *WsConn) close(code websocket.StatusCode, reason string) error {
	c.once.Do(func() {
		if err := c.conn.Close(code, reason); err != nil {
			c.conn.CloseNow()
		}
		close(c.closed)
	})
	return nil
}

func (c *WsConn) waitToClose() {
	go func() {
		select {
		case <-c.closed:
		case <-c.ctx.Done():
			c.close(websocket.StatusNormalClosure, "context done")
		}
	}()
}

func (c *WsConn) Close(ctx context.Context, args ...object.Object) object.Object {
	if err := c.close(websocket.StatusNormalClosure, "Close() call"); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func (c *WsConn) read() {
	go func() {
		for {
			msgType, buffer, err := c.conn.Read(c.ctx)
			if err != nil {
				c.reader <- object.NewList([]object.Object{object.Nil, object.Nil, object.NewError(err)})
				close(c.reader)
				break
			} else {
				c.reader <-object.NewList([]object.Object{object.NewInt(int64(msgType)), object.NewBufferFromBytes(buffer), object.Nil})
			}
		}
	}()
}

func (c *WsConn) Read(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.TypeErrorf("type error: ws.conn.read() takes exactly one arguments (%d given)", len(args))
	}
	timeout, errTimeout := object.AsInt(args[0])
	if errTimeout != nil {
		return object.NewError(errTimeout)
	}
	deadline := time.After(time.Second*time.Duration(timeout))

	select {
	case <- deadline:
		return object.NewList([]object.Object{object.Nil, object.Nil, object.Errorf("ws.conn.read() timeout")})
	case res := <-c.reader:
		return res
	}
}

func (c *WsConn) Write(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.TypeErrorf("type error: ws.conn.write() takes exactly two arguments (%d given)", len(args))
	}
	msgType, errMsgType := object.AsInt(args[0])
	if errMsgType != nil {
		return object.NewError(errMsgType)
	}

	buffer, errBuffer := object.AsBytes(args[1])
	if errBuffer != nil {
		return object.NewError(errBuffer)
	}

	err := c.conn.Write(ctx, websocket.MessageType(msgType), buffer)
	if err != nil {
		return object.NewError(err)
	}

	return object.Nil
}
