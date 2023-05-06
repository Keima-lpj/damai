# -*- coding: utf-8 -*-

import os  # 创建文件夹, 文件是否存在
import time  # time 计时
import pickle  # 保存和读取cookie实现免登陆的一个工具
from time import sleep
from selenium import webdriver  # 操作浏览器的工具
from selenium.webdriver.common.by import By
import argparse


# 大麦网主页
damai_url = 'https://www.damai.cn/'
# 登录
login_url = 'https://passport.damai.cn/login?ru=https%3A%2F%2Fwww.damai.cn%2F'
# 抢票目标页
target_url = 'https://detail.damai.cn/item.htm?spm=a2oeg.search_category.0.0.79664d15u8c2PP&id=711368998162&clicktitle=2023%E5%BC%A0%E4%BF%A1%E5%93%B2%E3%80%8C%E6%9C%AA%E6%9D%A5%E5%BC%8F2.0%E3%80%8D%E4%B8%96%E7%95%8C%E5%B7%A1%E5%9B%9E%E6%BC%94%E5%94%B1%E4%BC%9A-%E6%AD%A6%E6%B1%89%E7%AB%99'
# 选择城市
city = 1
# 选择场次   默认第一个场次
sessions = 1
# 选择票档  默认第一个票档，顺序从左往右从上往下
ticket_stalls = 1
# 选择数量
ticket_num = 1
# 驱动
chrome_driver = '/home/keima/Desktop/chromedriver/chromedriver'

