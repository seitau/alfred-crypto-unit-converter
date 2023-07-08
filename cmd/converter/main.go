package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	aw "github.com/deanishe/awgo"
	"github.com/shopspring/decimal"

	"golang.org/x/exp/slices"
)

var (
	decimalRegexp = regexp.MustCompile(`^\d*\.?\d*$`)
	logger        = log.New(os.Stderr, "converter", log.LstdFlags)
	wf            *aw.Workflow

	query string
)

// TODO jsonのunit名からinfo.plistを生成するようにする
// goのtemplateでbuild以下に配置する
// copy to clipboard
// send notiifcation

func init() {
	flag.StringVar(&query, "query", "", "search query")

	wf = aw.New()
}

type coin struct {
	Name  string
	Icon  string
	Units []*unit
}

type coins []*coin

type unit struct {
	Name     string
	Decimals int32
}

func (c *coin) findUnitByName(unitName string) *unit {
	for i := range c.Units {
		if u := c.Units[i]; u.Name == unitName {
			return u
		}
	}
	return nil
}

func (cs coins) filterByUnitName(unitName string) coins {
	filtered := make(coins, 0)
	for i := range cs {
		hasUnit := slices.ContainsFunc(cs[i].Units, func(u *unit) bool {
			return u.Name == unitName
		})
		if hasUnit {
			filtered = append(filtered, cs[i])
		}
	}
	return filtered
}

type conversionResult struct {
	unitName string
	icon     string
	value    string
}

func (c *conversionResult) title() string {
	return fmt.Sprintf("%s %s", c.value, c.unitName)
}

func convert(coins coins, baseUnitName string, queryValue string) ([]*conversionResult, error) {
	cs := coins.filterByUnitName(baseUnitName)
	if len(cs) == 0 {
		return nil, nil
	}
	results := make([]*conversionResult, 0)
	for i := range cs {
		baseUnit := coins[i].findUnitByName(baseUnitName)
		if baseUnit == nil {
			continue
		}
		for _, unit := range coins[i].Units {
			if unit.Name == baseUnitName {
				continue
			}
			diff := baseUnit.Decimals - unit.Decimals
			m := decimal.New(1, diff)
			d, err := decimal.NewFromString(queryValue)
			if err != nil {
				return nil, err
			}
			results = append(results, &conversionResult{
				unitName: unit.Name,
				icon:     coins[i].Icon,
				value:    m.Mul(d).String(),
			})
		}
	}
	return results, nil
}

func readCoins() (coins, error) {
	f, err := os.Open("./coins.json")
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var coins []*coin
	if err := json.Unmarshal(b, &coins); err != nil {
		return nil, err
	}
	return coins, nil
}

func run() {
	wf.Args()
	flag.Parse()

	coins, err := readCoins()
	if err != nil {
		wf.FatalError(err)
	}
	log.Println(coins)

	args := wf.Args()
	if len(args) == 0 {
		return
	}
	queryValue := args[0]
	if !decimalRegexp.MatchString(queryValue) {
		wf.FatalError(fmt.Errorf("invalid decimal value"))
	}
	keyword, ok := wf.Config.Env.Lookup("alfred_workflow_keyword")
	if !ok {
		wf.FatalError(fmt.Errorf("keyword not found"))
	}
	baseUnitName := keyword

	results, err := convert(coins, baseUnitName, queryValue)
	if err != nil {
		wf.FatalError(err)
	}

	for _, result := range results {
		item := wf.NewItem(result.title())
		item.Icon(&aw.Icon{
			Value: result.icon,
			Type:  aw.IconTypeImage,
		})
		item.Title(result.title())
		item.Subtitle(result.unitName)
	}

	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
