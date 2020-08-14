package util

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/andycai/void"
)

// 将普通地址格式(host:port)拆分
func SplitAddress(addr string) (host string, port int, err error) {

	var portStr string

	host, portStr, err = net.SplitHostPort(addr)

	if err != nil {
		return "", 0, err
	}

	port, err = strconv.Atoi(portStr)

	if err != nil {
		return "", 0, err
	}

	return
}

// 将host和端口合并为(host:port)格式的地址
func JoinAddress(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

// 修复ws没有实现所有net.conn方法，导致无法获取客服端地址问题.
type RemoteAddr interface {
	RemoteAddr() net.Addr
}

// 获取session远程的地址
func GetRemoteAddress(ses void.Session) (string, bool) {
	if c, ok := ses.Raw().(RemoteAddr); ok {
		return c.RemoteAddr().String(), true
	}

	return "", false
}

var (
	ErrInvalidPortRange = errors.New("invalid port range")
)

type Addr struct {
	Scheme  string
	Host    string
	MinPort int
	MaxPort int
	Path    string
}

// HostPortString return host:port
func (a *Addr) HostPortString(port int) string {
	return fmt.Sprintf("%s:%d", a.Host, port)
}

// String return scheme://host:port/path
func (a *Addr) String(port int) string {
	if a.Scheme == "" {
		return a.HostPortString(port)
	}

	return fmt.Sprintf("%s://%s:%d%s", a.Scheme, a.Host, port, a.Path)
}

// ParseAddr format: scheme://host:minPort~maxPort/path
func ParseAddr(addrStr string) (addr *Addr, err error) {
	addr = new(Addr)

	schemePos := strings.Index(addrStr, "://")

	// 移除scheme部分
	if schemePos != -1 {
		addr.Scheme = addrStr[:schemePos]
		addrStr = addrStr[schemePos+3:]
	}

	colonPos := strings.Index(addrStr, ":")

	if colonPos != -1 {
		addr.Host = addrStr[:colonPos]
	}

	addrStr = addrStr[colonPos+1:]

	rangePos := strings.Index(addrStr, "~")

	var minStr, maxStr string
	if rangePos != -1 {
		minStr = addrStr[:rangePos]

		slashPos := strings.Index(addrStr, "/")

		if slashPos != -1 {
			maxStr = addrStr[rangePos+1 : slashPos]
			addr.Path = addrStr[slashPos:]
		} else {
			maxStr = addrStr[rangePos+1:]
		}
	} else {
		slashPos := strings.Index(addrStr, "/")

		if slashPos != -1 {
			addr.Path = addrStr[slashPos:]
			minStr = addrStr[rangePos+1 : slashPos]
		} else {
			minStr = addrStr[rangePos+1:]
		}
	}

	// extract min port
	addr.MinPort, err = strconv.Atoi(minStr)
	if err != nil {
		return nil, ErrInvalidPortRange
	}

	if maxStr != "" {
		// extract max port
		addr.MaxPort, err = strconv.Atoi(maxStr)
		if err != nil {
			return nil, ErrInvalidPortRange
		}
	} else {
		addr.MaxPort = addr.MinPort
	}

	return
}

func DetectPort(addrStr string, fn func(addr *Addr, port int) (interface{}, error)) (interface{}, error) {
	addr, err := ParseAddr(addrStr)
	if err != nil {
		return nil, err
	}

	for port := addr.MinPort; port <= addr.MaxPort; port++ {
		listener, err := fn(addr, port)
		if err == nil {
			return listener, nil
		}

		if port == addr.MaxPort {
			return nil, err
		}
	}

	return nil, fmt.Errorf("unable to bind to %s", addrStr)
}
