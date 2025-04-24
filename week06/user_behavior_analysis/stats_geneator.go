package main

import (
	"time"
)

type Statistics struct {
	UserStats   map[string]UserStat
	ActionStats map[string]int
	MinuteStats map[time.Time]MinuteStat
}

type UserStat struct {
	Count       int
	First, Last time.Time
}

type MinuteStat struct {
	TimeWindow   time.Time
	ActiveUsers  map[string]struct{}
	TotalActions int
}

func GenerateStats(records []LogRecord) *Statistics {
	stats := &Statistics{
		UserStats:   make(map[string]UserStat),
		ActionStats: make(map[string]int),
		MinuteStats: make(map[time.Time]MinuteStat),
	}

	for _, r := range records {
		// 用户统计
		if user, exists := stats.UserStats[r.UserID]; exists {
			user.Count++
			if r.Timestamp.Before(user.First) {
				user.First = r.Timestamp
			}
			if r.Timestamp.After(user.Last) {
				user.Last = r.Timestamp
			}
			stats.UserStats[r.UserID] = user
		} else {
			stats.UserStats[r.UserID] = UserStat{
				Count: 1,
				First: r.Timestamp,
				Last:  r.Timestamp,
			}
		}

		// 行为统计
		stats.ActionStats[r.Action]++

		// 分钟统计
		truncated := r.Timestamp.Truncate(time.Minute)
		if _, exists := stats.MinuteStats[truncated]; !exists {
			stats.MinuteStats[truncated] = MinuteStat{
				TimeWindow:   truncated,
				ActiveUsers:  make(map[string]struct{}),
				TotalActions: 0,
			}
		}
		ms := stats.MinuteStats[truncated]
		ms.ActiveUsers[r.UserID] = struct{}{}
		ms.TotalActions++
		stats.MinuteStats[truncated] = ms
	}

	return stats
}
