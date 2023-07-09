package dag

import "strings"

const sep = "----"

type Link struct {
	From string
	To   string
}

func GetDeployOrder(links []Link) [][]string {
	_, from2ToMap, _ := getNodesAndMap(links)

	// 找出依赖环
	circles := [][]string{}
	for _, link := range links {
		to := link.To
		circles = append(circles, walkAndGeCircles([]string{link.From, to}, from2ToMap)...)
	}

	if len(circles) > 0 {
		circles = deduplicateCircles(circles)
	}

	// 将存在依赖环的节点改成拼接起来
	convertLinks := []Link{}
	for i, link := range links {
		for _, circle := range circles {
			if find(circle, link.From) {
				link.From = strings.Join(circle, sep)
			}
			if find(circle, link.To) {
				link.To = strings.Join(circle, sep)
			}
		}
		links[i] = link
		if link.From != link.To {
			convertLinks = append(convertLinks, link)
		}
	}

	result := walkLink(convertLinks, nil)

	// 将字符串转换成数字
	for i, group := range result {
		newGroup := []string{}
		for _, item := range group {
			if strings.Contains(item, sep) {
				newGroup = append(newGroup, strings.Split(item, sep)...)
			} else {
				newGroup = append(newGroup, item)
			}
		}
		result[i] = newGroup
	}

	return result
}

func findIndex(arr []string, target string) int {
	for i, item := range arr {
		if item == target {
			return i
		}
	}
	return -1
}

func find(arr []string, target string) bool {
	return findIndex(arr, target) > -1
}

func sliceArr(arr []string, from int, to int) []string {
	if from < 0 {
		panic("from must be positive integer")
	}

	if to == -1 {
		to = len(arr)
	}

	if to > len(arr) {
		to = len(arr)
	}

	res := []string{}
	for i := from; i < to; i++ {
		res = append(res, arr[i])
	}
	return res
}

func appendUnique(arr []string, target string) []string {
	if !find(arr, target) {
		arr = append(arr, target)
	}
	return arr
}

func walkAndGeCircles(chain []string, fromMap map[string][]string) [][]string {
	from := chain[len(chain)-1]
	toes := fromMap[from]
	circles := [][]string{}
	if len(toes) > 0 {
		for _, to := range toes {
			index := findIndex(chain, to)
			if index > -1 {
				// 存在依赖环
				circles = append(circles, sliceArr(chain, index, -1))
				break
			} else {
				circles = append(circles, walkAndGeCircles(append(chain, to), fromMap)...)
			}
		}
	}
	return circles
}

func deduplicateCircles(circles [][]string) [][]string {
	deduped := [][]string{}
	others := [][]string{}

	circle := circles[0]
	for i := 1; i < len(circles); i++ {
		item := circles[i]
		finded := false
		for _, node := range item {
			if find(circle, node) {
				// 如果两个环有重合，合并环
				finded = true
			}
		}
		if finded {
			for _, node := range item {
				circle = appendUnique(circle, node)
			}
		} else {
			others = append(others, item)
		}
	}
	deduped = append(deduped, circle)
	if len(others) > 0 {
		deduped = append(deduped, deduplicateCircles(others)...)
	}
	return deduped
}

func getNodesAndMap(links []Link) ([]string, map[string][]string, map[string][]string) {
	from2ToMap := map[string][]string{}
	to2FromMap := map[string][]string{}

	nodes := []string{}

	for _, link := range links {
		toes := from2ToMap[link.From]
		if len(toes) == 0 {
			toes = []string{}
		}
		from2ToMap[link.From] = append(toes, link.To)

		strings := from2ToMap[link.To]
		if len(strings) == 0 {
			strings = []string{}
		}
		from2ToMap[link.To] = strings

		froms := to2FromMap[link.To]
		if len(froms) == 0 {
			froms = []string{}
		}
		to2FromMap[link.To] = append(froms, link.From)

		nodes = appendUnique(nodes, link.From)
		nodes = appendUnique(nodes, link.To)
	}
	return nodes, from2ToMap, to2FromMap
}

func walkLink(links []Link, orphanNodes []string) [][]string {
	_, from2ToMap, to2FromMap := getNodesAndMap(links)
	leaf := []string{}

	if len(orphanNodes) > 0 {
		leaf = orphanNodes
	}

	nextSuspectedOrphanNodes := []string{}
	for from, toes := range from2ToMap {
		if len(toes) == 0 {
			leaf = appendUnique(leaf, from)

			// 查询依赖该叶子节点的节点是否为疑似孤立节点
			froms4Leaf := to2FromMap[from]
			if len(froms4Leaf) > 0 {
				for _, from4Leaf := range froms4Leaf {
					if to2FromMap[from4Leaf] == nil || len(to2FromMap[from4Leaf]) == 0 {
						nextSuspectedOrphanNodes = appendUnique(nextSuspectedOrphanNodes, from4Leaf)
					}
				}
			}
		}
	}

	// 排除叶子节点继续递归剩余部分
	nextLinks := []Link{}
	for _, link := range links {
		if !find(leaf, link.To) {
			nextLinks = append(nextLinks, link)
		}
	}

	result := [][]string{leaf}
	if len(nextLinks) > 0 {
		nextNodes, _, _ := getNodesAndMap(nextLinks)
		nextOrphanNodes := excludes(nextSuspectedOrphanNodes, nextNodes)
		result = append(result, walkLink(nextLinks, nextOrphanNodes)...)
	} else if len(nextSuspectedOrphanNodes) > 0 {
		nextOrphanNodes := excludes(nextSuspectedOrphanNodes, leaf)
		if len(nextOrphanNodes) > 0 {
			result = append(result, nextOrphanNodes)
		}
	}
	return result
}

func excludes(sourceArr []string, excludesArr []string) []string {
	res := []string{}
	for _, node := range sourceArr {
		if !find(excludesArr, node) {
			res = appendUnique(res, node)
		}
	}
	return res
}
