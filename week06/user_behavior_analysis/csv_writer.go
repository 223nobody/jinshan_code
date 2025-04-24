package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

func WriteUserStatsCSV(stats map[string]UserStat, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	if err := writer.Write([]string{"用户ID", "操作次数", "首次操作时间", "最后操作时间"}); err != nil {
		return err
	}

	// 排序用户ID
	var users []string
	for u := range stats {
		users = append(users, u)
	}
	sort.Strings(users)

	// 写入数据
	for _, userID := range users {
		stat := stats[userID]
		record := []string{
			userID,
			fmt.Sprintf("%d", stat.Count),
			stat.First.Format(time.RFC3339),
			stat.Last.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func WriteMinuteStatsCSV(minuteStats map[time.Time]MinuteStat, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头（严格匹配图片字段顺序）
	header := []string{"时间段", "活跃用户数", "操作总数"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("写入表头失败: %w", err)
	}

	// 将时间窗口转换为切片并排序
	timeWindows := make([]time.Time, 0, len(minuteStats))
	for t := range minuteStats {
		timeWindows = append(timeWindows, t)
	}
	sort.Slice(timeWindows, func(i, j int) bool {
		return timeWindows[i].Before(timeWindows[j])
	})

	// 时间格式化为 "YYYY-MM-DD HH:MM"（与图片示例一致）
	const timeFormat = "2006-01-02 15:04"

	// 写入数据行
	for _, tw := range timeWindows {
		stat := minuteStats[tw]
		record := []string{
			tw.Format(timeFormat),
			strconv.Itoa(len(stat.ActiveUsers)),
			strconv.Itoa(stat.TotalActions),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("写入数据行失败: %w", err)
		}
	}

	return nil
}
func WriteActionStatsCSV(actionStats map[string]int, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头（完全匹配图片格式）
	header := []string{"行为类型", "总次数"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("写入表头失败: %w", err)
	}

	// 对行为类型进行字母排序（与图片顺序一致）
	actions := make([]string, 0, len(actionStats))
	for action := range actionStats {
		actions = append(actions, action)
	}
	sort.Strings(actions)

	// 写入数据行
	for _, action := range actions {
		record := []string{
			action,
			strconv.Itoa(actionStats[action]),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("写入数据行失败: %w", err)
		}
	}

	return nil
}
