package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/tebeka/selenium"
)

const (
	ChormeDriverPath  = "/home/keima/Desktop/chromedriver/chromedriver"
	EdgeDriverPath    = "/home/keima/Desktop/edgedriver/msedgedriver"
	FirefoxDriverPath = "/home/keima/Desktop/geckodriver/geckodriver"
	port              = 8081
	ticketBuyCount    = 2
	cookieName        = "./cookies.txt"
)

// var wd selenium.WebDriver

// func main() {
// 	// 测试一下. 首先配置服务选项
// 	opt := []selenium.ServiceOption{
// 		selenium.Output(os.Stderr),         // 内容输出 这里是标准输出到控制台
// 		selenium.GeckoDriver(damai.Driver), // 配置驱动地址 chrome参考官网
// 	}
// 	// chrome请参考一下官方文档. 有部份不同
// 	config := damai.NewServiceConfig(damai.SeleniumPath, damai.Port, false)                      // 创建服务配置
// 	config.AddServiceOption(opt...)                                                              // 添加
// 	service, driver, _ := damai.NewWebDriver(config, damai.BrowserName)                          // 创建驱动
// 	defer service.Stop()                                                                         // 延时停止服务
// 	const url = "https://www.baidu.com"                                                          // 请求地址
// 	driver.Get(url)                                                                              // 打开浏览器并加载地址
// 	defer driver.QuitBrowser()                                                                   // 延时关闭浏览器
// 	s := damai.NewSliceMrg([]string{"FindElement"}, []string{selenium.ByID}, []string{"kw"}, 20) // 多步骤管理
// 	for _, v := range s {
// 		// 执行步骤execute
// 		}
// 	}

// }

func main() {

	wd, _ := chrome()

	// firefox()
	// 先写cookie
	// cookie(wd)

	// 再抢票
	getTicket(wd)

}

func chrome() (selenium.WebDriver, error) {
	var err error
	// Start a WebDriver server instance
	opts := []selenium.ServiceOption{
		// selenium.Output(os.Stderr), // Output debug information to STDERR.
	}
	// selenium.SetDebug(false)
	_, err = selenium.NewChromeDriverService(ChormeDriverPath, port, opts...)

	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))

	if err != nil {
		panic(err)
	}
	return wd, nil
}

func firefox() (selenium.WebDriver, error) {
	var err error
	opts := []selenium.ServiceOption{
		// selenium.StartFrameBuffer(), // Start an X frame buffer for the browser to run in.
		// selenium.GeckoDriver(FirefoxDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		// selenium.Output(os.Stderr),              // Output debug information to STDERR.
	}
	service, err := selenium.NewGeckoDriverService(FirefoxDriverPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", port))
	// wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		panic(err)
	}
	return wd, nil
}

