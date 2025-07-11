package ipipdnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/orestonce/Ip2regionTool/dbformat"
	"net"
)

type DbFormatIpipdnetData struct {
}

func (d DbFormatIpipdnetData) GetType() dbformat.DbFormatType {
	return dbformat.DbFormatType{
		ShowPriority: dbformat.ShowPriority_IpipdnetDat,
		NameForCmd:   "IpipdnetDat",
		Desc:         "IpipdnetDat, ipip.net的dat格式数据",
		ExtName:      "*.dat",
		SupportWrite: false,
	}
}

func (d DbFormatIpipdnetData) ReadData(data []byte) (list []dbformat.IpRangeItem, err error) {
	loc := new(Datax_locator)
	loc.init(data)
	return readDataWith(loc)
}

func (d DbFormatIpipdnetData) FormatAttach(attach dbformat.IpRangeAttach) (value string, err error) {
	return "", errors.New("implement me")
}

func (d DbFormatIpipdnetData) WriteData(list []dbformat.IpRangeItem) (data []byte, err error) {
	return nil, errors.New("implement me")
}

type DbFormatIpipdnetDataX struct {
}

func (d DbFormatIpipdnetDataX) GetType() dbformat.DbFormatType {
	return dbformat.DbFormatType{
		ShowPriority: dbformat.ShowPriority_IpipdnetDatX,
		NameForCmd:   "IpipdnetDatX",
		Desc:         "IpipdnetDatX, ipip.net的datx格式数据",
		ExtName:      "*.datx",
		SupportWrite: false,
	}
}

func (d DbFormatIpipdnetDataX) ReadData(data []byte) (list []dbformat.IpRangeItem, err error) {
	loc := new(Datax_locator)
	loc.initX(data)
	return readDataWith(loc)
}

func readDataWith(loc *Datax_locator) (list []dbformat.IpRangeItem, err error) {
	rangeList, localList := loc.Dump()
	if len(rangeList) != len(localList) {
		return nil, fmt.Errorf("rangeList (%v), localList (%v)", len(rangeList), len(localList))
	}

	for idx := 0; idx < len(rangeList); idx++ {
		local := localList[idx]

		list = append(list, dbformat.IpRangeItem{
			Origin:  "",
			LowU32:  dbformat.Ipv4ToUint32(rangeList[idx].Start),
			HighU32: dbformat.Ipv4ToUint32(rangeList[idx].End),
			Attach:  fmt.Sprintf("%v|0|%v|%v|%v", local.Country, local.City, local.Region, local.Isp),
			//CityId:  0,
			AttachObj: dbformat.IpRangeAttach{
				Country:  local.Country,
				Province: local.City,
				City:     local.Region,
				ISP:      local.Isp,
			},
		})
	}
	return list, nil
}

func (d DbFormatIpipdnetDataX) FormatAttach(attach dbformat.IpRangeAttach) (value string, err error) {
	return "", errors.New("implement me")
}

func (d DbFormatIpipdnetDataX) WriteData(list []dbformat.IpRangeItem) (data []byte, err error) {
	return nil, errors.New("implement me")
}

func init() {
	dbformat.RegisterDbFormat(DbFormatIpipdnetData{})
	dbformat.RegisterDbFormat(DbFormatIpipdnetDataX{})
}

type Datax_locator struct {
	index           [256]int
	indexData       []uint32
	textStartIndex  []int
	textLengthIndex []int
	textData        []byte
}

type Datax_Range struct {
	Start net.IP
	End   net.IP
}

// Find locationInfo by ip string
// It will return err when ipstr is not a valid format
func (loc *Datax_locator) Find(ipstr string) (info *LocationInfo, err error) {
	ip := net.ParseIP(ipstr)
	if ip == nil || ip.To4() == nil {
		err = ErrUnsupportedIP
		return
	}
	info = loc.FindByUint(binary.BigEndian.Uint32(ip.To4()))
	return
}

// Find locationInfo by uint32
func (loc *Datax_locator) FindByUint(ip uint32) (info *LocationInfo) {

	idx := loc.findTextIndex(ip, loc.index[ip>>24])
	start := loc.textStartIndex[idx]
	return newLocationInfo(loc.textData[start : start+loc.textLengthIndex[idx]])
}

