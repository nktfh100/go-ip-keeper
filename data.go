package main

import (
	"encoding/json"
	"os"
)

type DataIPs struct {
	IP  string   `json:"ip"`
	IPs []DataIP `json:"ips"`
}

type DataIP struct {
	IP    string `json:"ip"`
	Owner string `json:"owner"`
}

func LoadData() DataIPs {
	fileData, err := os.ReadFile("./ip-keeper-data.json")
	if err != nil {
		if os.IsNotExist(err) {
			return DataIPs{}
		}
		panic(err)
	}

	var data DataIPs

	err = json.Unmarshal(fileData, &data)
	if err != nil {
		panic(err)
	}
	return data
}

func SaveData(data DataIPs) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("./ip-keeper-data.json", jsonData, 0644)
	if err != nil {
		panic(err)
	}
}

func (data *DataIPs) removeIPByI(i int) {
	newIPs := append(data.IPs[:i], data.IPs[i+1:]...)
	data.IPs = newIPs
}

func (data *DataIPs) ipExists(ip string) bool {
	for _, ele := range data.IPs {
		if ele.IP == ip {
			return true
		}
	}
	return false
}
