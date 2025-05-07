package clickhouse

import (
	"context"
	"testing"

	// "github.com/AMuzykus/risor/limits"
	"github.com/AMuzykus/risor"
	"github.com/AMuzykus/risor/object"	

	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	conn := Open(context.Background(), object.NewMap(map[string]object.Object{
		"addr": object.NewString("dev.albacore.ru:9000"),
		"auth": object.NewMap(map[string]object.Object{
			"database": object.NewString("bmstu"),
			"username": object.NewString("albacore"),
			"password": object.NewString("Albacore11"),
		}),
	}))
	require.NotEqual(t, conn.Type(), object.ERROR, conn)

	var (
		query object.Object
		ok bool
	)

	if query, ok = conn.GetAttr("query"); !ok {
		t.Error("Can't get 'query' method of clickhouse.conn")
		t.FailNow()
	}

	queryFn, _ := query.(*object.Builtin)
	res := queryFn.Call(context.Background(), object.NewString("select count(*) from history"))

	t.Log(res.Type())
	t.Log(res)

	res, err := risor.Eval(context.Background(), `

import clickhouse

print("--- START ---")

request := http.get("127.0.0.1")

con := clickhouse.open({
	addr: "dev.albacore.ru:9000",
	auth: {
		database: "bmstu",
		username: "albacore",
		password: "Albacore11"
	}
})

res, err := try(func() {[con.query("select count(*) as rows from history"), nil]}, func(err) {[nil, err]})
print("err=", err)
print("res=", res)

print("--- STOP ---")
err == nil ? res[0]["rows"] : nil
	`, risor.WithGlobal("clickhouse", Module()))

	if err != nil {
		t.Log(err)
	} else {
		t.Log(res)
	}

	require.Equal(t, nil, err)
}

