package Ip2regionTool

import (
	"encoding/json"
	"os"
	"testing"
)

func TestReadGlobalRegionMap(t *testing.T) {
	m, errMsg := ReadGlobalRegionMap("testdata/global_region.csv")
	if errMsg != "" {
		t.Fatal(errMsg)
		return
	}
	if len(m) != 3452 {
		t.Fatal(len(m))
		return
	}
	errMsg = ConvertTxtToDb(ConvertTxtToDb_Req{
		TxtFileName:       "testdata/ip.merge.txt",
		DbFileName:        "testdata/ip2region.conv.db",
		RegionCsvFileName: "testdata/global_region.csv",
		Merge:             false,
	})
	if errMsg != "" {
		t.Fatal(errMsg)
		return
	}
	list1 := Must_ReadV1DataBlob_File("testdata/ip2region.conv.db")	// 转换出来的db
	list2 := Must_ReadV1DataBlob_File("testdata/ip2region.db")			// ip2region 原版db
	Must_ListEqual(t, list1, list2)
}

func Must_ListEqual(t *testing.T, list1 []IpRangeItem, list2 []IpRangeItem) {
	if len(list1) != len(list2) {
		panic("len(list1) != len(list2)")
	}
	for idx :=0; idx < len(list1); idx ++ {
		a, b := list1[idx], list2[idx]
		if a.LowU32 != b.LowU32 || a.HighU32 != b.HighU32 || a.Attach != b.Attach || a.CityId != b.CityId {
			ab, _ := json.Marshal(a)
			bb, _ := json.Marshal(b)
			t.Error("Must_ListEqual ", idx, string(ab), string(bb))
		}
	}
}

func Must_ReadV1DataBlob_File(fileName string) (list []IpRangeItem) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	list, errMsg := ReadV1DataBlob(data)
	if errMsg != "" {
		panic(errMsg)
	}
	return list
}
