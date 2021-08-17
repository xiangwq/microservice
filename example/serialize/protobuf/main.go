package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"microservice/example/serialize/protobuf/test"
)

func main() {
	fmt.Println(1111)
	var person test.Person
	person.Id = 111
	person.Name = "test"

	var phone test.Phone
	phone.Number = "13600000000"
	person.Phones = append(person.Phones, &phone)

	data, err := proto.Marshal(&person)
	if err != nil {
		fmt.Println("marshal failed")
		return
	}

	ioutil.WriteFile("./test.dat", data, 0777)

	var person2 test.Person

	data2, err := ioutil.ReadFile("./test.dat")
	if err != nil {
		fmt.Println("read failed")
		return
	}

	proto.Unmarshal(data2, &person2)
	fmt.Printf("person2 %v", person2)
}
