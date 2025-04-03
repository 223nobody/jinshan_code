package main

import "fmt"

// Animal 接口定义
type Animal interface {
	Speak() string
	Move() string
	Name() string
}

// Dog 结构体
type Dog struct {
	name string
}

// Cat 结构体
type Cat struct {
	name string
}

// Bird 结构体
type Bird struct {
	name string
}

// Dog 实现 Animal 接口的 Speak 方法
func (d Dog) Speak() string {
	return "汪汪汪！"
}

// Dog 实现 Animal 接口的 Move 方法
func (d Dog) Move() string {
	return "用四条腿跑"
}

// Dog 实现 Name 方法
func (d Dog) Name() string {
	return d.name
}

// Cat 实现 Animal 接口的 Speak 方法
func (c Cat) Speak() string {
	return "喵喵喵！"
}

// Cat 实现 Animal 接口的 Move 方法
func (c Cat) Move() string {
	return "优雅地走猫步"
}

// Cat 实现 Name 方法
func (c Cat) Name() string {
	return c.name
}

// Bird 实现 Animal 接口的 Speak 方法
func (b Bird) Speak() string {
	return "叽叽喳喳！"
}

// Bird 实现 Animal 接口的 Move 方法
func (b Bird) Move() string {
	return "拍打翅膀飞"
}

// Bird 实现 Name 方法
func (b Bird) Name() string {
	return b.name
}

func main() {
	animals := []Animal{
		Dog{name: "Buddy"},
		Cat{name: "Whiskers"},
		Bird{name: "Tweety"},
	}

	for _, animal := range animals {
		fmt.Printf("%s 说: %s, %s 移动方式: %s\n", animal.(interface{ Name() string }).Name(), animal.Speak(), animal.(interface{ Name() string }).Name(), animal.Move())
	}
}
