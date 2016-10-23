package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var websites []string

func main() {
	//main websites
	websites = []string{
		"http://www.jd.com", "http://www.taobao.com", "http://www.tmall.com", "http://www.sohu.com",
		"http://www.163.com", "http://www.ifeng.com", "http://www.suning.com", "http://www.ctrip.com",
		"http://www.gome.com.cn", "http://www.btime.com", "http://www.pconline.com.cn", "http://www.ly.com",
		"http://www.zol.com.cn", "http://www.autohome.com.cn", "http://www.anjuke.com", "http://www.amazon.cn",
	}
	//并发数量
	pnum := 10
	goRequest(pnum, websites)
}

func goRequest(pnum int, websites []string) {
	total := len(websites)
	if pnum <= 0 { //如果为0，默认全部并发
		pnum = total
	}
	if pnum > total {
		pnum = total
	}

	startTime := time.Now().UnixNano()
	fetchData := make(map[string]string, total) //反馈抓取后的数据结果
	execChans := make(chan bool, pnum)          //控制并发数量的通道，pnum指定大小，超过则阻塞
	doneChans := make(chan bool, 1)             //用来传递完成信号

	for i := 0; i < total; i++ {
		go request(i, websites[i], execChans, doneChans, fetchData) //并发执行
	}
	for i := 0; i < total; i++ {
		r := <-doneChans //完成一个，同时获取下一个任务
		<-execChans      //读取下一个任务
		if !r {
			log.Printf("第%d项,URL:%s获取失败", i, websites[i])
		}
	}

	close(doneChans)
	close(execChans)
	processTime := float32(time.Now().UnixNano()-startTime) / 1e9 //总耗时
	log.Printf("全部完成，并发数量：%d, 耗时：%.3fs", pnum, processTime)
	for k, v := range fetchData {
		log.Printf("URL:%s, data:%s", k, v)
	}
	log.Printf("data: %q", fetchData)
}

func request(i int, url string, execChans chan bool, doneChans chan bool, fetchData map[string]string) {
	execChans <- true //放在函数开始处，用来阻塞执行
	log.Printf("No: %02d, url: %s, start...", i, url)
	isOk := false
	startTime := time.Now().UnixNano()
	resp, _ := http.Get(url)

	defer func() {
		resp.Body.Close()
		doneChans <- isOk
		processTime := float32(time.Now().UnixNano()-startTime) / 1e9
		log.Printf("No: %02d, url: %s, status: %t, time: %.3fs end.", i, url, isOk, processTime)
	}()

	body, err := ioutil.ReadAll(resp.Body)
	len := len(body)
	log.Printf("No: %02d, len: %d", i, len)
	fetchData[url] = fmt.Sprintf("len: %d", len)
	if err == nil {
		isOk = true
	}
}
