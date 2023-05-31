package helper

import (
	"strconv"
	"strings"
)

func ConvertStringToBulkInt64(str string) ([]int64, error) {
	var ids []int64
	strSplit := strings.Split(str, ",")
	for _, id := range strSplit {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return []int64{}, err
		}
		ids = append(ids, int64(idInt))
	}

	return ids, nil
}