# class Concert:
class Concert:
    # 初始化加载
    def __init__(self):
        self.status = 0  # 状态, 表示当前操作执行到了哪个步骤
        self.login_method = 1  # {0:模拟登录, 1:cookie登录}自行选择登录的方式
        # 设置浏览器选项
        option = webdriver.ChromeOptions()
        option.add_argument('user-agent="Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"')
        option.add_experimental_option('useAutomationExtension', False)  # 去掉开发者警告
        option.add_experimental_option('excludeSwitches', ['enable-automation'])
        option.add_argument('--disable-blink-features=AutomationControlled')
        
        self.driver = webdriver.Chrome(executable_path=chrome_driver, options=option)  # 当前浏览器驱动对象
        
    # cookies: 登录网站时出现的 记录用户信息用的
    def set_cookies(self):
        """cookies: 登录网站时出现的 记录用户信息用的"""
        self.driver.get(damai_url)
        print('###请点击登录###')
        # 我没有点击登录,就会一直延时在首页, 不会进行跳转
        while self.driver.title.find('大麦网-全球演出赛事官方购票平台') != -1:
            sleep(1)
        print('###请扫码登录###')
        # 没有登录成功
        while self.driver.title != '大麦网-全球演出赛事官方购票平台-100%正品、先付先抢、在线选座！':
            sleep(1)
        print('###扫码成功###')
        # get_cookies: driver里面的方法
        pickle.dump(self.driver.get_cookies(), open('cookies.pkl', 'wb'))
        print('###cookie保存成功###')
        self.driver.get(target_url)

    # 假如说我现在本地有 cookies.pkl 那么 直接获取
    def get_cookie(self):
        """假如说我现在本地有 cookies.pkl 那么 直接获取"""
        cookies = pickle.load(open('cookies.pkl', 'rb'))
        for cookie in cookies:
            cookie_dict = {
                'domain': '.damai.cn',  # 必须要有的, 否则就是假登录
                'name': cookie.get('name'),
                'value': cookie.get('value')
            }
            self.driver.add_cookie(cookie_dict)
        print('###载入cookie###')

    def login(self):
        """登录"""
        if self.login_method == 0:
            self.driver.get(login_url)
            print('###开始登录###')
        elif self.login_method == 1:
            # 创建文件夹, 文件是否存在
            if not os.path.exists('cookies.pkl'):
                self.set_cookies()  # 没有文件的情况下, 登录一下
            else:
                self.driver.get(target_url)  # 跳转到抢票页
                self.get_cookie()  # 并且登录

    def enter_concert(self):
        """打开浏览器"""
        print('###打开浏览器,进入大麦网###')
        # 调用登录
        self.login()  # 先登录再说
        self.driver.refresh()  # 刷新页面
        self.status = 2  # 登录成功标识
        print('###登录成功###')
        # 处理弹窗
        if self.isElementExist('/html/body/div[2]/div[2]/div/div/div[3]/div[2]'):
            self.driver.find_element(By.XPATH, '/html/body/div[2]/div[2]/div/div/div[3]/div[2]').click()

    # 二. 抢票并且下单
    def choose_ticket(self):
        """选票操作"""
        if self.status == 2:
            print('=' * 30)
            print('###开始城市，场次，票档，数量的选择###')

            # 这里先找城市
            cityXpath = "/html/body/div[2]/div[1]/div[1]/div[1]/div[1]/div[2]/div[3]/div[1]/div[%d]" % (city)
            self.driver.find_element(By.XPATH, cityXpath).click()
            time.sleep(1)

            # 再找场次
            sessXpath = "/html/body/div[2]/div[1]/div[1]/div[1]/div[1]/div[2]/div[4]/div[3]/div[2]/div/div[%d]" % (sessions)
            self.driver.find_element(By.XPATH, sessXpath).click()
            time.sleep(1)

            # 再找票档
            ticketStallsXpath = "/html/body/div[2]/div[1]/div[1]/div[1]/div[1]/div[2]/div[4]/div[5]/div[2]/div/div[%d]" % (ticket_stalls)
            self.driver.find_element(By.XPATH, ticketStallsXpath).click()
            time.sleep(1)

            # 再找数量
            if ticket_num > 1:
                ticketNumXpath = "/html/body/div[2]/div[1]/div[1]/div[1]/div[1]/div[2]/div[4]/div[6]/div[2]/div/div/a[2]"
                ticketNum = self.driver.find_element(By.XPATH, ticketNumXpath)
                for x in range (1, ticket_num):
                    ticketNum.click()
            time.sleep(1)        

            while self.driver.title.find("确认订单") == -1:
                try:
                    buybutton = self.driver.find_element(By.CLASS_NAME, 'buy-link').text
                    if buybutton == '提交缺货登记':
                        self.status = 2  # 没有进行更改操作
                        self.driver.get(target_url)  # 刷新页面 继续执行操作
                    elif buybutton == '立即预定':
                        # 点击立即预定
                        self.driver.find_element('buybtn').click()
                        self.status = 3
                    elif buybutton == '立即购买':
                        self.driver.find_element(By.CLASS_NAME, 'buybtn').click()
                        self.status = 4
                    elif buybutton == '选座购买':
                        self.driver.find_element(By.CLASS_NAME, 'buybtn').click()
                        self.status = 5
                    elif buybutton == "不，立即购买":
                        self.driver.find_element(By.CLASS_NAME, 'buy-link').click()
                        self.status = 5
                    elif buybutton == "不，立即预订":
                        self.driver.find_element(By.CLASS_NAME, 'buy-link').click()
                        self.status = 5
                except:
                    print('###没有跳转到订单结算界面###')
                    return False
                title = self.driver.title
                if title == '选座购买':
                    # 选座购买的逻辑
                    self.choice_seats()
                elif title == '订单确认页':
                    # 实现下单的逻辑
                    while True:
                        time.sleep(1)
                        # 如果标题为确认订单
                        print('正在加载.......')

                        # 找到第一位观演人，点击按钮选中
                        # 根据前面选择的数量，循环勾选对应的观演人。 注意这里只能从第一个人开始选
                        for x in range (1, ticket_num+1):
                            print('开始买第%d张票' % x)
                            ticketNumXpath = '//*[@id="dmViewerBlock_DmViewerBlock"]/div[2]/div[1]/div[%d]' % (x)
                            self.driver.find_element(By.XPATH, ticketNumXpath).click()
                        time.sleep(0.5)
                        # 下单
                        print('正在下单.......')
                        self.driver.find_element(By.XPATH, '//*[@id="dmOrderSubmitBlock_DmOrderSubmitBlock"]/div[2]/div[1]/div[2]/div[3]/div[2]').click()
                        # 如果当前购票人信息存在 就点击
                        # if self.isElementExist('//*[@id="container"]/div/div[9]/button'):
                            # 下单操作
                        # self.check_order()
                        time.sleep(20)
                        return True

    def choice_seats(self):
        """选择座位"""
        while self.driver.title == '选座购买':
            while self.isElementExist('//*[@id="app"]/div[2]/div[2]/div[1]/div[2]/img'):
                print('请快速选择你想要的座位!!!')
            while self.isElementExist('//*[@id="app"]/div[2]/div[2]/div[2]/div'):
                self.driver.find_element(By.XPATH, '//*[@id="app"]/div[2]/div[2]/div[2]/button').click()

    def check_order(self):
        """下单操作"""
        if self.status in [3, 4, 5]:
            print('###开始确认订单###')
            time.sleep(1)
            try:
                # 默认选第一个购票人信息
                self.driver.find_element(By.XPATH, '//*[@id="container"]/div/div[2]/div[2]/div[1]/div/label').click()
            except Exception as e:
                print('###购票人信息选中失败, 自行查看元素位置###')
                print(e)
            # 最后一步提交订单
            time.sleep(0.5)  # 太快了不好, 影响加载 导致按钮点击无效
            self.driver.find_element(By.XPATH, '//*[@id="container"]/div/div[9]/button').click()
            time.sleep(20)

    def isElementExist(self, element):
        """判断元素是否存在"""
        flag = True
        browser = self.driver
        try:
            browser.find_element(By.XPATH, element)
            return flag
        except:
            flag = False
            return flag

    def finish(self):
        """抢票完成, 退出"""
        self.driver.quit()


if __name__ == '__main__':
    con = Concert()
    try:
        # 从命令行中获取对应的抢票详情页面
        parser = argparse.ArgumentParser()
        parser.add_argument("-a", "--inputA", help="请输入你需要抢票的详情页面地址：", dest="argA", type=str, default="xxx")
        parser.add_argument("-b", "--inputB", help="请输入你需要抢票的城市", dest="argB", type=int, default="1")
        parser.add_argument("-c", "--inputC", help="请输入你需要抢票的场次", dest="argC", type=int, default="1")
        parser.add_argument("-d", "--inputD", help="请输入你需要抢票的票档", dest="argD", type=int, default="1")
        parser.add_argument("-e", "--inputE", help="请输入你需要抢票的数量", dest="argE", type=int, default="1")
        args = parser.parse_args()
        target_url = args.argA
        city = args.argB
        sessions = args.argC
        ticket_stalls = args.argD
        ticket_num = args.argE
        print('初始化详情、城市、场次、票档、数量：', target_url, city, sessions, ticket_stalls, ticket_num)
        con.enter_concert()  # 打开浏览器
        # 下单选票，如果失败则等1s后刷新页面重新选
        if con.choose_ticket() == False:
            time.sleep(1)
            con.driver.refresh()
            con.choose_ticket()
        con.finish()    
    except Exception as e:
        print(e)
        con.finish()

