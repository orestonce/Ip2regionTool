package Ip2regionTool

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

type TxtToXdbReq struct {
	SrcFile      string
	DstFile      string
	IndexPolicyS string
}

func TxtToXdb(req TxtToXdbReq) (errMsg string) {
	indexPolicy, err := IndexPolicyFromString(req.IndexPolicyS)
	if err != nil {
		return "indexPolicy " + req.IndexPolicyS
	}
	maker, err := NewMaker(indexPolicy, req.SrcFile, req.DstFile)
	if err != nil {
		return fmt.Sprintf("failed to create %s", err)
	}
	defer maker.End()

	err = maker.Init()
	if err != nil {
		return fmt.Sprintf("failed Init: %s", err)
	}

	err = maker.Start()
	if err != nil {
		return fmt.Sprintf("failed Start: %s", err)
	}

	err = maker.End()
	if err != nil {
		return fmt.Sprintf("failed End: %s", err)
	}
	return ""
}

//
//type XdbToTxtReq struct {
//	InputXdbFile  string
//	OutputTxtFile string
//}
//
//func XdbToTxt(req XdbToTxtReq) (errMsg string) {
//	xdbData, err := ioutil.ReadFile(req.InputXdbFile)
//	if err != nil {
//		return "XdbToTxt 读取错误: " + err.Error()
//	}
//	list, errMsg := ReadV2DataBlob(xdbData)
//	if errMsg != "" {
//		return "XdbToTxt 解析错误: " + errMsg
//	}
//	txtData := WriteV1DataTxt(list)
//	err = ioutil.WriteFile(req.OutputTxtFile, txtData, 0666)
//	if err != nil {
//		return "XdbToTxt 写入错误: " + err.Error()
//	}
//	return ""
//}

func ReadV2DataBlob(b []byte) (list []IpRangeItem, errMsg string) {
	const hLen = 256 + 512*1024
	if len(b) < hLen {
		return nil, "header length error: " + strconv.Itoa(len(b))
	}
	startIndexPtr := binary.LittleEndian.Uint32(b[8:])
	endIndexPtr := binary.LittleEndian.Uint32(b[12:])
	if startIndexPtr < hLen || startIndexPtr >= uint32(len(b)) || endIndexPtr >= uint32(len(b)) || startIndexPtr > endIndexPtr {
		return nil, "startIndexPtr/endIndexPtr length error: " + strconv.Itoa(int(startIndexPtr)) + ", " + strconv.Itoa(int(endIndexPtr))
	}
	for ptr := startIndexPtr; ptr <= endIndexPtr; ptr += SegmentIndexSize {
		indexBuff := b[ptr:]
		var item IpRangeItem
		item.LowU32 = binary.LittleEndian.Uint32(indexBuff)
		item.HighU32 = binary.LittleEndian.Uint32(indexBuff[4:])
		dataLen := int(binary.LittleEndian.Uint16(indexBuff[8:]))
		dataPtr := int(binary.LittleEndian.Uint32(indexBuff[10:]))
		if dataPtr+dataLen > len(b) {
			return nil, "dataPtr/dataLen error: " + strconv.Itoa(dataPtr) + ", " + strconv.Itoa(dataLen)
		}
		item.Attach = string(b[dataPtr : dataPtr+dataLen])
		list = append(list, item)
	}
	return list, ""
}
