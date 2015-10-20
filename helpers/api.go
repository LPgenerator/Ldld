package helpers

import (
	"strings"
)


func LxcInfoToInterface(info string) interface{} {
	result := map[string]interface{}{}
	if info == "" { return result }
	for _, line := range strings.Split(info, "\n") {
		data := strings.Split(line, ":")
		result[strings.Trim(data[0], " ")] = strings.Trim(data[1], " ")
	}
	return result
}

func LxcListToInterface(list string) interface{} {
	result := []interface{}{}
	if list == "" { return result }
	for i, line := range strings.Split(list, "\n") {
		var arr []string
		for _, d := range strings.Split(line, "  ") {
			d = strings.Trim(d, " ")
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
			"memory": arr[6],
			"ram": arr[7],
			"swap": arr[8],
		})
	}
	return result
}

func LxcOutToList(list string) interface{} {
	if list == "" { return []interface{}{} }
	return strings.Split(list, "\n")
}

func ImagesToInterface(list string) interface{} {
	if list == "" { return []interface{}{} }
	return strings.Split(list, "\n")
}
