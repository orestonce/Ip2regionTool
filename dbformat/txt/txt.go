package txt

import (
	"bytes"
	"github.com/orestonce/Ip2regionTool/dbformat"
	"net"
	"strings"
)

type DBFormatTxt struct {
}

func (DBFormatTxt) GetType() dbformat.DbFormatType {
	return dbformat.DbFormatType{
		ShowPriority: dbformat.ShowPriority_txt,
		NameForCmd:   "txt",
		Desc:         "txt,纯文本格式",
		ExtName:      "*.txt",
		SupportWrite: true,
	}
}

func (DBFormatTxt) ReadData(data []byte) (list []dbformat.IpRangeItem, err error) {
	for _, one := range strings.Split(string(data), "\n") {
		one = strings.TrimSpace(one)
		if one == `` {
			continue
		}
		temp := strings.Split(one, `|`)
		if len(temp) < 2 {
			continue
		}
		sip := dbformat.Ipv4ToUint32(net.ParseIP(temp[0]))
		eip := dbformat.Ipv4ToUint32(net.ParseIP(temp[1]))
		list = append(list, dbformat.IpRangeItem{
			Origin:  one,
			LowU32:  sip,
			HighU32: eip,
			Attach:  strings.Join(temp[2:], `|`),
		})
	}
	return list, nil
}

func (d DBFormatTxt) FormatAttach(attach dbformat.IpRangeAttach) (value string, err error) {
	return strings.Join([]string{
		attach.Country,
		"0",
		attach.Province,
		attach.City,
		attach.ISP,
	}, "|"), nil
}

func (DBFormatTxt) WriteData(list []dbformat.IpRangeItem) (data []byte, err error) {
	buf := bytes.NewBuffer(nil)
	for _, one := range list {
		buf.WriteString(dbformat.Uint32ToIpv4(one.LowU32).String() + `|` + dbformat.Uint32ToIpv4(one.HighU32).String() + `|` + one.Attach + "\n")
	}
	return buf.Bytes(), nil
}

func init() {
	dbformat.RegisterDbFormat(DBFormatTxt{})
}
