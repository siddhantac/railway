package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

func parseJSON(file string) ([]interface{}, error) {
	d, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err = json.Unmarshal(d, &m); err != nil {
		return nil, err
	}

	stn, ok := m["stations"]
	if !ok {
		return nil, fmt.Errorf("'stations' key not found")
	}

	stns, ok := stn.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'stations' is of wrong type, expected an array")
	}

	//_, ok := stns[0].(string)
	//if !ok {
	////fmt.Printf("%T", strStns)
	//}

	return stns, nil
}

func main() {
	var config string
	var up, down bool
	var trains int

	flag.StringVar(&config, "d", "", "data filename (json)")
	flag.IntVar(&trains, "n", 1, "number of trains to run")
	flag.BoolVar(&up, "up", false, "run trains only in UP direction")
	flag.BoolVar(&down, "down", false, "run trains only in DOWN direction")
	flag.Parse()

	stnNames, err := parseJSON(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	line := GreenLine
	stns := make([]*Station, 0)

	for _, n := range stnNames {
		name := n.(string)
		s := NewStation(name)
		s.AddLine(line)
		stns = append(stns, s)
	}

	Connect(line, stns...)

	for _, stn := range stns {
		stn.Run()
	}

	rand.Seed(time.Now().UnixNano())

	if up {
		stns[0].StartService(line, trains)
	}

	if down {
		stns[len(stns)-1].StartService(line, trains)
	}

	time.Sleep(time.Second * 20)
}