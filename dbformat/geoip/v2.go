package geoip

import (
	"errors"
	"github.com/orestonce/Ip2regionTool/dbformat"
	"github.com/oschwald/geoip2-golang"
	"github.com/oschwald/maxminddb-golang/v2"
	"net"
	"net/netip"
)

type DbFormatGeoip struct {
}

func (d DbFormatGeoip) GetType() dbformat.DbFormatType {
	return dbformat.DbFormatType{
		ShowPriority: dbformat.ShowPriority_Maxmind,
		NameForCmd:   "Maxmind",
		Desc:         "Maxmind, maxmind的mmdb格式数据",
		ExtName:      "*.mmdb",
		SupportWrite: false,
	}
}

func DecodeResult(result maxminddb.Result) (attachObj dbformat.IpRangeAttach, err error) {
	var city geoip2.City
	var country geoip2.Country
	var asn geoip2.ASN

	err1 := result.Decode(&country)
	if err1 != nil && attachObj.Country == "" {
		attachObj.Country = country.RegisteredCountry.Names["zh-CN"]
	}
	err2 := result.Decode(&city)
	for _, lang := range []string{"zh-CN", "en"} {
		for _, names := range []map[string]string{country.RegisteredCountry.Names, country.Country.Names, city.RegisteredCountry.Names, city.Country.Names} {
			if attachObj.Country != "" {
				break
			}
			attachObj.Country = names[lang]
		}
		if attachObj.Country != "" {
			break
		}
	}
	for _, lang := range []string{"zh-CN", "en"} {
		for _, subdivision := range city.Subdivisions {
			if attachObj.Province != "" {
				break
			}
			attachObj.Province = subdivision.Names[lang]
		}
		if attachObj.Province != "" {
			break
		}
	}
	for _, lang := range []string{"zh-CN", "en"} {
		if attachObj.City != "" {
			break
		}
		attachObj.City = city.City.Names[lang]
	}

	err3 := result.Decode(&asn)
	if err3 == nil && attachObj.ISP == "" {
		attachObj.ISP = asn.AutonomousSystemOrganization
	}

	if err1 != nil && err2 != nil && err3 != nil {
		return attachObj, errors.New("result.Decode error " + result.Prefix().String() + " " + err1.Error())
	}
	return attachObj, nil
}

func (DbFormatGeoip) ReadData(data []byte) (list []dbformat.IpRangeItem, err error) {
	var db *maxminddb.Reader
	db, err = maxminddb.FromBytes(data)
	if err != nil {
		return nil, err
	}

	var result maxminddb.Result
	for result = range db.Networks() {
		network := result.Prefix()
		if network.Addr().Is4() == false {
			continue
		}
		var attachObj dbformat.IpRangeAttach
		attachObj, err = DecodeResult(result)
		if err != nil {
			return nil, err
		}
		list = append(list, dbformat.IpRangeItem{
			Origin:    "",
			LowU32:    dbformat.Ipv4ToUint32(FirstIP(network)),
			HighU32:   dbformat.Ipv4ToUint32(LastIP(network)),
			Attach:    "",
			AttachObj: attachObj,
		})
	}

	return list, nil
}

func (d DbFormatGeoip) FormatAttach(attach dbformat.IpRangeAttach) (value string, err error) {
	return "", errors.New("implement me")
}

func (d DbFormatGeoip) WriteData(list []dbformat.IpRangeItem) (data []byte, err error) {
	return nil, errors.New("implement me")
}

func init() {
	dbformat.RegisterDbFormat(DbFormatGeoip{})
}

// FirstIP 返回 Prefix 中的第一个 IP 地址（网络地址）
func FirstIP(prefix netip.Prefix) net.IP {
	return net.ParseIP(prefix.Addr().String())
}

// LastIP 返回 Prefix 中的最后一个 IP 地址（广播地址或网络末尾地址）
func LastIP(prefix netip.Prefix) net.IP {
	addr := prefix.Addr()
	bits := prefix.Bits()

	// 计算主机位数（总位数 - 网络位数）
	var hostBits int
	if addr.Is4() {
		hostBits = 32 - bits
	} else { // IPv6
		hostBits = 128 - bits
	}

	// 如果主机位数为 0（如 /32, /128），直接返回网络地址
	if hostBits == 0 {
		return net.ParseIP(addr.String())
	}

	// 计算需要添加的偏移量（2^hostBits - 1）
	maxHosts := uint64(1<<hostBits) - 1

	// 处理 IPv4（32 位）
	if addr.Is4() {
		// 将 IP 转换为 32 位无符号整数
		ip32 := addr.As4()
		var ipUint uint32
		for i := 0; i < 4; i++ {
			ipUint = (ipUint << 8) | uint32(ip32[i])
		}

		// 加上偏移量
		ipUint += uint32(maxHosts)

		// 转换回 netip.Addr
		var result [4]byte
		for i := 3; i >= 0; i-- {
			result[i] = byte(ipUint & 0xFF)
			ipUint >>= 8
		}
		return net.IPv4(result[0], result[1], result[2], result[3])
	}

	return net.ParseIP(addr.String()) // ipv6不考虑

	//// 处理 IPv6（128 位）
	//ip128 := addr.As16()
	//// 注意：对于大范围的 IPv6 地址（如 /64），直接计算可能不实用
	//// 这里仅处理较小的范围
	//for i := uint64(0); i < maxHosts; i++ {
	//	// 逐位加 1（简化实现，实际应使用更高效的大数加法）
	//	for j := 15; j >= 0; j-- {
	//		if ip128[j] < 255 {
	//			ip128[j]++
	//			break
	//		}
	//		ip128[j] = 0
	//	}
	//}
	//return netip.AddrFrom16(ip128)
}
