package orm

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/ziutek/mymysql/godrv"
	"reflect"
	"strings"
)

type Model struct {
	tableName         string
	memoizedFields    []string
	noPkFields        []string
	structure         interface{}
	Conn              *sql.DB
	IncludesUpdatedAt bool
	IncludesCreatedAt bool
}

var connectionString string

func SetConnectionString(auth string) {
	connectionString = auth
}

func NewModel(assign interface{}) *Model {
	var f []string
	var noPk []string
	db, dbErr := sql.Open("mymysql", connectionString)
	if dbErr != nil {
		panic(dbErr)
	}

	structValue := reflect.ValueOf(assign)
	modelStructType := structValue.Type()

	tableName, _ := modelStructType.FieldByName("TableName")
	var i int
	for i = 0; i < modelStructType.NumField(); i += 1 {
		if modelStructType.Field(i).Name != "TableName" {
			names := strings.Split(fmt.Sprintf("%v.%v", tableName.Tag, modelStructType.Field(i).Tag), ":")
			if len(names) > 1 && (names[1] == "pk" || names[1] == "PK" || names[1] == "Pk") {
				f = append(f, names[0])
			} else {
				noPk = append(noPk, names[0])
				f = append(f, fmt.Sprintf("%v.%v", tableName.Tag, modelStructType.Field(i).Tag))
			}
		}
	}

	m := Model{tableName: fmt.Sprintf("%v", tableName.Tag), Conn: db, memoizedFields: f, noPkFields: noPk, IncludesUpdatedAt: true, IncludesCreatedAt: true}
	return &m
}

type Fieldable interface {
	Fields() []interface{}
	FieldsNoPk() []interface{}
	SetPk(int64)
}

func infoLog(message string) {
	debug := true
	debug = false
	if debug == true {
		println(message)
	}
}

type H map[string]string

func (self *Model) Create(keyValues H) error {
	var values []interface{}
	var keys []string
	var marks []string
	for k, v := range keyValues {
		marks = append(marks, "?")
		keys = append(keys, k)
		values = append(values, v)
	}

	sqlText := `INSERT IGNORE ` + self.tableName + `(` + strings.Join(keys, ", ") + `) VALUES (` + strings.Join(marks, ", ") + `);`

	insertStatement, insertErr := self.Conn.Prepare(sqlText)
	if insertErr != nil {
		return insertErr
	}
	defer insertStatement.Close()

	_, insertExecErr := insertStatement.Exec(values...)
	if insertExecErr != nil {
		return insertExecErr
	}
	return nil
}

func (self *Model) Truncate() {
	sql := "TRUNCATE TABLE " + self.tableName + ";"
	infoLog(sql)
	trunPlacements, trunPlacementsErr := self.Conn.Prepare(sql)
	if trunPlacementsErr != nil {
		panic(trunPlacementsErr)
	}
	defer trunPlacements.Close()

	_, insertErr := trunPlacements.Exec()
	if insertErr != nil {
		panic(insertErr)
	}
}

func (self *Model) Count() int {
	sql := "SELECT count(*) from " + self.tableName + ";"
	infoLog(sql)
	row := self.Conn.QueryRow(sql)

	var theCount int
	err := row.Scan(&theCount)
	if err != nil {
		panic(err)
	}
	return theCount
}

func (self *Model) CountWhere(params H) int {
	var whereClause []string
	for k, v := range params {
		whereClause = append(whereClause, k+" = "+v)
	}
	sql := "SELECT count(id) FROM " + self.tableName + " WHERE " + strings.Join(whereClause, " AND ") + ";"
	infoLog(sql)
	row := self.Conn.QueryRow(sql)

	var theCount int
	err := row.Scan(&theCount)
	if err != nil {
		panic(err)
	}
	return theCount
}

type Query struct {
	SelectClause string
	FromClause   string
	JoinClause   string
	WhereClause  string
	LimitClause  string
	model        *Model
	Params       []interface{}
}

func (self *Model) First(object Fieldable) error {
	query := new(Query)
	query.model = self

	query.SelectClause = "SELECT " + strings.Join(self.memoizedFields, ", ")
	query.FromClause = " FROM " + self.tableName + " "
	query.LimitClause = `LIMIT 1`

	err := query.Find(object)
	return err
}