func (loc *Datax_locator) Dump() (rs []Datax_Range, locs []*LocationInfo) {

	rs = make([]Datax_Range, 0, len(loc.indexData))
	locs = make([]*LocationInfo, 0, len(loc.indexData))

	for i := 1; i < len(loc.indexData); i++ {
		s, e := loc.indexData[i-1], loc.indexData[i]
		off := loc.textStartIndex[i]
		l := newLocationInfo(loc.textData[off : off+loc.textLengthIndex[i]])
		rs = append(rs, Datax_Range{ipOf(s), ipOf(e)})
		locs = append(locs, l)
	}
	return
}

func ipOf(n uint32) net.IP {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return net.IP(b)
}

// binary search
func (loc *Datax_locator) findTextIndex(ip uint32, start int) int {

	end := len(loc.indexData) - 1
	for start < end {
		mid := (start + end) / 2
		if ip > loc.indexData[mid] {
			start = mid + 1
		} else {
			end = mid
		}
	}

	if loc.indexData[end] >= ip {
		return end
	} else {
		return start
	}

}

func (loc *Datax_locator) init(data []byte) {

	offset := int(binary.BigEndian.Uint32(data[:4]))
	textOff := offset - 1024

	loc.textData = data[textOff:]
	for i := 0; i < 256; i++ {
		off := 4 + i*4
		loc.index[i] = int(binary.LittleEndian.Uint32(data[off : off+4]))
	}

	nidx := (textOff - 4 - 1024) / 8

	loc.indexData = make([]uint32, nidx)
	loc.textStartIndex = make([]int, nidx)
	loc.textLengthIndex = make([]int, nidx)

	for i := 0; i < nidx; i++ {
		off := 4 + 1024 + i*8
		loc.indexData[i] = binary.BigEndian.Uint32(data[off : off+4])
		loc.textStartIndex[i] = int(uint32(data[off+4]) | uint32(data[off+5])<<8 | uint32(data[off+6])<<16)
		loc.textLengthIndex[i] = int(data[off+7])
	}
	return
}

// datx format
func (loc *Datax_locator) initX(data []byte) {

	offset := int(binary.BigEndian.Uint32(data[:4]))
	textOff := offset - 256*256*4
	loc.textData = data[textOff:]
	for i := 0; i < 256; i++ {
		// datx格式使用了ipv4的前两个字节做为索引字段，出于对data格式兼容考虑这里只使用首字节做为索引字段
		// 由于我们使用二分查找, 这个点上认为对性能不会有任何影响
		off := 4 + i*256*4
		loc.index[i] = int(binary.LittleEndian.Uint32(data[off : off+4]))
	}

	nidx := (textOff - 4 - 256*256*4) / 9

	loc.indexData = make([]uint32, nidx)
	loc.textStartIndex = make([]int, nidx)
	loc.textLengthIndex = make([]int, nidx)

	for i := 0; i < nidx; i++ {
		off := 4 + 256*256*4 + i*9
		loc.indexData[i] = binary.BigEndian.Uint32(data[off : off+4])
		loc.textStartIndex[i] = int(uint32(data[off+4]) | uint32(data[off+5])<<8 | uint32(data[off+6])<<16)
		loc.textLengthIndex[i] = int(uint32(data[off+8]) | uint32(data[off+7])<<8)
	}
	return
}

func newLocationInfo(str []byte) *LocationInfo {

	var info *LocationInfo

	fields := bytes.Split(str, []byte("\t"))
	if len(fields) < 4 {
		panic("unexpected ip info:" + string(str))
	}
	info = &LocationInfo{
		Country: string(fields[0]),
		Region:  string(fields[1]),
		City:    string(fields[2]),
	}
	if len(fields) >= 5 {
		info.Isp = string(fields[4])
	}

	{
		if len(info.Country) == 0 {
			info.Country = Null
		}
		if len(info.Region) == 0 {
			info.Region = Null
		}
		if len(info.City) == 0 {
			info.City = Null
		}
		if len(info.Isp) == 0 {
			info.Isp = Null
		}
	}

	return info
}
