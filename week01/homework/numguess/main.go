package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func write(filename string, content string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content + "\n")
	if err != nil {
		fmt.Println("写入文件失败:", err)
	}
}

func main() {
	var choice int = 0
	var choice2 int = 0
	var chance int = 0
	var num int = 0
	j := 1
	filename := "../../homework/game.txt"
	os.WriteFile(filename, []byte(""), 0644)
Loop1:
	for {
		fmt.Printf("\n\t第%d次游戏开始", j)
		write(filename, fmt.Sprintf("\t第%d次游戏开始", j))
		j++
		fmt.Println("\n欢迎来到猜数字游戏！")
		fmt.Println("1.计算机将在1到100之间随机选择一个数字。")
		fmt.Println("2.您可以选择难度级别(简单、中等、困难)，不同的难度对应不同的猜测机会。")
		fmt.Println("3.请输入您的猜测")
		fmt.Println()
		fmt.Println("请选择难度级别（简单/中等/困难）：")
		fmt.Println("1.简单( 10 次机会)")
		fmt.Println("2.中等( 5 次机会)")
		fmt.Println("3.困难( 3 次机会)")
		fmt.Print("\n请输入选择：")
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			chance = 10
			write(filename, "用户选择简单模式")
		case 2:
			chance = 5
			write(filename, "用户选择中等模式")
		case 3:
			chance = 3
			write(filename, "用户选择困难模式")
		default:
			fmt.Println("输入错误，请重新选择难度")
			write(filename, "输入错误，请重新选择难度\n")
			continue
		}
		fmt.Println("开始游戏：")
		start := time.Now()
		write(filename, fmt.Sprint("游戏开始时间"+start.String()))
		// rand.Seed(time.Now().UnixNano())
		randnum := rand.Intn(100) + 1
		write(filename, fmt.Sprintf("系统生成随机数 %d ", randnum))
	Loop:
		for i := 1; i <= chance; i++ {
			fmt.Printf("\n第 %d 次猜测，请输入您的数字(1 - 100)：", i)
			_, err := fmt.Scanln(&num)
			write(filename, fmt.Sprintf("第 %d 次猜测，用户输入 %d ", i, num))
			if err != nil || num < 1 || num > 100 {
				fmt.Println("输入错误，请输入1到100之间的整数！")
				write(filename, "输入错误，请输入1到100之间的整数！")
				i--
				continue
			}
			switch {
			case num < randnum:
				fmt.Print("您猜的数字小了")
				write(filename, "您猜的数字小了")
			case num > randnum:
				fmt.Print("您猜的数字大了")
				write(filename, "您猜的数字大了")
			default:
				fmt.Printf("恭喜您猜对了！您在第 %d 次猜测中成功", i)
				write(filename, fmt.Sprintf("恭喜您猜对了！您在第 %d 次猜测中成功", i))
				break Loop
			}
		}
		fmt.Println("\n本次游戏结束")
		write(filename, "本次游戏结束")
		fmt.Printf("本次游戏用时 %v\n", time.Since(start.Round(time.Second)))
		write(filename, fmt.Sprintf("本次游戏用时 %v\n", time.Since(start.Round(time.Second))))
		fmt.Println("\n\n请选择：")
		fmt.Println("1.重新开始")
		fmt.Println("2.退出游戏")
		for {
			fmt.Scanln(&choice2)
			if choice2 == 1 {
				break
			}
			if choice2 == 2 {
				fmt.Println("用户退出游戏")
				break Loop1
			} else if choice != 1 {
				fmt.Println("输入错误，请重新输入")
			}
		}

	}
}
