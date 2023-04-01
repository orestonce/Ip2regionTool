package Ip2regionTool

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"io/ioutil"
	"math"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ConvertDbToTxt_Req struct {
	DbFileName  string
	TxtFileName string
	Merge       bool
	DbVersion   int
}

func GetDbVersionByName(name string) int {
	if strings.HasSuffix(strings.ToLower(name), ".db") {
		return 1
	}
	if strings.HasSuffix(strings.ToLower(name), ".xdb") {
		return 2
	}
	return 1
}

func ConvertDbToTxt(req ConvertDbToTxt_Req) (errMsg string) {
	if req.DbVersion == 0 {
		req.DbVersion = GetDbVersionByName(req.DbFileName)
	}
	stat, err := os.Stat(req.DbFileName)
	if err != nil {
		return "文件状态错误: " + req.DbFileName + "," + err.Error()
	}
	if stat.Size() > 1000*1024*1024 {
		return "不支持超过1000MB的db文件: " + strconv.Itoa(int(stat.Size()))
	}
	dbFileContent, err := ioutil.ReadFile(req.DbFileName)
	if err != nil {
		return "读取db文件失败: " + req.DbFileName + ", " + err.Error()
	}
	var list []IpRangeItem
	if req.DbVersion == 1 {
		list, errMsg = ReadV1DataBlob(dbFileContent)
	} else {
		list, errMsg = ReadV2DataBlob(dbFileContent)
	}
	if errMsg != `` {
		return "文件数据错误: " + errMsg
	}
	if req.Merge {
		list = MergeIpRangeList(list)
	}
	errMsg = VerifyIpRangeList(VerifyIpRangeListRequest{
		DataInfoList:     list,
		VerifyFullUint32: true,
		VerifyFiled7:     req.DbVersion == 1, // 只有版本1才需要验证字段数为7
	})
	if errMsg != `` {
		return "验证文件数据失败: " + errMsg
	}
	data := WriteV1DataTxt(list)
	err = ioutil.WriteFile(req.TxtFileName, data, 0777)
	if err != nil {
		return "输出文件写入失败: " + err.Error()
	}
	return ""
}

type ConvertTxtToDb_Req struct {
	TxtFileName       string
	DbFileName        string
	RegionCsvFileName string
	Merge             bool
}

func ConvertTxtToDb(req ConvertTxtToDb_Req) (errMsg string) {
	stat, err := os.Stat(req.TxtFileName)
	if err != nil {
		return "文件状态错误: " + req.DbFileName + "," + err.Error()
	}
	if stat.Size() > 1000*1024*1024 {
		return "不支持超过1000MB的db文件: " + strconv.Itoa(int(stat.Size()))
	}
	var globalRegionMap map[string]uint32
	if req.RegionCsvFileName != "" {
		globalRegionMap, errMsg = ReadGlobalRegionMap(req.RegionCsvFileName)
		if errMsg != "" {
			return errMsg
		}
	}
	txtFileContent, err := ioutil.ReadFile(req.TxtFileName)
	if err != nil {
		return "读取db文件失败: " + req.TxtFileName + ", " + err.Error()
	}
	list := ReadV1DataTxt(txtFileContent)
	if errMsg != `` {
		return "文件数据错误: " + errMsg
	}
	if req.Merge {
		list = MergeIpRangeList(list)
	}
	if len(globalRegionMap) > 0 {
		for idx, one := range list {
			cityId := GetCityId(one.Attach, globalRegionMap)
			list[idx].CityId = cityId
		}
	}
	errMsg = VerifyIpRangeList(VerifyIpRangeListRequest{
		DataInfoList:     list,
		VerifyFullUint32: true,
		VerifyFiled7:     true,
	})
	if errMsg != `` {
		return "验证文件数据失败: " + errMsg
	}
	data := WriteV1DataBlob(list)
	err = ioutil.WriteFile(req.DbFileName, data, 0777)
	if err != nil {
		return "输出文件写入失败: " + err.Error()
	}
	return ""
}

