package main

import (
	"fmt"

	"github.com/hnakamur/commango/jsonutil"
	"github.com/hnakamur/commango/modules/redhat/service"
	"github.com/hnakamur/commango/modules/redhat/yum"
)

func main() {
	result, err := yum.Installed("ntp")
	if err != nil {
		panic(err)
	}
	json, err := jsonutil.Encode(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(json)

	if result.Rc == 1 {
		result, err = yum.Install("ntp")
		if err != nil {
			panic(err)
		}
		json, err := jsonutil.Encode(result)
		if err != nil {
			panic(err)
		}
		fmt.Println(json)
	}

	result, err = service.Status("ntpd")
	json, err = jsonutil.Encode(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(json)
	if result.Rc == 3 {
		result, err = service.Start("ntpd")
		if err != nil {
			panic(err)
		}
		json, err := jsonutil.Encode(result)
		if err != nil {
			panic(err)
		}
		fmt.Println(json)
	}

	result, err = service.AutoStartEnabled("ntpd")
	if err != nil {
		panic(err)
	}
	json, err = jsonutil.Encode(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(json)
	if result.Rc == 1 {
		result, err = service.EnableAutoStart("ntpd")
		if err != nil {
			panic(err)
		}
		json, err = jsonutil.Encode(result)
		if err != nil {
			panic(err)
		}
		fmt.Println(json)
	}
}
