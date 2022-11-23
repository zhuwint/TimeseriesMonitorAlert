package sensor

import "timeseries/pkg/api"

// 根据 location 列表 生成4级 location 树
func buildIndexTree(array []*api.LocationNode) []*api.LocationNode {
	maxLen := len(array)
	var isVisit = make([]bool, maxLen)

	var root []*api.LocationNode
	for i := 0; i < maxLen; i++ {
		if array[i].Level == 1 {
			root = append(root, array[i])
			continue
		} else if array[i].Level == 2 {
			for j := 0; j < maxLen; j++ {
				if array[j].Location1Id == array[i].Location1Id && array[j].Level == 1 {
					array[j].Children = append(array[j].Children, array[i])
					isVisit[i] = true
				}
			}
		} else if array[i].Level == 3 {
			for j := 0; j < maxLen; j++ {
				if array[j].Location1Id == array[i].Location1Id && array[j].Location2Id == array[i].Location2Id && array[j].Level == 2 {
					array[j].Children = append(array[j].Children, array[i])
					isVisit[i] = true
				}
			}
		} else {
			for j := 0; j < maxLen; j++ {
				if array[j].Location1Id == array[i].Location1Id && array[j].Location2Id == array[i].Location2Id && array[j].Location3Id == array[i].Location3Id && array[j].Level == 3 {
					array[j].Children = append(array[j].Children, array[i])
					isVisit[i] = true
				}
			}
		}
	}
	return root
}
