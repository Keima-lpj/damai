package damai

import (
	"fmt"
	"reflect"
	"time"

	"github.com/tebeka/selenium"
)

const (
	// seleniumJar包存放位置
	SeleniumPath = "/home/keima/Desktop/chromedriver/selenium-server-4.7.1.jar"
	// 驱动存放位置
	Driver = "/home/keima/Desktop/chromedriver/chromedriver"
	// 打开操作浏览器的名称(名称与驱动一致)
	BrowserName = "chrome"
	// 端口
	Port = 9091
)

// 用于管理服务配置
type ServiceConfig struct {
	SeleniumPath string
	Port         int
	Config       []selenium.ServiceOption
	SetDebug     bool
}

// 用于管理步骤
type CaseMrg struct {
	BySelect, Element string
	MethodName        string
	Timeout           time.Duration
	Value             []reflect.Value
}

// 用于管理driver
type WebDriver struct {
	Driver selenium.WebDriver
}

// 添加服务选项
func (s *ServiceConfig) AddServiceOption(opt ...selenium.ServiceOption) {
	s.Config = append(s.Config, opt...)
}

// Get方法
func (d WebDriver) Get(url string) error {
	return d.Driver.Get(url)
}

// 退出浏览器方法
func (d WebDriver) QuitBrowser() error {
	return d.Driver.Quit()
}

// 构造函数 初始化配置
func NewServiceConfig(seleniumPath string, Port int, setDebug bool) ServiceConfig {
	return ServiceConfig{
		SeleniumPath: seleniumPath,
		Port:         Port,
		Config:       append(([]selenium.ServiceOption)(nil)),
		SetDebug:     setDebug,
	}
}

// 构造函数 初始化服务
func NewService(config ServiceConfig) *selenium.Service {
	selenium.SetDebug(config.SetDebug)
	service, err := selenium.NewSeleniumService(config.SeleniumPath, config.Port, config.Config...)
	if err != nil {
		panic(err)
	}
	return service
}

// 构造函数 初始化driver,和服务
func NewWebDriver(config ServiceConfig, browser string) (service *selenium.Service, w WebDriver, err error) {
	service = NewService(config)
	opt := selenium.Capabilities{"browserName": browser}
	w.Driver, err = selenium.NewRemote(opt, fmt.Sprintf("http://localhost:%d/wd/hub", config.Port))
	if err != nil {
		panic(err)
	}
	return service, w, nil
}

/*
	轮询等待 二次封装

等待元素可见与python中的wait.until一致

	WebDriverWait(self.driver, timeout=timeout).until(lambda x: x.find_element(("id", el)), message=msg)
*/
func WaitElementTimeout(d WebDriver, timeout time.Duration, by, value string) error {
	var IsElementDisplay = func(wd selenium.WebDriver) (bool, error) {
		ele, err := wd.FindElement(by, value)
		if err != nil {
			return false, err
		}
		b, errs := ele.IsDisplayed()
		if errs != nil {
			return false, err
		}
		return b, nil

	}
	if err := d.Driver.WaitWithTimeoutAndInterval(IsElementDisplay, timeout, 500*time.Millisecond); err != nil {
		return err
	}
	return nil
}

// 管理步骤
func NewSliceMrg(names, BySelect, Ele []string, timeout time.Duration) []CaseMrg {
	DataList := make([]CaseMrg, 0)
	if len(names) >= 1 {
		for i, v := range names {
			d := CaseMrg{
				BySelect:   BySelect[i],
				Element:    Ele[i],
				MethodName: v,
				Timeout:    timeout,
				Value:      []reflect.Value{reflect.ValueOf(BySelect[i]), reflect.ValueOf(Ele[i])},
			}
			DataList = append(DataList, d)
		}
	}
	return DataList

}

// 通过反射拿到方法名称然后执行
func exeCommand(v reflect.Value, methodsName string, val []reflect.Value) []reflect.Value {
	m := v.MethodByName(methodsName)
	return m.Call(val)
}

// 执行用例
func Run(d WebDriver, dr CaseMrg) error {
	// 反射
	rfv := reflect.ValueOf(d.Driver)
	// 等待元素出现
	if err := WaitElementTimeout(d, 50*time.Second, dr.BySelect, dr.Element); err != nil {
		fmt.Println("等待超時:", err)
		time.Sleep(10 * time.Second)
		return err
	}
	// 执行
	result := exeCommand(rfv, dr.MethodName, dr.Value)
	// 打印返回数据(selenium.Element func)
	fmt.Println(result[0])
	return nil
}
