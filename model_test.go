package orm

import (
	// "fmt"
	. "github.com/KarateCode/helpers"
	"testing"
	"time"
)

const dateLayout = "2006-01-02"

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

var ManualStats *Model

type ManualStat struct {
	TableName     string    `manual_stats`
	Id            int64     `id:pk`
	ForDate       time.Time `for_date`
	ClientImps    float64   `imps`
	ClientClicks  float64   `clicks`
	ClientConvs   float64   `conversions`
	ClientRevenue float64   `revenue`
}

func (self *ManualStat) Fields() []interface{} {
	var results []interface{}
	return append(results, &self.Id, &self.ForDate, &self.ClientImps, &self.ClientClicks, &self.ClientConvs, &self.ClientRevenue)
}
func (self *ManualStat) FieldsNoPk() []interface{} {
	var results []interface{}
	return append(results, &self.ForDate, &self.ClientImps, &self.ClientClicks, &self.ClientConvs, &self.ClientRevenue)
}
func (self *ManualStat) SetPk(pk int64) {
	self.Id = pk
}

// func (self *[]ManualStat) FieldsAt(index int) []interface{} {
// var results []interface{}
// return append(results, &self.Id, &self.ForDate, &self.ClientImps, &self.ClientClicks, &self.ClientConvs, &self.ClientRevenue)
// return 
// }

func TestNewModelNoFields(*testing.T) {
	SetConnectionString("central_test/root/")
	ManualStats = NewModel(ManualStat{})

	ManualStats.Truncate()
	startDate, err := time.Parse(dateLayout, "2012-09-01")
	checkError(err)

	mstat := ManualStat{
		ForDate:       startDate,
		ClientRevenue: 3.00,
	}
	err = ManualStats.Save(&mstat)
	checkError(err)
	// fmt.Printf("\nmstat%+v\n", mstat)
	ShouldEqual(1, ManualStats.Count())
	stat := ManualStat{}
	ManualStats.Where("id=?", 1).Find(&stat)
	ShouldEqual(3.0, stat.ClientRevenue)
}

func TestFindAll(*testing.T) {
	SetConnectionString("central_test/root/")
	ManualStats = NewModel(ManualStat{})

	ManualStats.Truncate()
	startDate, err := time.Parse(dateLayout, "2012-09-01")
	checkError(err)

	mstat := ManualStat{
		ForDate:       startDate,
		ClientRevenue: 3.00,
	}
	err = ManualStats.Save(&mstat)
	checkError(err)
	mstat2 := ManualStat{
		ForDate:       startDate,
		ClientRevenue: 3.00,
	}
	err = ManualStats.Save(&mstat2)
	checkError(err)
	ShouldEqual(2, ManualStats.Count())

	var stats []ManualStat
	ManualStats.All().FindAll(&stats)
	ShouldEqual(2, len(stats))
	ShouldEqual(3.0, stats[0].ClientRevenue)
	// fmt.Printf("\nstats%+v\n", stats)
	// ManualStats.All().FindAll()
}

var Advertisers *Model

type Advertiser struct {
	TableName   string    `advertisers`
	Id          int64     `id:pk` //if the table's PrimaryKey is not id ,should add `PK` to ident
	Name        string    `name`
	Type        string    `type`
	ApiId       int64     `api_id`
	RefreshedAt time.Time `refreshed_at`
	// Tolerance   sql.NullFloat64 `tolerance`
	EntityType string `entity_type`
	Active     bool   `active`
}

func (self *Advertiser) Fields() []interface{} {
	var results []interface{}
	return append(results, &self.Id, &self.Name, &self.Type, &self.ApiId, &self.RefreshedAt, &self.EntityType, &self.Active)
	// if using 'pk' in the tag, don't include it here:
	// return append(results, &self.Name, &self.Type, &self.ApiId, &self.RefreshedAt, &self.EntityType, &self.Active)
}
func (self *Advertiser) FieldsNoPk() []interface{} {
	var results []interface{}
	return append(results, &self.Name, &self.Type, &self.ApiId, &self.RefreshedAt, &self.EntityType, &self.Active)
	// if using 'pk' in the tag, don't include it here:
	// return append(results, &self.Name, &self.Type, &self.ApiId, &self.RefreshedAt, &self.EntityType, &self.Active)
}
func (self *Advertiser) SetPk(pk int64) {
	self.Id = pk
}
func TestDoNotSetIdToZero(*testing.T) {
	SetConnectionString("central_test/root/")
	Advertisers = NewModel(Advertiser{})
	Advertisers.Truncate()

	a := Advertiser{
		ApiId:  62018,
		Name:   "OKC Chickasaw",
		Type:   "ApxAdvertiser",
		Active: true,
	}
	var id int64 = 1

	checkError(Advertisers.Save(&a))
	ShouldEqual(id, a.Id)
	// fmt.Printf("a: %+v\n", a)

	a = Advertiser{
		ApiId:  62018,
		Name:   "OKC Chickasaw",
		Type:   "ApxAdvertiser",
		Active: true,
	}
	checkError(Advertisers.Save(&a))
	checkError(Advertisers.First(&a))
	ShouldEqual(id, a.Id)
	// fmt.Printf("a: %+v\n", a)
}
