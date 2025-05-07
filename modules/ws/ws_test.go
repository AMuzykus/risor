package ws

import (
	"context"
	// "errors"
	"strings"

	// "fmt"
	"net/http"
	"testing"
	"time"

	// "github.com/AMuzykus/risor/limits"
	"github.com/AMuzykus/risor"
	// "github.com/AMuzykus/risor/object"

	"github.com/coder/websocket"
	// "github.com/coder/websocket/wsjson"

	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {

	handler := func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			t.Errorf("can't accept ws connection: %v", err)
		}
		defer c.CloseNow()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(30))
		defer cancel()
	
		for {
			_, buffer, err := c.Read(ctx)
			if err != nil {
				t.Logf("reading ws msg failed: %v", err)
				break			
			}
			t.Logf("ws msg received: %s", string(buffer))
			if string(buffer) == "STOP" {
				break
			}
			time.Sleep(1500*time.Millisecond)
			if strings.HasPrefix(string(buffer), "PING") {
				c.Write(ctx, websocket.MessageText, []byte("PONG"))
			}
		}
	
		c.Close(websocket.StatusNormalClosure, "")
	}
	
	http.HandleFunc("/", handler)
	go func() {http.ListenAndServe(":8090", nil)}()

	// ----------------------

	/*
	time.Sleep(2*time.Second)

	conn := Dial(context.Background(), object.NewString("ws://127.0.0.1:8090/"))

	require.NotEqual(t, cПonn.Type(), object.ERROR, conn)

	var (
		write object.Object
		ok bool
	)

	if write, ok = conn.GetAttr("write"); !ok {
		t.Error("Can't get 'write' method of ws.conn")
		t.FailNow()
	}

	writeFn, _ := write.(*object.Builtin)

	for i := range(10) {
		res := writeFn.Call(context.Background(), object.NewInt(1), object.NewString(fmt.Sprintf("PING%d",i)))
		require.NotEqual(t, res.Type(), object.ERROR, res)
	}
	writeFn.Call(context.Background(), object.NewInt(1), object.NewString("STOP"))
	*/

	res, err := risor.Eval(context.Background(), `

		// import ws

		print("--- START")
		time.sleep(1)

		func test() {
			con := ws.dial("ws://127.0.0.1:8090/")
			defer con.close()

			for i := range(10) {
				print('--- iteration # {i}')
				err_write := con.write(1, 'PING{i}')
				if err_write != nil {
					print(err_write)
					break
				}
				_, resp, err_read := con.read(1)
				if err_read != nil {
					print(err_read)
				} else {
					print('------ received response: {string(resp)}')
				}
			}
			con.write(1, 'STOP')
		}

		test()

		print("--- STOP")

		"Ok"

		`, risor.WithGlobal("ws", Module()))

	require.Equal(t, nil, err)

	t.Logf("result: %v", res)
}

