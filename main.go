package main

import (
	"JCRandomProxy/Conf"
	"JCRandomProxy/Proxy"
	"log"
	"net"
	"runtime/debug"
)

/**
*
* Author: JC0o0l
* email: jerryzvs@163.com
* wechat: JC_SecNotes
 */

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Conf.InitConfig()
}

func main() {
	// 监听TCP连接
	l, err := net.Listen("tcp", ":"+Conf.Port)
	if err != nil {
		log.Panic(err)
	}

	for {
		// 接收TCP连接，返回一个net.Conn
		client, err := l.Accept()
		if err != nil {
			log.Panic("Panic", err)
		}
		// 收到请求后，调用handle进行处理
		go handle(client)
	}
}

func handle(client net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Panic(err)
			debug.PrintStack()
		}
	}()
	if client == nil {
		return
	}

	log.Println("JCTLog: client tcp tunnel connection: ", client.LocalAddr().String(), "->", client.RemoteAddr().String())
	defer client.Close()
	// 随机取出一个代理
	paddr, ptype, _ := Proxy.GetAProxy()
	proxyAddr := ptype + "://" + paddr
	// 验证代理是否有效
	checkaddr := "http://myip.ipip.net"
	if Proxy.CheckProxy(proxyAddr, checkaddr) {
		log.Println(" 代理有效 ",proxyAddr)
		// 有效，使用端口转发
		PortForward(client,paddr)
	}else{
		log.Println(" 代理无效 ",proxyAddr)
		// 无效，使用自身代理
		lproxy(client)
	}
}

