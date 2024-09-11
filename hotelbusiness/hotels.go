//go:build !solution

package hotelbusiness

import "sort"

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	loadMap := make(map[int]int)
	for _, guest := range guests {
		loadMap[guest.CheckInDate] += 1
		loadMap[guest.CheckOutDate] -= 1
	}

	var dates []int
	for date := range loadMap {
		dates = append(dates, date)
	}
	sort.Ints(dates)

	var loads []Load
	currentLoad := 0
	for _, date := range dates {
		loadChange := loadMap[date]
		currentLoad += loadChange
		if loadChange != 0 {
			load := Load{StartDate: date, GuestCount: currentLoad}
			loads = append(loads, load)
		}
	}

	return loads
}