func cookie(wd selenium.WebDriver) {
	// 先打开网页，然后登录，并将相关的cookie写入到本地
	// 打开网页
	if err := wd.Get("https://passport.damai.cn/login?ru=https%3A%2F%2Fwww.damai.cn%2F"); err != nil {
		panic(err)
	}
	// 等到30秒进行登录操作
	time.Sleep(30 * time.Second)
	dictCookies, err := wd.GetCookies()
	if err != nil {
		panic(fmt.Sprintf("获取cookies失败”, err:%s", err))
	}
	// 将cookie保存至本地
	f, _ := os.OpenFile(cookieName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	cookieJson, _ := json.Marshal(dictCookies)
	f.Write(cookieJson)
}

func getTicket(wd selenium.WebDriver) {
	// 将cookie读出来然后写入
	cookieByte, _ := os.ReadFile(cookieName)
	// fmt.Println(string(cookieByte))
	var cookiesDicts []selenium.Cookie
	err := json.Unmarshal(cookieByte, &cookiesDicts)
	// 打开网页
	if err := wd.Get("https://detail.damai.cn/item.htm?spm=a2oeg.search_category.0.0.79664d15u8c2PP&id=711368998162&clicktitle=2023%E5%BC%A0%E4%BF%A1%E5%93%B2%E3%80%8C%E6%9C%AA%E6%9D%A5%E5%BC%8F2.0%E3%80%8D%E4%B8%96%E7%95%8C%E5%B7%A1%E5%9B%9E%E6%BC%94%E5%94%B1%E4%BC%9A-%E6%AD%A6%E6%B1%89%E7%AB%99"); err != nil {
		panic(err)
	}

	for _, v := range cookiesDicts {
		if v.Domain == "www.damai.cn" {
			v.Domain = ".damai.cn"
		}
		if v.Expiry == 0 {
			v.Expiry = math.MaxUint32
		}
		err = wd.AddCookie(&v)
	}
	wd.Refresh()

	// 获取城市相关按钮，并选择武汉
	elems, err := wd.FindElements(selenium.ByCSSSelector, ".cityitem")
	if err != nil {
		panic(fmt.Sprintf("选择城市失败, err:%s", err))
	}
	for _, v := range elems {
		city, _ := v.Text()
		if city == "武汉站" {
			if err := v.Click(); err != nil {
				panic(fmt.Sprintf("点击城市失败, err:%s", err))
			}
			break
		}
	}

	// 场次默认无需指定

	// 票档指定一下
	elems, err = wd.FindElements(selenium.ByCSSSelector, ".skuname")
	if err != nil {
		panic(fmt.Sprintf("选择城市失败, err:%s", err))
	}
	for _, v := range elems {
		// 从前往后选，直到选到有票为止
		if _, err := v.FindElement(selenium.ByCSSSelector, ".notticket"); err != nil {
			// 如果这里报错，则证明没找到对应的class，代表有票，可以点击
			if err = v.Click(); err != nil {
				panic(fmt.Sprintf("点击票档错误, err:%s", err))
			}
		}
	}

	// 数量修改
	// elem, err := wd.FindElement(selenium.ByCSSSelector, ".cafe-c-input-number-input")
	// if err != nil {
	// 	panic(fmt.Sprintf("选择数量失败, err:%s", err))
	// }

	// count, err := elem.Text()
	// if err != nil {
	// 	panic(fmt.Sprintf("查看张数失败, err:%s", err))
	// }
	// 如果数量是1，那么点击加号
	//if count == "1" && ticketBuyCount > 1 {
	elem, err := wd.FindElement(selenium.ByCSSSelector, ".cafe-c-input-number-handler-up")
	if err != nil {
		panic(fmt.Sprintf("选择+号失败, err:%s", err))
	}
	for i := 2; i <= ticketBuyCount; i++ {
		if err = elem.Click(); err != nil {
			panic(fmt.Sprintf("点击+号失败, err:%s", err))
		}
	}
	//}

	time.Sleep(1 * time.Second)

	// 点击“不，立即预订”
	elem, err = wd.FindElement(selenium.ByCSSSelector, ".buy-link")
	if err != nil {
		panic(fmt.Sprintf("选择“不，立即预订失败”, err:%s", err))
	}
	if err = elem.Click(); err != nil {
		panic(fmt.Sprintf("点击“不，立即预订失败”, err:%s", err))
	}

	time.Sleep(200 * time.Second)
}

// elem, err := wd.FindElement(selenium.ByCSSSelector, "#code")
// 	if err != nil {
// 		panic(err)
// 	}
// 	// Remove the boilerplate code already in the text box.
// 	if err := elem.Clear(); err != nil {
// 		panic(err)
// 	}

// 	// Enter some new code in text box.
// 	err = elem.SendKeys(`
//         package main
//         import "fmt"
//         func main() {
//             fmt.Println("Hello WebDriver!")
//         }
//     `)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Click the run button.
// 	btn, err := wd.FindElement(selenium.ByCSSSelector, "#run")
// 	if err != nil {
// 		panic(err)
// 	}
// 	if err := btn.Click(); err != nil {
// 		panic(err)
// 	}

// 	time.Sleep(5 * time.Second)

// 	// Wait for the program to finish running and get the output.
// 	outputDiv, err := wd.FindElement(selenium.ByCSSSelector, ".Playground-output")
// 	if err != nil {
// 		panic(err)
// 	}

// 	var output string
// 	for {
// 		output, err = outputDiv.Text()
// 		if err != nil {
// 			panic(err)
// 		}
// 		if output != "Waiting for remote server..." {
// 			break
// 		}
// 		time.Sleep(time.Millisecond * 100)
// 	}

// 	fmt.Printf("%s", strings.Replace(output, "\n\n", "\n", -1))