func (self *Model) Save(object Fieldable) error {
	stmt, err := self.PrepareInsert(self.memoizedFields)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, errs := stmt.Exec(object.Fields()...)
	if errs != nil {
		return errs
	}
	pk, pkErr := result.LastInsertId()
	if pkErr != nil {
		return pkErr
	}
	object.SetPk(pk)
	return nil
}

func (self *Model) Insert(object Fieldable) error {
	stmt, err := self.PrepareInsert(self.noPkFields)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, errs := stmt.Exec(object.FieldsNoPk()...)
	if errs != nil {
		return errs
	}
	pk, pkErr := result.LastInsertId()
	if pkErr != nil {
		return pkErr
	}
	object.SetPk(pk)
	return nil
}

func (self *Model) FindOrCreate(conditions H, obj Fieldable) {
	var values []interface{}
	var marks []string
	for k, v := range conditions {
		marks = append(marks, k+"=?")
		values = append(values, v)
	}

	err := self.Where(strings.Join(marks, ", "), values...).Find(obj)
	if err == nil {
		return
	}

	self.Create(conditions)
}

func (self *Model) All() *Query {
	query := new(Query)
	query.model = self
	// query.Params = params

	query.SelectClause = "SELECT " + strings.Join(self.memoizedFields, ", ")
	query.FromClause = " FROM " + self.tableName + " "
	// query.WhereClause = " WHERE " + where

	return query
}

func (self *Model) Where(where string, params ...interface{}) *Query {
	query := new(Query)
	query.model = self
	query.Params = params

	query.SelectClause = "SELECT " + strings.Join(self.memoizedFields, ", ")
	query.FromClause = " FROM " + self.tableName + " "
	query.WhereClause = " WHERE " + where

	return query
}

func (self *Model) From(from string) *Query {
	query := new(Query)
	query.model = self

	query.SelectClause = "SELECT " + strings.Join(self.memoizedFields, ", ")
	query.FromClause = " FROM " + from + " "
	query.WhereClause = " "

	return query
}

func (self *Model) Join(join string) *Query {
	query := new(Query)
	query.model = self

	query.SelectClause = "SELECT " + strings.Join(self.memoizedFields, ", ")
	query.FromClause = " FROM " + self.tableName + " "
	query.JoinClause = " JOIN " + join
	query.WhereClause = " "

	return query
}

func (query *Query) From(from string) *Query {
	query.JoinClause = " FROM " + from
	return query
}

func (query *Query) Join(join string) *Query {
	query.JoinClause = " JOIN " + join
	return query
}

func (query *Query) Where(where string, params ...interface{}) *Query {
	query.Params = params
	query.WhereClause = " WHERE " + where
	return query
}

func (query *Query) Find(object Fieldable) error {
	sql := query.SelectClause + query.FromClause + query.JoinClause + query.WhereClause + query.LimitClause + ";"
	infoLog(sql)

	stmt, err := query.model.Conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRow(query.Params...)

	selectError := row.Scan(object.Fields()...)
	if selectError != nil {
		return selectError
	}

	return nil
}

func (query *Query) FindAll(fieldables interface{}) error {
	sqlText := query.SelectClause + query.FromClause + query.JoinClause + query.WhereClause + query.LimitClause + ";"
	infoLog(sqlText)

	stmt, err := query.model.Conn.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, errs := stmt.Query(query.Params...)
	if errs != nil {
		return errs
	}

	sliceValue := reflect.Indirect(reflect.ValueOf(fieldables))
	if sliceValue.Kind() != reflect.Slice {
		return errors.New("needs a pointer to a slice")
	}
	sliceElementType := sliceValue.Type().Elem()

	for rows.Next() {
		newValue := reflect.New(sliceElementType)
		fields := newValue.MethodByName("Fields").Call(nil)

		var Values []interface{}
		refValue := reflect.Indirect(reflect.ValueOf(&Values))
		refValue.Set(reflect.AppendSlice(refValue, fields[0]))
		rows.Scan(Values...)

		sliceValue.Set(reflect.Append(sliceValue, reflect.Indirect(reflect.ValueOf(newValue.Interface()))))
	}
	return nil
}
