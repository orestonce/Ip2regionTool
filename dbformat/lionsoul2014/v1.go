package lionsoul2014

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"errors"
	"github.com/orestonce/Ip2regionTool/dbformat"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"
)

type DbFormatLinsoul2014v1 struct {
	globalRegionMap map[string]uint32
}

func (DbFormatLinsoul2014v1) GetType() dbformat.DbFormatType {
	return dbformat.DbFormatType{
		ShowPriority: dbformat.ShowPriority_Linsoul2014v1,
		NameForCmd:   "Linsoul2014v1",
		Desc:         "Linsoul2014v1, Linsoul2014第一版db格式",
		ExtName:      "*.db",
		SupportWrite: true,
	}
}

func (DbFormatLinsoul2014v1) NeedVerifyFiled7() bool {
	return true
}

func (DbFormatLinsoul2014v1) ReadData(data []byte) (list []dbformat.IpRangeItem, err error) {
	if len(data) < 8 {
		return nil, errors.New("数据文件大小至少为8字节")
	}
	fp := getUint32(data, 0)
	lp := getUint32(data, 4)
	if fp < 8 || lp > uint32(len(data))-12 || fp > lp || (lp-fp)%12 != 0 {
		return nil, errors.New("lp, fp 指针异常 " + strconv.Itoa(int(fp)) + ", " + strconv.Itoa(int(lp)))
	}

	var dataInfoList []dbformat.IpRangeItem

	for idx := fp; idx <= lp; idx += 12 {
		ptr := getUint32(data, idx+8)
		attach := ``
		dataLen := (ptr >> 24) & 0xFF
		if dataLen > math.MaxUint8 {
			return nil, errors.New("附加数据长度异常1 " + strconv.Itoa(int(dataLen)))
		}
		//var cityId uint32
		if dataLen > 0 {
			if dataLen < 4 {
				return nil, errors.New("附加数据长度异常2 " + strconv.Itoa(int(dataLen)))
			}
			ptr = ptr & 0x00FFFFFF
			if ptr+dataLen > uint32(len(data)) {
				return nil, errors.New("附加数据长度异常3 " + strconv.Itoa(int(ptr)) + "," + strconv.Itoa(int(dataLen)) + "," + strconv.Itoa(len(data)))
			}
			attach = string(data[ptr+4 : ptr+dataLen])
			//cityId = binary.LittleEndian.Uint32(data[ptr:])
		}
		temp := strings.Split(attach, "|")
		var attachObj dbformat.IpRangeAttach
		if len(temp) >= 5 {
			attachObj.Country = temp[0]
			attachObj.Province = temp[2]
			attachObj.City = temp[3]
			attachObj.ISP = temp[4]
		}
		dataInfoList = append(dataInfoList, dbformat.IpRangeItem{
			LowU32:  getUint32(data, idx),
			HighU32: getUint32(data, idx+4),
			//Attach:    attach,
			AttachObj: attachObj,
			//CityId:  cityId,
		})
	}
	return dataInfoList, nil
}

func (d DbFormatLinsoul2014v1) FormatAttach(attach dbformat.IpRangeAttach) (value string, err error) {
	return strings.Join([]string{
		attach.Country,
		"0",
		attach.Province,
		attach.City,
		attach.ISP,
	}, "|"), nil
}

func (obj DbFormatLinsoul2014v1) WriteData(list []dbformat.IpRangeItem) (data []byte, err error) {
	//if len(obj.globalRegionMap) > 0 {
	//	for idx, one := range list {
	//		cityId := GetCityId(one.Attach, obj.globalRegionMap)
	//		list[idx].CityId = cityId
	//	}
	//}

	idxMap := map[string]uint32{}
	data = make([]byte, 8)
	for _, one := range list {
		if idxMap[one.Attach] > 0 {
			continue
		}
		idxMap[one.Attach] = uint32(len(data)) | uint32((len(one.Attach)+4)<<24)
		cityIdBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(cityIdBytes, 0) // one.CityId)
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
	return data, nil
}

func (this *DbFormatLinsoul2014v1) LoadGlobalRegionMap(regionCsv string) (err error) {
	regionData, err := ioutil.ReadFile(regionCsv)
	if err != nil {
		return errors.New("读取region.csv失败: " + err.Error())
	}
	var recordAll [][]string
	recordAll, err = csv.NewReader(bytes.NewReader(regionData)).ReadAll()
	if err != nil {
		return errors.New("读取region.csv失败2: " + err.Error())
	}
	this.globalRegionMap = map[string]uint32{}
	for _, line := range recordAll {
		if len(line) != 5 {
			continue
		}
		cityId, _ := strconv.Atoi(line[0])
		name := line[2]
		this.globalRegionMap[name] = uint32(cityId)
	}
	return nil
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

func getUint32(b []byte, offset uint32) uint32 {
	return binary.LittleEndian.Uint32(b[offset:])
}

func init() {
	dbformat.RegisterDbFormat(DbFormatLinsoul2014v1{})
}
