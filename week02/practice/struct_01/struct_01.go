package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name  string
	Age   int
	Email string
}

func main() {
	p := NewPerson("付坤", 20, "2141024586@qq.com")
	p.PrintPerson()

	jsonData, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("\n" + string(jsonData))

	jsonData1, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("\n" + string(jsonData1))
}

func NewPerson(name string, age int, email string) Person {
	p := Person{
		Name:  name,
		Age:   age,
		Email: email,
	}
	return p
}

func (p Person) PrintPerson() {
	fmt.Println("姓名是：", p.Name)
	fmt.Println("年龄是：", p.Age)
	fmt.Println("邮箱是：", p.Email)
}
