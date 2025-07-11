package Ip2regionTool

import (
	"fmt"
	"github.com/orestonce/Ip2regionTool/dbformat"
	_ "github.com/orestonce/Ip2regionTool/dbformat/geoip"
	_ "github.com/orestonce/Ip2regionTool/dbformat/ipipdnet"
	_ "github.com/orestonce/Ip2regionTool/dbformat/lionsoul2014"
	_ "github.com/orestonce/Ip2regionTool/dbformat/txt"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ConvertDbReq struct {
	FromName         string
	FromType         string
	ToName           string
	ToType           string
	VerifyFullUint32 bool
	FillFullUint32   bool
	MergeIpRange     bool
}

func ConvertDb(req ConvertDbReq) (errMsg string) {
	start := time.Now()

	from := dbformat.GetDbFormatByType(req.FromType)
	if from == nil {
		return "未找到文件类型 " + strconv.Quote(req.FromType)
	}

	to := dbformat.GetDbFormatByType(req.ToType)
	if to == nil {
		return "未找到文件类型 " + strconv.Quote(req.ToType)
	}

	if to.GetType().SupportWrite == false {
		return "不支持写入此类型数据 " + req.ToType
	}

	fromContent, err := ioutil.ReadFile(req.FromName)
	if err != nil {
		return "读取来源文件错误 " + err.Error()
	}

	var list []dbformat.IpRangeItem
	list, err = from.ReadData(fromContent)
	if err != nil {
		return "解析来源文件错误 " + err.Error()
	}
	fmt.Println("ReadData", len(list), time.Since(start))
	start = time.Now()

	sort.Slice(list, func(i, j int) bool {
		a, b := list[i], list[j]
		return a.LowU32 < b.LowU32
	})
	errMsg = VerifyIpRangeList(list)
	if errMsg != `` {
		return "验证文件数据失败: " + errMsg
	}

	if req.FillFullUint32 {
		list = FillFullUint32(list)
	}
	if req.VerifyFullUint32 {
		if len(list) == 0 || list[0].LowU32 != 0 {
			return "ip 范围缺失[0.0.0.0, ~]"
		}
		if list[len(list)-1].HighU32 != math.MaxUint32 {
			return "ip 范围缺失 [~, 255.255.255.255]"
		}
		for idx := 0; idx < len(list)-1; idx++ {
			left := list[idx]
			right := list[idx+1]

			if left.HighU32+1 != right.LowU32 {
				return "ip范围缺失, [" + dbformat.Uint32ToIpv4(left.HighU32+1).String() + `, ` + dbformat.Uint32ToIpv4(right.LowU32-1).String() + "]"
			}
		}
	}
	fmt.Println("sort,verify", len(list), time.Since(start))
	start = time.Now()

	for idx := range list {
		list[idx].Attach, err = to.FormatAttach(list[idx].AttachObj)
		if err != nil {
			return "格式化失败 " + err.Error()
		}
	}
	if tempObj, ok := to.(dbformat.DbNeedVerifyFiled7); ok && tempObj.NeedVerifyFiled7() {
		for _, one := range list {
			if len(strings.Split(one.Attach, `|`)) != 5 {
				return "ip范围信息错误，需要有7个字段: " + one.Origin
			}
		}
	}
	if req.MergeIpRange {
		list = MergeIpRangeList(list)
	}
	if len(list) == 0 {
		return "未读取到任何数据"
	}
	fmt.Println("format,merge", len(list), time.Since(start))
	start = time.Now()

	var toData []byte
	toData, err = to.WriteData(list)
	if err != nil {
		return "输出文件编码失败 " + err.Error()
	}
	err = ioutil.WriteFile(req.ToName, toData, 0777)
	if err != nil {
		return "输出文件写入失败: " + err.Error()
	}
	fmt.Println("Write", len(list), time.Since(start))
	start = time.Now()

	return ""
}

func FillFullUint32(list []dbformat.IpRangeItem) []dbformat.IpRangeItem {
	if len(list) == 0 {
		return nil
	}
	var other []dbformat.IpRangeItem

	first := list[0]
	if first.LowU32 != 0 {
		other = append(other, dbformat.IpRangeItem{
			LowU32:  0,
			HighU32: first.LowU32 - 1,
		})
	}
	other = append(other, first)

	// 只管item前面的是否有空缺
	for idx := 1; idx < len(list); idx++ {
		item := list[idx]
		if item.LowU32 != list[idx-1].HighU32+1 {
			other = append(other, dbformat.IpRangeItem{
				LowU32:  list[idx-1].HighU32 + 1,
				HighU32: item.LowU32 - 1,
			})
		}
		other = append(other, item)
	}

	// last已经放进other了，不用重新填充
	last := list[len(list)-1]
	if last.HighU32 != math.MaxUint32 {
		other = append(other, dbformat.IpRangeItem{
			LowU32:  last.HighU32 + 1,
			HighU32: math.MaxUint32,
		})
	}
	return other
}

type VerifyIpRangeListRequest struct {
	DataInfoList []dbformat.IpRangeItem
	VerifyFiled7 bool // 验证是否每行都有7个字段
}

func VerifyIpRangeList(list []dbformat.IpRangeItem) (errMsg string) {
	for idx := 0; idx < len(list)-1; idx++ {
		left := list[idx]
		right := list[idx+1]

		if left.LowU32 >= right.LowU32 {
			return "ip范围未排序: " + left.Origin
		}
	}
	for _, one := range list {
		if one.LowU32 > one.LowU32 {
			return "ip范围信息错误, 第一个ip必须小于等于第二个ip: " + one.Origin
		}
	}
	return ""
}

func MergeIpRangeList(list []dbformat.IpRangeItem) []dbformat.IpRangeItem {
	listLen := len(list)
	merge := make([]dbformat.IpRangeItem, 0, listLen)

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
