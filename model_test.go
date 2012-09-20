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

func TestNewModelNoFields(*testing.T) {
	SetConnectionString("central_test/root/")
	ManualStats = NewModel("manual_stats", 
		func() Fieldable {return new(ManualStat)})
		
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
}