func ReadGlobalRegionMap(regionCsv string) (globalRegionMap map[string]uint32, errMsg string) {
	regionData, err := ioutil.ReadFile(regionCsv)
	if err != nil {
		return nil, "读取region.csv失败: " + err.Error()
	}
	var recordAll [][]string
	recordAll, err = csv.NewReader(bytes.NewReader(regionData)).ReadAll()
	if err != nil {
		return nil, "读取region.csv失败2: " + err.Error()
	}
	globalRegionMap = map[string]uint32{}
	for _, line := range recordAll {
		if len(line) != 5 {
			continue
		}
		cityId, _ := strconv.Atoi(line[0])
		name := line[2]
		globalRegionMap[name] = uint32(cityId)
	}
	return globalRegionMap, ""
}

func GetCityId(region string, globalRegionMap map[string]uint32) uint32 {
	var p = strings.Split(region, "|")
	if len(p) != 5 {
		return 0
	}
	var key string
	for i := 3; i >= 0; i-- {
		if p[i] == "0" {
			continue
		}
		key = p[i]
		return globalRegionMap[key]
	}
	return 0
}

type IpRangeItem struct {
	Origin  string
	LowU32  uint32
	HighU32 uint32
	Attach  string
	CityId  uint32
}

func ReadV1DataTxt(data []byte) (list []IpRangeItem) {
	for _, one := range strings.Split(string(data), "\n") {
		one = strings.TrimSpace(one)
		if one == `` {
			continue
		}
		temp := strings.Split(one, `|`)
		if len(temp) < 2 {
			continue
		}
		sip := ipToUint32(net.ParseIP(temp[0]))
		eip := ipToUint32(net.ParseIP(temp[1]))
		list = append(list, IpRangeItem{
			Origin:  one,
			LowU32:  sip,
			HighU32: eip,
			Attach:  strings.Join(temp[2:], `|`),
		})
	}
	return list
}

func WriteV1DataTxt(list []IpRangeItem) (data []byte) {
	buf := bytes.NewBuffer(nil)
	for _, one := range list {
		buf.WriteString(uint32ToIp(one.LowU32).String() + `|` + uint32ToIp(one.HighU32).String() + `|` + one.Attach + "\n")
	}
	return buf.Bytes()
}

func WriteV1DataBlob(list []IpRangeItem) (data []byte) {
	idxMap := map[string]uint32{}
	data = make([]byte, 8)
	for _, one := range list {
		if idxMap[one.Attach] > 0 {
			continue
		}
		idxMap[one.Attach] = uint32(len(data)) | uint32((len(one.Attach)+4)<<24)
		cityIdBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(cityIdBytes, one.CityId)
		data = append(data, cityIdBytes...)
		data = append(data, one.Attach...)
	}
	binary.LittleEndian.PutUint32(data[0:], uint32(len(data)))
	sort.Slice(list, func(i, j int) bool {
		return list[i].LowU32 < list[j].LowU32
	})
	for _, one := range list {
		tmp := make([]byte, 12)
		binary.LittleEndian.PutUint32(tmp[0:], one.LowU32)
		binary.LittleEndian.PutUint32(tmp[4:], one.HighU32)
		binary.LittleEndian.PutUint32(tmp[8:], idxMap[one.Attach])
		data = append(data, tmp...)
	}
	binary.LittleEndian.PutUint32(data[4:], uint32(len(data)-12))
	return data
}

