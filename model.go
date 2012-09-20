package orm

import (
	"fmt"
	"strings"
	"database/sql"
	"reflect"
	_ "github.com/ziutek/mymysql/godrv"
)

type Model struct {
	tableName      string
	memoizedFields []string
	structure      interface{}
	Conn           *sql.DB
	NewInstance    func() Fieldable
}

var connectionString string

func SetConnectionString(auth string) {
	connectionString = auth
}

func NewModel(tName string, assign func() Fieldable) *Model {
	var fields []string
	db, dbErr := sql.Open("mymysql", connectionString)
	if dbErr != nil { panic(dbErr) }
	
	structValue := reflect.ValueOf(assign()).Elem()
	modelStructType := structValue.Type()
	// println(modelStructType.Kind().String())
	
	// var results []interface{}
	// return append(results, &self.Id, &self.ForDate, &self.ClientImps, &self.ClientClicks, &self.ClientConvs, &self.ClientRevenue)
	
	var i int
	for i=0; i<modelStructType.NumField(); i+=1 {
		// println(fmt.Sprintf("%v", modelStructType.Field(i).Tag))
		fields = append(fields, fmt.Sprintf("%v", modelStructType.Field(i).Tag))
	}
	m := Model{tableName:tName, Conn:db, NewInstance:assign, memoizedFields:fields}
	// fmt.Printf("\nfields%+v\n", fields)
	return &m
}

type Fieldable interface {
	Fields() []interface{}
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
	for k, v := range(keyValues) {
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
		whereClause = append(whereClause, k + " = " + v)
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
	FromClause string
	WhereClause string
	model *Model
	Params []interface{}
}

func (self *Model) Save(object Fieldable) error {
	stmt, err := self.PrepareInsert(self.memoizedFields)
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	// fmt.Printf("\nobject.Fields()%v\n", object.Fields())
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

func (self *Model) FindOrCreate(conditions H, obj Fieldable)  {
	var values []interface{}
	var marks []string
	for k, v := range(conditions) {
		marks = append(marks, k+"=?")
		values = append(values, v)
	}
	
	err := self.Where(strings.Join(marks, ", "), values...).Find(obj)
	if err == nil {
		return
	}
	
	self.Create(conditions)
}

func (self *Model) All() (*Query) {
	query := new(Query)
	query.model = self
	// query.Params = params
	
	query.SelectClause = "SELECT " + strings.Join(self.memoizedFields, ", ")
	query.FromClause = " FROM " + self.tableName + " "
	// query.WhereClause = " WHERE " + where
	
	return query
}

func (self *Model) Where(where string, params ...interface{}) (*Query) {
	query := new(Query)
	query.model = self
	query.Params = params
	
	query.SelectClause = "SELECT " + strings.Join(self.memoizedFields, ", ")
	query.FromClause = " FROM " + self.tableName + " "
	query.WhereClause = " WHERE " + where
	
	return query
}

func (query *Query) Find(object Fieldable) error {
	sql := query.SelectClause + query.FromClause + query.WhereClause + ";"
	infoLog(sql)
	
	stmt, err := query.model.Conn.Prepare(sql)
	if err != nil {return err}
	defer stmt.Close()
	
	row := stmt.QueryRow(query.Params...)
	
	selectError := row.Scan(object.Fields()...)
	if selectError != nil {return selectError}

	return nil
}

func (query *Query) FindAll() ([]Fieldable, error) {
	sqlText := query.SelectClause + query.FromClause + query.WhereClause + ";"
	infoLog(sqlText)
	
	stmt, err := query.model.Conn.Prepare(sqlText)
	if err != nil {return nil, err}
	defer stmt.Close()
	
	rows, errs := stmt.Query(query.Params...)
	if errs != nil { return nil, errs }
	var fieldables []Fieldable
	
	i := 0
	for rows.Next() {
		ph := query.model.NewInstance()
		fieldables = append(fieldables, ph)
		rows.Scan(fieldables[i].Fields()...)
		i++
	}
	return fieldables, nil
}
