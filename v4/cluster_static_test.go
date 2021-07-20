package gosql

import (
	"bytes"
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func TestCluster(t *testing.T) {
	c, cleanup, err := mockStaticCluster("config/demo.yml")
	if nil != err {
		t.Fatal(err)
	}
	defer cleanup()

	c.Write("test", []byte("test"), func(db *sql.DB, node int, table *bytes.Buffer) error {
		fmt.Println(table.String())

		buf := c.TableSelector().Pick("user", node)
		defer ReleaseBuffer(buf)
		fmt.Println(buf.String())
		return nil
	})
}

func TestMonthlyCluster(t *testing.T) {
	c, cleanup, err := mockMonthlyStaticCluster("config/demo.yml")
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

func TestDailyCluster(t *testing.T) {
	c, cleanup, err := mockDailyStaticCluster("config/demo.yml")
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

func TestAnnuallyCluster(t *testing.T) {
	c, cleanup, err := mockAnnuallyStaticCluster("config/demo.yml")
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
