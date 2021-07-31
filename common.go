package common

import (
	"bufio"
	"errors"
	"io"
	"net"
	"os"
	"reflect"
)

func IsFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}

func IsDir(p string) bool {
	fi, e := os.Stat(p)
	if e != nil {
		return false
	}
	return fi.IsDir()
}

// 判断接口是否初始化 如果未初始化，返回true
func IsStructureEmpty(x, y interface{}) bool {
	return reflect.DeepEqual(x, y)
}

// 原文件名称:srcPath
// 新文件名称:dstPath
func CopyFile(srcPath, dstPath string) error {
	//打开文件
	srcFile, err := os.OpenFile(srcPath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	//关闭文件句柄
	defer srcFile.Close()
	defer dstFile.Close()

	//将文件写入缓冲区
	srcBuf := bufio.NewReader(srcFile)
	dstBuf := bufio.NewWriter(dstFile)
	//调用系统复制方法
	_, err = io.Copy(dstBuf, srcBuf)
	if err != nil {
		return err
	}
	//落盘
	err = dstBuf.Flush()
	if err != nil {
		return err
	}
	return nil
}

// 判断文件夹或者文件是否存在
func IsFileExists(file string) error {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return err
	}
	return nil
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}
	return ip
}

// 获取本机ip
func GetLocalIp() (net.IP, error) {
	ip, err := externalIP()
	if err != nil {
		return nil, err
	}
	return ip, nil
}

func GetLocalIpStr() string {
	ipnet, _ := GetLocalIp()
	return ipnet.String()
}
