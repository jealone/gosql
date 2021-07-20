package gosql

import (
	"bytes"
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func TestMonthlyDynamicCluster(t *testing.T) {
	c, cleanup, err := mockMonthlyDynamicCluster("config/demo.yml")
	if nil != err {
		t.Fatal(err)
	}
	defer cleanup()

	c.Write("test", []byte("test"), func(db *sql.DB, node int, table *bytes.Buffer) error {
		fmt.Println(table.String())

		buf := c.TableSelector().Pick("user", node, time.Now().Unix())
		defer ReleaseBuffer(buf)
		fmt.Println(buf.String())
		return nil
	}, time.Now().Unix())
}
