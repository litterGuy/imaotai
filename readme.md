# imaotai自动预约

app碰到无法预约，预约一直失败要求重试。所以有了想找一个预约脚本的想法。

偶然间发现了这个项目[https://github.com/oddfar/campus-imaotai](https://github.com/oddfar/campus-imaotai)
，实现了自己的需求。但是硬件配置有限，每次启动程序太吃资源了、于是就有了这个项目。

## 项目介绍

- 用户信息配置到yml文件
- 使用sqlite存储预约需要的信息
- 每天9点20定时任务去跑预约

## 填写配置文件

1. 获取所在地经纬度

```
imaotai --address *****
```

根据返回的结果值，将其填写到config.yml文件

2. 发送登录短信验证码

```
imaotai --codephone *****
```

执行成功后进行登录

3. 登录获取用户信息

```
imaotai --phone ***** --code ****
```

根据返回的结果值，将其填写到config.yml文件

4. reserveType 设置，1-预约本市出货量最大的门店 2-预约你的位置附近门店

至此，配置文件处理完毕。

## 启动应用

```
imaotai 
```

如果无异常，可等待预约完成后去app查看结果

## 打包linux

未解决windows下交叉编译的问题。因此增加dockerfile曲线处理一下。

1. 打包镜像

```
docker build -t imaotai .
```

2. 执行命令操作

```
// 获取地址经纬度
docker run --name imaotai imaotai --address "清华大学"

// 发送登录短信验证码
docker run --name imaotai imaotai --codephone *****

// 登录获取用户信息
docker run --name imaotai imaotai --phone ***** --code ****
```

3. 配置config.yml文件
4. 启动容器

```
docker run --name imaotai -d -v ${config.yml}:/app/config.yml imaotai
```

${config.yml} 为本地的配置文件

5. 把编译文件复制到本地

```
docker cp imaotai:/app/imaotai .
```

## 增加pushplus推送

去[官网](http://www.pushplus.plus/push2.html)获取到token和topic后，填入到配置文件


## 更新
增加配置参数crossCity.
```
0-跨市 1-不跨市
```

### dockerfile 在阿里云构建过慢
在dockerfile添加
```
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories 
```
替换成阿里云资源

## 未添加事项

- [ ] 没有指定哪种酒预约。目前是预约查询到的全部
- [ ] 暂未处理好错误打印信息，没有将错误log记录文件
- [ ] 尚未做好测试，目前是否存在其他问题暂不知道
- [x] 增加类似plusplus这种消息推送
- [x] windows打包linux失败