# Golang使用goquery抓取京东商品价格Demo

本文主要介绍基于goquery实现一个简单的京东商品价格数据爬虫

## 分析京东网页分页地址

- 进入京东首页搜索mac，获取浏览器地址，下面的地址是无法分页抓取的

  https://search.jd.com/Search?keyword=mac&enc=utf-8&wq=mac&pvid=7faaf245e0a74dfcbe6d80da0ffd9647

- 点击页面底部分页按钮获取数据分页链接地址，继续点击分页链接，分析发现只是`page`参数和`s`参数会变化。多次尝试后发现只有`keyword`、`enc`、`wq`、`page`对页面搜索数据有影响

  https://search.jd.com/Search?keyword=mac&enc=utf-8&qrst=1&rt=1&stop=1&vt=2&wq=mac&page=3&s=58&click=0

  > keyword和wq为商品搜索内容这里使用 mac为关键字
  >
  > enc 指定返回数据字符集 这里为utf-8

- 所以最终使用的抓取地址就是下面这个了，获取下一页数据只需要page加1就可以

  https://search.jd.com/Search?keyword=mac&enc=utf-8&wq=mac&page=1

## 分页网页结构

- 获取主要内容组件 分析得到 id为"J_goodsList" 的div为主要内容列表
- 里面UL列表为每个商品展示数据
- LI节点中 class为"p-price"的div为价格数据
- LI节点中class为"p-name p-name-type-2"的div为商品名称数据

## 安装依赖包

此项目只是抓取数据存入本地文件，所以第三方包只依赖goquery

```shell
github.com/PuerkitoBio/goquery
```

## 编码

```go
package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"jdSpider/util"
	"log"
	"net/http"
	"strings"
)

func main() {
	// 测试功能只抓取10页数据
	for i := 1; i < 10; i++ {
		url := fmt.Sprintf("https://search.jd.com/Search?keyword=mac&enc=utf-8&wq=mac&page=%d", i)
		fetchData(url)
	}

}

/**
抓取数据
*/
func fetchData(url string) {
	fmt.Println(url)
	client := http.Client{}
	request, err := http.NewRequest("GET", url, strings.NewReader("name=cjb"))
	if err != nil {
		log.Println(err)
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Linux; U; Android 5.1; zh-cn; m1 metal Build/LMY47I) AppleWebKit/537.36 (KHTML, like Gecko)Version/4.0 Chrome/37.0.0.0 MQQBrowser/7.6 Mobile Safari/537.36")

	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}
	// 使用NewDocumentFromResponse方式获取获取数据，是应为某些网页会有防止爬取限制，需要设置Header防止被限制
	doc, err := goquery.NewDocumentFromResponse(response)
	/**
	1：获取ID为J_goodsList 的div节点
	2：获取ul节点
	3：获取li节点列表
	*/
	doc.Find("div[id=\"J_goodsList\"]").Find("ul").Find("li").Each(func(i int, selection *goquery.Selection) {
		// 获取class为p-name p-name-type-2 的div节点，然后获取em子节点的文字内容作为商品标题
		title := selection.Find("div[class=\"p-name p-name-type-2\"]").Find("em").Text()
		// 获取class为p-price的节点，然后获取i标签中的文字作为价格
		price := selection.Find(".p-price").Find("i").Text()
		// 列表中有部分内容是广告内容，不属于标准商品数据，这里排除掉
		if len(title) > 1 {
			// 把获取到的数据追加到jdprive.txt 文件中，格式为  商品名称+四个制表符+价格+换行
			util.AppendToFile("jdprive.txt", title+"\t\t\t\t"+price+"\n")
		}

	})
}
```

## 说明

京东网页HTML结构可能随时会变动，这个程序你看到的时候可能已经无法抓取到数据。但是按照上面的说明可以很容易抓取新网页内容的数据。
jdprive.txt 文件为抓取的内容数据
项目doc目录下面的jd.html 文件是我写这个程序时候保存的网页数据