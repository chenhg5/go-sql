package dialect

import "fmt"

type commonDialect struct {
	delimiter string
}

func (c commonDialect) Insert(comp *SQLComponent) string {
	comp.prepareInsert(c.delimiter)
	return comp.Statement
}

func (c commonDialect) Delete(comp *SQLComponent) string {
	comp.Statement = "delete from " + comp.TableName + comp.getWheres(c.delimiter)
	return comp.Statement
}

func (c commonDialect) Update(comp *SQLComponent) string {
	comp.prepareUpdate(c.delimiter)
	return comp.Statement
}

func (c commonDialect) Count(comp *SQLComponent) string {
	comp.prepareUpdate(c.delimiter)
	return comp.Statement
}

func (c commonDialect) Select(comp *SQLComponent) string {
	comp.Statement = "select " + comp.getFields(c.delimiter) + " from " + comp.TableName + comp.getJoins(c.delimiter) +
		comp.getWheres(c.delimiter) + comp.getGroupBy() + comp.getOrderBy() + comp.getLimit() + comp.getOffset()
	return comp.Statement
}

func (c commonDialect) ShowColumns(table string) string {
	return fmt.Sprintf("select column_name, udt_name from information_schema.columns where table_name = '%s'", table)
}

func (c commonDialect) GetName() string {
	return "common"
}

func (c commonDialect) ShowTables() string {
	return "show tables"
}

func (c commonDialect) GetDelimiter() string {
	return c.delimiter
}
