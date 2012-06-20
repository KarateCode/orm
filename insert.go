package orm

import (
	"strings"
	"database/sql"
	_ "github.com/ziutek/mymysql/godrv"
)

type Stmt struct {
	stmt *sql.Stmt
}

func (self *Model) PrepareInsert(fields []string) (*Stmt, error) {
	var questionMarks []string
	var updateClause []string
	for i := 0; i<len(fields); i++ {
		questionMarks = append(questionMarks, "?")
		updateClause = append(updateClause, fields[i]+"=?")
	}
	sqlText := `INSERT INTO ` + self.tableName + ` (` + strings.Join(fields, ", ") + `, updated_at, created_at)
		         VALUES (` + strings.Join(questionMarks, ", ") + `, NOW(), NOW())
		         ON DUPLICATE KEY UPDATE ` + strings.Join(updateClause, ", ") + `, updated_at = NOW()`
	insertStatement, insertErr := self.Conn.Prepare(sqlText)
	return &Stmt{stmt:insertStatement}, insertErr
}

func (self *Stmt) Exec(args ...interface{}) (sql.Result, error) {
	result, err := self.stmt.Exec(append(args, args...)...)
	return result, err
}

func (self *Stmt) Close() {
	self.stmt.Close()
}
