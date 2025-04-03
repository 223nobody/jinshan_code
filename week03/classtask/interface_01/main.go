package main

import (
	"fmt"
	"math"
)

// Shape 接口定义
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Circle 结构体
type Circle struct {
	Radius float64
}

// Rectangle 结构体
type Rectangle struct {
	Width  float64
	Height float64
}

// Triangle 结构体
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// Circle 实现 Shape 接口的 Area 方法
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Circle 实现 Shape 接口的 Perimeter 方法
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Rectangle 实现 Shape 接口的 Area 方法
func (r Rectangle) Area() float64 {
	return r.Height * r.Width
}

// Rectangle 实现 Shape 接口的 Perimeter 方法
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Height + r.Width)
}

// Triangle 实现 Shape 接口的 Area 方法
func (t Triangle) Area() float64 {
	s := t.Perimeter() / 2
	area := math.Sqrt(s * (s - t.SideA) * (s - t.SideB) * (s - t.SideC)) // 海伦公式
	return area
}

// Triangle 实现 Shape 接口的 Perimeter 方法
func (t Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

func main() {
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 4, Height: 6},
		Triangle{SideA: 3, SideB: 4, SideC: 5},
	}

	for _, shape := range shapes {
		fmt.Printf("图形面积: %.2f, 图形周长: %.2f\n", shape.Area(), shape.Perimeter())
	}
}
