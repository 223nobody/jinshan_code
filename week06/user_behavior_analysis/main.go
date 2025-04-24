package main

import (
	"flag"
	"log"
)

func main() {
	inputFile := flag.String("input", "user_actions.log", "输入文件的路径")
	outputPrefix := flag.String("output", "", "输出文件的前缀")
	flag.Parse()

	// 1. 解析日志
	records, _ := ParseLogFile(*inputFile)

	// 2. 生成统计
	stats := GenerateStats(records)

	// 3. 写入CSV
	if err := WriteUserStatsCSV(stats.UserStats, *outputPrefix+"user_statistics.csv"); err != nil {
		log.Fatalf("用户统计写入失败: %v", err)
	}
	if err := WriteActionStatsCSV(stats.ActionStats, *outputPrefix+"action_statistics.csv"); err != nil {
		log.Fatalf("行为统计写入失败: %v", err)
	}
	if err := WriteMinuteStatsCSV(stats.MinuteStats, *outputPrefix+"minute_statistics.csv"); err != nil {
		log.Fatalf("分钟统计写入失败: %v", err)
	}

	log.Println("处理完成")
}
