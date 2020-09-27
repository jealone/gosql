package gosql

import (
	"database/sql"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func New(config *Config) (db *sql.DB, err error) {

	db, err = sql.Open("mysql", config.GetConfigString())
	return db, err
}

type Config struct {
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Name     string `json:"db_name,omitempty"`
	Charset  string `json:"charset,omitempty"`
}

func (c *Config) GetConfigString() string {
	if "" == c.Charset {
		c.Charset = "utf8"
	}
	return c.User + ":" + c.Password + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.Name + "?charset=" + c.Charset
}

func FetchAll(rows *sql.Rows, total int) ([]map[string]string, error) {

	columns, err := rows.Columns()
	if nil != err {
		return nil, err
	}

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	var result map[string]string
	data := make([]map[string]string, 0, total)

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if nil != err {
			return nil, err
		}
		var rst string
		result = make(map[string]string, len(columns))
		for i, col := range values {
			switch item := col.(type) {
			case int64:
				rst = strconv.FormatInt(item, 10)
			case string:
				rst = item
			case []byte:
				rst = string(item)
			case float64:
				strconv.FormatFloat(item, 'f', -1, 64)
			case bool:
				strconv.FormatBool(item)
			case time.Time:
				rst = item.Format(timeFormat)
			case nil:
				rst = ""
			}
			result[columns[i]] = rst
		}
		data = append(data, result)
	}

	return data, nil

}

func FetchRow(rows *sql.Rows) (map[string]string, error) {

	columns, err := rows.Columns()
	if nil != err {
		return nil, err
	}

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	result := make(map[string]string)

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if nil != err {
			return nil, err
		}
		var rst string
		for i, col := range values {
			if nil == col {
				continue
			}
			switch item := col.(type) {
			case int64:
				rst = strconv.FormatInt(item, 10)
			case string:
				rst = item
			case []byte:
				rst = string(item)
			case float64:
				strconv.FormatFloat(item, 'f', -1, 64)
			case bool:
				strconv.FormatBool(item)
			case time.Time:
				rst = item.Format(timeFormat)
			}
			result[columns[i]] = rst
		}
		break
	}

	return result, nil

}
