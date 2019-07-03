package main

import "fmt"

type IstioRoutingInterface interface {
	TrafficHeaders(headers map[string]string) error
	TrafficPercentage(percentage int32) error
}

type IstioTraffic struct {
	Name string
}

func (t IstioTraffic) TrafficHeaders(headers map[string]string) error {
	fmt.Println("headers:", headers)
	return nil
}

func (t IstioTraffic) TrafficPercentage(percentage int32) error {
	fmt.Println("%:", percentage)
	return nil
}

func main() {
	var obj IstioRoutingInterface
	headers := map[string]string{
		"oi":  "oi2",
		"oi3": "oi",
	}
	custom := IstioTraffic{"Headers"}
	obj = custom
	obj.TrafficHeaders(headers)
	obj.TrafficPercentage(2)
}
