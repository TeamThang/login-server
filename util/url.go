package util

// 组合两个路径
func PathJoin(url1 string, url2 string) string {
	if string(url1[len(url1)-1]) == "/" && string(url2[0]) == "/" {
		return url1[:len(url1)-1] + url2
	}
	if string(url1[len(url1)-1]) == "/" || string(url2[0]) == "/" {
		return url1 + url2
	}
	return url1 + "/" + url2
}
