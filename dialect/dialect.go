package dialect

import (
	"strings"
)

// Dialect is methods set of different driver.
type Dialect interface {
	// GetName get dialect's name
	GetName() string

	// ShowColumns show columns of specified table
	ShowColumns(table string) string

	// ShowTables show tables of database
	ShowTables() string

	// Insert
	Insert(comp *SQLComponent) string

	// Delete
	Delete(comp *SQLComponent) string

	// Update
	Update(comp *SQLComponent) string

	// Select
	Select(comp *SQLComponent) string

	// GetDelimiter return the delimiter of Dialect.
	GetDelimiter() string
}

// GetDialectByDriver return the Dialect of given driver.
func GetDialectByDriver(driver string) Dialect {
	switch driver {
	case "mysql":
		return mysql{
			commonDialect: commonDialect{delimiter: "`"},
		}
	case "mssql":
		return mssql{
			commonDialect: commonDialect{delimiter: "`"},
		}
	case "postgresql":
		return postgresql{
			commonDialect: commonDialect{delimiter: `"`},
		}
	case "sqlite":
		return sqlite{
			commonDialect: commonDialect{delimiter: "`"},
		}
	default:
		return commonDialect{delimiter: "`"}
	}
}

// H is a shorthand of map.
type H map[string]interface{}

// SQLComponent is a sql components set.
type SQLComponent struct {
	Fields     []string
	Functions  []string
	TableName  string
	Wheres     []Where
	Leftjoins  []Join
	Args       []interface{}
	Order      string
	Offset     string
	Limit      string
	WhereRaws  string
	UpdateRaws []RawUpdate
	Statement  string
	Values     H
}

// Where contains the operation and field.
type Where struct {
	Operation string
	Field     string
	Qmark     string
}

// Join contains the table and field and operation.
type Join struct {
	Table     string
	FieldA    string
	Operation string
	FieldB    string
}

// RawUpdate contains the expression and arguments.
type RawUpdate struct {
	Expression string
	Args       []interface{}
}

// *******************************
// internal help function
// *******************************

func (sql *SQLComponent) getLimit() string {
	if sql.Limit == "" {
		return ""
	}
	return " limit " + sql.Limit + " "
}

func (sql *SQLComponent) getOffset() string {
	if sql.Offset == "" {
		return ""
	}
	return " offset " + sql.Offset + " "
}

func (sql *SQLComponent) getOrderBy() string {
	if sql.Order == "" {
		return ""
	}
	return " order by " + sql.Order + " "
}

func (sql *SQLComponent) getJoins(delimiter string) string {
	if len(sql.Leftjoins) == 0 {
		return ""
	}
	joins := ""
	for _, join := range sql.Leftjoins {
		joins += " left join " + wrap(delimiter, join.Table) + " on " + join.FieldA + " " + join.Operation + " " + join.FieldB + " "
	}
	return joins
}

func (sql *SQLComponent) getFields(delimiter string) string {
	if len(sql.Fields) == 0 {
		return "*"
	}
	fields := ""
	if len(sql.Leftjoins) == 0 {
		for k, field := range sql.Fields {
			if sql.Functions[k] != "" {
				fields += sql.Functions[k] + "(" + wrap(delimiter, field) + "),"
			} else {
				fields += wrap(delimiter, field) + ","
			}
		}
	} else {
		for _, field := range sql.Fields {
			arr := strings.Split(field, ".")
			if len(arr) > 1 {
				fields += arr[0] + "." + wrap(delimiter, arr[1]) + ","
			} else {
				fields += wrap(delimiter, field) + ","
			}
		}
	}
	return fields[:len(fields)-1]
}

func wrap(delimiter, field string) string {
	if field == "*" {
		return "*"
	}
	return delimiter + field + delimiter
}

func (sql *SQLComponent) getWheres(delimiter string) string {
	if len(sql.Wheres) == 0 {
		if sql.WhereRaws != "" {
			return " where " + sql.WhereRaws
		}
		return ""
	}
	wheres := " where "
	var arr []string
	for _, where := range sql.Wheres {
		arr = strings.Split(where.Field, ".")
		if len(arr) > 1 {
			wheres += arr[0] + "." + wrap(delimiter, arr[1]) + " " + where.Operation + " " + where.Qmark + " and "
		} else {
			wheres += wrap(delimiter, where.Field) + " " + where.Operation + " " + where.Qmark + " and "
		}
	}

	if sql.WhereRaws != "" {
		return wheres + sql.WhereRaws
	}
	return wheres[:len(wheres)-5]
}

func (sql *SQLComponent) prepareUpdate(delimiter string) {
	fields := ""
	args := make([]interface{}, 0)

	if len(sql.Values) != 0 {

		for key, value := range sql.Values {
			fields += wrap(delimiter, key) + " = ?, "
			args = append(args, value)
		}

		if len(sql.UpdateRaws) == 0 {
			fields = fields[:len(fields)-2]
		} else {
			for i := 0; i < len(sql.UpdateRaws); i++ {
				if i == len(sql.UpdateRaws)-1 {
					fields += sql.UpdateRaws[i].Expression + " "
				} else {
					fields += sql.UpdateRaws[i].Expression + ","
				}
				args = append(args, sql.UpdateRaws[i].Args...)
			}
		}

		sql.Args = append(args, sql.Args...)
	} else {
		if len(sql.UpdateRaws) == 0 {
			panic("prepareUpdate: wrong parameter")
		} else {
			for i := 0; i < len(sql.UpdateRaws); i++ {
				if i == len(sql.UpdateRaws)-1 {
					fields += sql.UpdateRaws[i].Expression + " "
				} else {
					fields += sql.UpdateRaws[i].Expression + ","
				}
				args = append(args, sql.UpdateRaws[i].Args...)
			}
		}
		sql.Args = append(args, sql.Args...)
	}

	sql.Statement = "update " + sql.TableName + " set " + fields + sql.getWheres(delimiter)
}

func (sql *SQLComponent) prepareInsert(delimiter string) {
	fields := " ("
	quesMark := "("

	for key, value := range sql.Values {
		fields += wrap(delimiter, key) + ","
		quesMark += "?,"
		sql.Args = append(sql.Args, value)
	}
	fields = fields[:len(fields)-1] + ")"
	quesMark = quesMark[:len(quesMark)-1] + ")"

	sql.Statement = "insert into " + sql.TableName + fields + " values " + quesMark
}
