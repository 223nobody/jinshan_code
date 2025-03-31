package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Student struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

func saveStudents(filename string, students []Student) {
	// 使用os.Create会清空原有内容
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("无法创建文件:", err)
		return
	}
	defer file.Close()

	// 使用JSON编码器
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(students); err != nil {
		fmt.Println("保存数据失败:", err)
	}
}

func loadStudents(filename string) []Student {
	var students []Student

	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&students); err != nil {
		return nil
	}

	return students
}

func main() {
	var choice, age int
	var name, gender string
	filename := "student.txt"
	students := loadStudents(filename)
Loop:
	for {
		fmt.Println("\n欢迎来到学生信息管理系统")
		fmt.Println("请选择操作类型")
		fmt.Println("1.录入学生信息")
		fmt.Println("2.查询学生信息")
		fmt.Println("3.修改学生信息")
		fmt.Println("4.删除学生信息")
		fmt.Println("5.退出学生信息管理系统")
		fmt.Scanln(&choice)
		fmt.Println()
		switch choice {
		case 1:
			fmt.Print("请输入将要录入的学生姓名：")
			fmt.Scanln(&name)
			fmt.Print("请输入将要录入的学生性别：")
			fmt.Scanln(&gender)
			fmt.Print("请输入将要录入的学生年龄：")
			fmt.Scanln(&age)
			student := Student{
				Name:   name,
				Gender: gender,
				Age:    age,
			}
			students = append(students, student)
		case 2:
			for key, student := range students {
				fmt.Print(key+1, " 学生姓名：", student.Name, " 学生性别：", student.Gender, " 学生年龄：", student.Age, "\n")
			}
		case 3:
			fmt.Print("请输入将要修改的学生姓名：")
			fmt.Scanln(&name)
			for key, student := range students {
				if student.Name == name {
					fmt.Print(key+1, " 学生姓名：", student.Name, " 学生性别：", student.Gender, " 学生年龄：", student.Age, "\n")
					fmt.Print("请输入学生性别：")
					fmt.Scanln(&gender)
					fmt.Print("请输入学生年龄：")
					fmt.Scanln(&age)
					students[key].Gender = gender
					students[key].Age = age
				}
			}

		case 4:
			fmt.Print("请输入将要删除的学生姓名：")
			fmt.Scanln(&name)
			for key, student := range students {
				if student.Name == name {
					students = append(students[:key], students[key+1:]...)
				}
			}
		case 5:
			saveStudents(filename, students)
			break Loop
		default:
			fmt.Println("输入的操作数不存在，请重新输入")
			continue
		}

	}
}