func ReadV1DataBlob(b []byte) (list []IpRangeItem, errMsg string) {
	if len(b) < 8 {
		return nil, "数据文件大小至少为8字节"
	}
	fp := getUint32(b, 0)
	lp := getUint32(b, 4)
	if fp < 8 || lp > uint32(len(b))-12 || fp > lp || (lp-fp)%12 != 0 {
		return nil, "lp, fp 指针异常 " + strconv.Itoa(int(fp)) + ", " + strconv.Itoa(int(lp))
	}

	var dataInfoList []IpRangeItem

	for idx := fp; idx <= lp; idx += 12 {
		ptr := getUint32(b, idx+8)
		attach := ``
		dataLen := (ptr >> 24) & 0xFF
		if dataLen > math.MaxUint8 {
			return nil, "附加数据长度异常1 " + strconv.Itoa(int(dataLen))
		}
		var cityId uint32
		if dataLen > 0 {
			if dataLen < 4 {
				return nil, "附加数据长度异常2 " + strconv.Itoa(int(dataLen))
			}
			ptr = ptr & 0x00FFFFFF
			if ptr+dataLen > uint32(len(b)) {
				return nil, "附加数据长度异常3 " + strconv.Itoa(int(ptr)) + "," + strconv.Itoa(int(dataLen)) + "," + strconv.Itoa(len(b))
			}
			attach = string(b[ptr+4 : ptr+dataLen])
			cityId = binary.LittleEndian.Uint32(b[ptr:])
		}
		dataInfoList = append(dataInfoList, IpRangeItem{
			LowU32:  getUint32(b, idx),
			HighU32: getUint32(b, idx+4),
			Attach:  attach,
			CityId:  cityId,
		})
	}
	return dataInfoList, ""
}

type VerifyIpRangeListRequest struct {
	DataInfoList     []IpRangeItem
	VerifyFullUint32 bool // 验证是否全部的uint32的ip都已覆盖
	VerifyFiled7     bool // 验证是否每行都有7个字段
}

func VerifyIpRangeList(req VerifyIpRangeListRequest) (errMsg string) {
	for idx := 0; idx < len(req.DataInfoList)-1; idx++ {
		left := req.DataInfoList[idx]
		right := req.DataInfoList[idx+1]

		if left.LowU32 >= right.LowU32 {
			return "ip范围未排序: " + left.Origin
		}
	}
	for _, one := range req.DataInfoList {
		if one.LowU32 > one.LowU32 {
			return "ip范围信息错误, 第一个ip必须小于等于第二个ip: " + one.Origin
		}
		if req.VerifyFiled7 && len(strings.Split(one.Attach, `|`)) != 5 {
			return "ip范围信息错误，需要有7个字段: " + one.Origin
		}
	}
	if req.VerifyFullUint32 {
		if len(req.DataInfoList) == 0 || req.DataInfoList[0].LowU32 != 0 {
			return "ip 范围缺失[0.0.0.0, ~]"
		}
		if req.DataInfoList[len(req.DataInfoList)-1].HighU32 != math.MaxUint32 {
			return "ip 范围缺失 [~, 255.255.255.255]"
		}
		for idx := 0; idx < len(req.DataInfoList)-1; idx++ {
			left := req.DataInfoList[idx]
			right := req.DataInfoList[idx+1]

			if left.HighU32+1 != right.LowU32 {
				return "ip范围缺失, [" + uint32ToIp(left.HighU32+1).String() + `, ` + uint32ToIp(right.LowU32-1).String() + "]"
			}
		}
	}
	return ""
}

func uint32ToIp(ip uint32) net.IP {
	var tmp = make([]byte, 4)
	binary.BigEndian.PutUint32(tmp, ip)
	return net.IPv4(tmp[0], tmp[1], tmp[2], tmp[3])
}

func ipToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return binary.BigEndian.Uint32([]byte{ip[0], ip[1], ip[2], ip[3]})
}

func getUint32(b []byte, offset uint32) uint32 {
	return binary.LittleEndian.Uint32(b[offset:])
}

func MergeIpRangeList(list []IpRangeItem) []IpRangeItem {
	listLen := len(list)
	merge := make([]IpRangeItem, 0, listLen)

	for idx := 0; idx < listLen; idx++ {
		mergeLen := len(merge)
		if idx > 0 && merge[mergeLen-1].Attach == list[idx].Attach && merge[mergeLen-1].HighU32+1 == list[idx].LowU32 {
			merge[mergeLen-1].HighU32 = list[idx].HighU32
			continue
		}

		merge = append(merge, list[idx])
	}

	return merge
}
