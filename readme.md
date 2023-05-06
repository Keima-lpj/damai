# 这是一个大麦网抢票的脚本

**分为python实现和golang实现。推荐使用python版本，golang目前存在问题**

**目前暂时只支持chrome，需根据需要下载对应的chromedriver**

```
    https://chromedriver.chromium.org/downloads
```

**使用方法**

``` shell
    # -a 表示抢票详情页
    # -b 表示城市，为数字，从左往右从上往下。抢票前请先去详情页确认，下面的参数同上
    # -c 表示场次
    # -d 表示票档
    # -e 表示数量。注意，选择了数量之后，当前账户录入的观影人数量需大于等于此数量，否则会导致失败
    # -f 表示使用浏览器的driver
    python3 damai.py -a "https://detail.damai.cn/item.htm?id=711954279503" -b 10 -c 2 -d 1 -e 1 -f "/home/keima/Desktop/chromedriver/chromedriver"
```
**注意，如果登录保存的cookie文件失效，需要删除本地的cookies.pkl文件**




