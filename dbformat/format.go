package dbformat

import (
	"encoding/binary"
	"net"
	"sort"
	"strings"
)

type IpRangeItem struct {
	Origin    string
	LowU32    uint32
	HighU32   uint32
	Attach    string
	AttachObj IpRangeAttach
	//CityId  uint32
}

type IpRangeAttach struct {
	Country  string // 国家
	Province string // 省
	City     string // 市
	ISP      string // 网络提供商
}

type DbFormatType struct {
	ShowPriority int    // 显示优先级
	NameForCmd   string // 命令行使用
	Desc         string // 描述
	ExtName      string // 扩展名
	SupportWrite bool   // 是否支持生成文件
}

const (
	ShowPriority_txt           = iota
	ShowPriority_Linsoul2014v1 = iota
	ShowPriority_Linsoul2014v2 = iota
	ShowPriority_IpipdnetIpdb  = iota
	ShowPriority_Maxmind       = iota
)

type DBFormat interface {
	GetType() DbFormatType
	ReadData(data []byte) (list []IpRangeItem, err error)
	FormatAttach(attach IpRangeAttach) (value string, err error)
	WriteData(list []IpRangeItem) (data []byte, err error)
}

type DbNeedVerifyFiled7 interface {
	NeedVerifyFiled7() bool
}

var gDbFormatList []DBFormat

func RegisterDbFormat(format DBFormat) {
	gDbFormatList = append(gDbFormatList, format)

	sort.Slice(gDbFormatList, func(i, j int) bool {
		a, b := gDbFormatList[i], gDbFormatList[j]
		return a.GetType().ShowPriority < b.GetType().ShowPriority
	})
}

func GetDbFormatByType(desc string) DBFormat {
	desc = strings.ToLower(desc)

	for _, one := range gDbFormatList {
		t := one.GetType()
		if strings.ToLower(t.Desc) == desc || strings.ToLower(t.NameForCmd) == desc {
			return one
		}
	}
	return nil
}

func GetDbTypeList() (list []DbFormatType) {
	for _, one := range gDbFormatList {
		list = append(list, one.GetType())
	}
	return list
}

func Uint32ToIpv4(ip uint32) net.IP {
	var tmp = make([]byte, 4)
	binary.BigEndian.PutUint32(tmp, ip)
	return net.IPv4(tmp[0], tmp[1], tmp[2], tmp[3])
}

func Ipv4ToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return binary.BigEndian.Uint32([]byte{ip[0], ip[1], ip[2], ip[3]})
}
