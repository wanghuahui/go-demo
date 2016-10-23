package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func main() {
	s := "http://api.test.com/path1/path2/like?Id=1234&param1=2222&param2=abcd"
	//解析url
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	//协议
	fmt.Printf("scheme: %s\n", u.Scheme)
	// User，包含所有认证信息
	if u.User != nil {
		fmt.Printf("user: %s\n", u.User)
		fmt.Printf("user name: %s\n", u.User.Username())
		p, _ := u.User.Password()
		fmt.Printf("user password: %s\n", p)
	}

	// Host 同时包括主机名和端口信息，如果端口存在，使用Split提取
	fmt.Printf("host: %s\n", u.Host)
	h := strings.Split(u.Host, ":")
	fmt.Printf("host name: %s\n", h[0])
	// 判断是否包含端口
	if len(h) > 1 {
		port, _ := strconv.Atoi(h[1])
		fmt.Printf("port: %d\n", port)
	}

	// 路径
	fmt.Printf("path: %s\n", u.Path)
	// 参数信息
	fmt.Printf("paramters: %s\n", u.RawQuery)
	// 将查询参数解析为map类型
	m, _ := url.ParseQuery(u.RawQuery)
	fmt.Printf("values: %s\n", m)
	// map以查询字符串为键，对应字符串切片为值，第一个值为[0]
	ID, _ := strconv.Atoi(m["Id"][0])
	param1 := m["param1"][0]
	param2 := m["param2"][0]
	fmt.Printf("Id: %d, param1: %s, param2: %s\n", ID, param1, param2)
}
