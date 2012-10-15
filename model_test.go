package orm

import (
	// "fmt"
	"testing"
	. "github.com/KarateCode/helpers"
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
	TableName string `manual_stats`
	Id int64 `id`
	ForDate time.Time `for_date`
	ClientImps float64 `client_imps`
	ClientClicks float64 `client_clicks`
	ClientConvs float64 `client_convs`
	ClientRevenue float64 `client_revenue`
}

func (self *ManualStat) Fields() []interface{} {
	var results []interface{}
	return append(results, &self.Id, &self.ForDate, &self.ClientImps, &self.ClientClicks, &self.ClientConvs, &self.ClientRevenue)
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
		ForDate: startDate,
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
		ForDate: startDate,
		ClientRevenue: 3.00,
	}
	err = ManualStats.Save(&mstat)
	checkError(err)
	mstat2 := ManualStat{
		ForDate: startDate,
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
