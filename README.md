代理封装 edge浏览器tts 接口 和百度 tts接口
主要用于 阅读 软件听书自用

# 如何使用
## 运行
```bash
go run main.go --port 2580 --host 0.0.0.0 --ct baidu
```
## 阅读配置
```js
http://your.host:2580,{
    "method": "POST",
    "body": "tex={{java.encodeURI(java.encodeURI(speakText))}}&spd={{(speakSpeed + 5) / 5+ 2}}&vn=5003&v=10"
}
```



