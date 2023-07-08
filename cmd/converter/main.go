package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	aw "github.com/deanishe/awgo"
)

var (
	logger = log.New(os.Stderr, "converter", log.LstdFlags)
	wf     *aw.Workflow

	query string
)

// TODO jsonのunit名からinfo.plistを生成するようにする
// copy to clipboard
// send notiifcation

func init() {
	flag.StringVar(&query, "query", "", "search query")

	wf = aw.New()
}

type coin struct {
	Name  string
	Icon  string
	Units []unit
}

type unit struct {
	Name     string
	Decimals uint64
}

func readUnits() ([]coin, error) {
	f, err := os.Open("./coins.json")
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var coins []coin
	if err := json.Unmarshal(b, &coins); err != nil {
		return nil, err
	}
	return coins, nil
}

func run() {
	wf.Args()
	flag.Parse()

	coins, err := readUnits()
	if err != nil {
		wf.FatalError(err)
	}
	log.Println(coins)

	args := wf.Args()
	if len(args) == 0 {
		return
	}
	keyword, ok := wf.Config.Env.Lookup("alfred_workflow_keyword")
	if !ok {
		wf.FatalError(fmt.Errorf("keyword not found"))
	}
	baseUnit := keyword
	log.Println(baseUnit)

	// TODO jsonを読み取ってbase unitからunitsを取得し、baseunit以外に変換したものを表示する

	item := wf.NewItem(args[0])
	item.Icon(&aw.Icon{
		Value: "icons/ethereum.png",
		Type:  aw.IconTypeImage,
	})
	item.Title("0.0000001 ether")
	item.Subtitle("ether")
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
