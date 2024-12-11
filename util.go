package socks5

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// ParseAddress format address x.x.x.x:xx to raw address.
// addr contains domain length
func ParseAddress(address string) (a byte, addr []byte, port []byte, err error) {
	var h, p string
	h, p, err = net.SplitHostPort(address)
	if err != nil {
		return
	}
	ip := net.ParseIP(h)
	if ip4 := ip.To4(); ip4 != nil {
		a = ATYPIPv4
		addr = []byte(ip4)
	} else if ip6 := ip.To16(); ip6 != nil {
		a = ATYPIPv6
		addr = []byte(ip6)
	} else {
		a = ATYPDomain
		addr = []byte{byte(len(h))}
		addr = append(addr, []byte(h)...)
	}
	i, _ := strconv.Atoi(p)
	port = make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(i))
	return
}

// bytes to address
// addr contains domain length
func ParseBytesAddress(b []byte) (a byte, addr []byte, port []byte, err error) {
	if len(b) < 1 {
		err = errors.New("Invalid address")
		return
	}
	a = b[0]
	if a == ATYPIPv4 {
		if len(b) < 1+4+2 {
			err = errors.New("Invalid address")
			return
		}
		addr = b[1 : 1+4]
		port = b[1+4 : 1+4+2]
		return
	}
	if a == ATYPIPv6 {
		if len(b) < 1+16+2 {
			err = errors.New("Invalid address")
			return
		}
		addr = b[1 : 1+16]
		port = b[1+16 : 1+16+2]
		return
	}
	if a == ATYPDomain {
		if len(b) < 1+1 {
			err = errors.New("Invalid address")
			return
		}
		l := int(b[1])
		if len(b) < 1+1+l+2 {
			err = errors.New("Invalid address")
			return
		}
		addr = b[1 : 1+1+l]
		port = b[1+1+l : 1+1+l+2]
		return
	}
	err = errors.New("Invalid address")
	return
}

// ToAddress format raw address to x.x.x.x:xx
// addr contains domain length
func ToAddress(a byte, addr []byte, port []byte) string {
	var h, p string
	if a == ATYPIPv4 || a == ATYPIPv6 {
		h = net.IP(addr).String()
	}
	if a == ATYPDomain {
		if len(addr) < 1 {
			return ""
		}
		if len(addr) < int(addr[0])+1 {
			return ""
		}
		h = string(addr[1:])
	}
	p = strconv.Itoa(int(binary.BigEndian.Uint16(port)))
	return net.JoinHostPort(h, p)
}

// Address return request address like ip:xx
func (r *Request) Address() string {
	var s string
	if r.Atyp == ATYPDomain {
		s = bytes.NewBuffer(r.DstAddr[1:]).String()
	} else {
		s = net.IP(r.DstAddr).String()
	}
	p := strconv.Itoa(int(binary.BigEndian.Uint16(r.DstPort)))
	return net.JoinHostPort(s, p)
}

// Address return request address like ip:xx
func (r *Reply) Address() string {
	var s string
	if r.Atyp == ATYPDomain {
		s = bytes.NewBuffer(r.BndAddr[1:]).String()
	} else {
		s = net.IP(r.BndAddr).String()
	}
	p := strconv.Itoa(int(binary.BigEndian.Uint16(r.BndPort)))
	return net.JoinHostPort(s, p)
}

// Address return datagram address like ip:xx
func (d *Datagram) Address() string {
	var s string
	if d.Atyp == ATYPDomain {
		s = bytes.NewBuffer(d.DstAddr[1:]).String()
	} else {
		s = net.IP(d.DstAddr).String()
	}
	p := strconv.Itoa(int(binary.BigEndian.Uint16(d.DstPort)))
	return net.JoinHostPort(s, p)
}

func GetRandomIPFromCidrs(cidrs []string) (string, error) {
	if len(cidrs) == 0 {
		return "", fmt.Errorf("no CIDRs provided")
	}

	// 初始化随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 随机选择一个CIDR
	selectedCIDR := cidrs[r.Intn(len(cidrs))]

	// 解析选定的CIDR
	_, ipNet, err := net.ParseCIDR(selectedCIDR)
	if err != nil {
		return "", fmt.Errorf("invalid CIDR: %v", err)
	}

	// 获取CIDR的掩码和网络地址
	mask := ipNet.Mask
	networkIP := ipNet.IP

	// 计算可用IP地址的数量
	ones, bits := mask.Size()
	maxHosts := new(big.Int).Lsh(big.NewInt(1), uint(bits-ones))
	maxHosts.Sub(maxHosts, big.NewInt(1))

	// 生成随机数
	randomInt := new(big.Int)
	randomInt.Rand(r, maxHosts)

	// 计算随机IP
	randomIP := make(net.IP, len(networkIP))
	copy(randomIP, networkIP)

	// 将随机数添加到网络地址以获得最终IP
	for i := len(randomIP) - 1; i >= 0; i-- {
		randomIP[i] |= byte(randomInt.Int64() & 0xFF)
		randomInt.Rsh(randomInt, 8)
	}

	// 检查生成的IP是否在CIDR范围内
	if !ipNet.Contains(randomIP) {
		return "", fmt.Errorf("error, generated IP is not in CIDR range")
	}

	// 如果是IPv6地址，添加方括号
	if randomIP.To4() == nil {
		return fmt.Sprintf("[%s]", randomIP.String()), nil
	}

	return randomIP.String(), nil
}
