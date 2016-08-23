package helpers

import (
	"strings"
)


func LxcInfoToInterface(info string) interface{} {
	result := map[string]interface{}{}
	if info == "" { return result }
	for _, line := range strings.Split(info, "\n") {
		data := strings.Split(line, ":")
		key := strings.Trim(data[0], " ")
		if val, ok := result[key].(string); ok {
		    result[key] = val + ", " + strings.Trim(data[1], " ")
		} else {
			result[key] = strings.Trim(data[1], " ")
		}
	}
	return result
}

func LxcListToInterface(list string) interface{} {
	result := []interface{}{}
	if list == "" { return result }
	for i, line := range strings.Split(list, "\n") {
		var arr []string
		for _, d := range strings.Split(line, " ") {
			d = strings.TrimSpace(d)
			if d != "" {
				arr = append(arr, d)
			}
		}
		if len(arr) < 2 || i == 0 { continue }
		result = append(result, map[string]interface{}{
			"name": arr[0],
			"state": arr[1],
			"ipv4": arr[2],
			"ipv6": arr[3],
			"autostart": arr[4],
			"pid": arr[5],
			"ram": arr[6],
			"swap": arr[7],
		})
	}
	return result
}

func LxcOutToList(list string) interface{} {
	if list == "" { return []interface{}{} }
	return strings.Split(list, "\n")
}
