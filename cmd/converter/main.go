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

// TODO generate info.plist from json

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
	Decimals string
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

func (cs *coins) convert(fromUnitName string, value string) ([]*conversionResult, error) {
	results := make([]*conversionResult, 0)
	for _, c := range cs.filterByUnitName(fromUnitName) {
		res, err := c.convert(fromUnitName, value)
		if err != nil {
			return nil, err
		}
		results = append(results, res...)
	}
	return results, nil
}

func (c *coin) convert(fromUnitName string, value string) ([]*conversionResult, error) {
	fromUnit := c.findUnitByName(fromUnitName)
	if fromUnit == nil {
		return nil, nil
	}

	results := make([]*conversionResult, 0)
	for _, unit := range c.Units {
		if unit.Name == fromUnitName {
			continue
		}
		converted, err := convert(fromUnit, unit, value)
		if err != nil {
			return nil, err
		}
		results = append(results, &conversionResult{
			unitName: unit.Name,
			icon:     c.Icon,
			value:    converted,
		})
	}
	return results, nil
}

func convert(fromUnit, toUnit *unit, value string) (string, error) {
	val, err := decimal.NewFromString(value)
	if err != nil {
		return "", err
	}
	fromDecimal, err := decimal.NewFromString(fromUnit.Decimals)
	if err != nil {
		return "", err
	}
	toDecimal, err := decimal.NewFromString(toUnit.Decimals)
	if err != nil {
		return "", err
	}
	propotion := fromDecimal.Div(toDecimal)
	return val.Mul(propotion).String(), nil
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
	fromUnitName := keyword

	log.Println("fromUnitName:", fromUnitName)
	results, err := coins.convert(fromUnitName, queryValue)
	if err != nil {
		wf.FatalError(err)
	}

	for _, result := range results {
		wf.NewItem(result.title()).
			Icon(&aw.Icon{
				Value: result.icon,
				Type:  aw.IconTypeImage,
			}).
			Title(result.title()).
			Subtitle(result.unitName).
			UID(result.icon + result.unitName).
			Valid(true).
			Arg(result.value)
	}

	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
