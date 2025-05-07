package clickhouse

import (
	"context"
	"fmt"
	"sync"
	"time"

	"reflect"

	"github.com/AMuzykus/risor/arg"
	"github.com/AMuzykus/risor/errz"
	"github.com/AMuzykus/risor/object"
	"github.com/AMuzykus/risor/op"

	ch "github.com/ClickHouse/clickhouse-go/v2"	
)

const CH_CONN = object.Type("ch.conn")

type ChConn struct {
	ctx    context.Context
	conn   ch.Conn
	once   sync.Once
	closed chan bool
}

func (c *ChConn) Type() object.Type {
	return CH_CONN
}

func (c *ChConn) Inspect() string {
	return "pgx.conn()"
}

func (c *ChConn) Interface() interface{} {
	return c.conn
}

func (c *ChConn) Value() ch.Conn {
	return c.conn
}

func (c *ChConn) Equals(other object.Object) object.Object {
	return object.NewBool(c == other)
}

func (c *ChConn) IsTruthy() bool {
	return true
}

func (c *ChConn) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "query":
		return object.NewBuiltin("clickhouse.conn.query", c.Query), true
	case "exec", "execute": // "exec" for backwards compatibility
		return object.NewBuiltin("clickhouse.conn.execute", c.Exec), true
	case "close":
		return object.NewBuiltin("clickhouse.conn.close", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("clickhouse.conn.close", 0, args); err != nil {
				return err
			}
			if err := c.Close(); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	}
	return nil, false
}

func (c *ChConn) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("type error: pgx.conn object has no attribute %q", name)
}

func (c *ChConn) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for pgx.conn: %v", opType)
}

func (c *ChConn) Close() error {
	var err error
	c.once.Do(func() {
		err = c.conn.Close()
		close(c.closed)
	})
	return err
}

func (c *ChConn) waitToClose() {
	go func() {
		select {
		case <-c.closed:
		case <-c.ctx.Done():
			c.conn.Close()
		}
	}()
}

func (c *ChConn) Cost() int {
	return 8
}

func (c *ChConn) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal pgx.conn")
}

func New(ctx context.Context, conn ch.Conn) *ChConn {
	obj := &ChConn{
		ctx:    ctx,
		conn:   conn,
		closed: make(chan bool),
	}
	obj.waitToClose()
	return obj
}

func (c *ChConn) Query(ctx context.Context, args ...object.Object) object.Object {
	// The arguments should include a query string and zero or more query args
	if len(args) < 1 {
		return object.TypeErrorf("type error: pgx.conn.query() one or more arguments (%d given)", len(args))
	}
	query, errObj := object.AsString(args[0])
	if errObj != nil {
		return errObj
	}

	// Build list of query args as their Go types
	var queryArgs []interface{}
	for _, queryArg := range args[1:] {
		queryArgs = append(queryArgs, queryArg.Interface())
	}

	// Start the query
	rows, err := c.conn.Query(ctx, query, queryArgs...)
	if err != nil {
		return object.NewError(err)
	}
	defer rows.Close()

	// The field descriptions will tell us how to decode the result values
	var (
		columnTypes = rows.ColumnTypes()
		columnNames = rows.Columns()
		values        = make([]interface{}, len(columnTypes))
	)
	for i := range columnTypes {
		values[i] = reflect.New(columnTypes[i].ScanType()).Interface()
	}

	var results []object.Object

	// Transform each result row into a Risor map object
	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return object.NewError(err)
		}
		row := map[string]object.Object{}
		for colIndex, value := range values {
			key := columnNames[colIndex]
			var val object.Object
			switch value := value.(type) {
			case *string:
				val = object.NewString(*value)
			case *time.Time:
				val = object.NewTime(*value)
			case *int64:
				val = object.NewInt(*value)
			case *uint64:
				val = object.NewInt(int64(*value))
			case *float64:
				val = object.NewFloat(*value)
			}
			if val == nil {
				return object.TypeErrorf("type error: clickhouse.conn.query() encountered unsupported type: %T", value)
			}
			if !object.IsError(val) {
				row[key] = val
			} else {
				row[key] = object.NewString(fmt.Sprintf("__error__%s", value))
			}
		}
		results = append(results, object.NewMap(row))
	}
	return object.NewList(results)
}

func (c *ChConn) Exec(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.TypeErrorf("type error: clickhouse.conn.exec() one or more arguments (%d given)", len(args))
	}
	query, errObj := object.AsString(args[0])
	if errObj != nil {
		return errObj
	}
	var queryArgs []interface{}
	if len(args) == 2 {
		if list, ok := args[1].(*object.List); ok {
			for _, item := range list.Value() {
				queryArgs = append(queryArgs, item.Interface())
			}
		} else {
			queryArgs = append(queryArgs, args[1].Interface())
		}
	} else {
		for _, queryArg := range args[1:] {
			queryArgs = append(queryArgs, queryArg.Interface())
		}
	}
	err := c.conn.Exec(ctx, query, queryArgs...)
	if err != nil {
		return object.NewError(err)
	}
	return object.Nil
}
