package lionsoul2014

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/orestonce/Ip2regionTool/dbformat"
	"strconv"
	"strings"
)

type DbFormatLinsoul2014v2 struct {
	indexPolicy IndexPolicy
}

func (DbFormatLinsoul2014v2) GetType() dbformat.DbFormatType {
	return dbformat.DbFormatType{
		ShowPriority: dbformat.ShowPriority_Linsoul2014v2,
		NameForCmd:   "Linsoul2014v2",
		Desc:         "Linsoul2014v2, Linsoul2014 第二版xdb格式",
		ExtName:      "*.xdb",
		SupportWrite: true,
	}
}

func (DbFormatLinsoul2014v2) ReadData(data []byte) (list []dbformat.IpRangeItem, err error) {
	const hLen = 256 + 512*1024
	if len(data) < hLen {
		return nil, errors.New("header length error: " + strconv.Itoa(len(data)))
	}
	startIndexPtr := binary.LittleEndian.Uint32(data[8:])
	endIndexPtr := binary.LittleEndian.Uint32(data[12:])
	if startIndexPtr < hLen || startIndexPtr >= uint32(len(data)) || endIndexPtr >= uint32(len(data)) || startIndexPtr > endIndexPtr {
		return nil, errors.New("startIndexPtr/endIndexPtr length error: " + strconv.Itoa(int(startIndexPtr)) + ", " + strconv.Itoa(int(endIndexPtr)))
	}
	for ptr := startIndexPtr; ptr <= endIndexPtr; ptr += SegmentIndexSize {
		indexBuff := data[ptr:]
		var item dbformat.IpRangeItem
		item.LowU32 = binary.LittleEndian.Uint32(indexBuff)
		item.HighU32 = binary.LittleEndian.Uint32(indexBuff[4:])
		dataLen := int(binary.LittleEndian.Uint16(indexBuff[8:]))
		dataPtr := int(binary.LittleEndian.Uint32(indexBuff[10:]))
		if dataPtr+dataLen > len(data) {
			return nil, errors.New("dataPtr/dataLen error: " + strconv.Itoa(dataPtr) + ", " + strconv.Itoa(dataLen))
		}
		item.AttachObj = decodeIpRangeAttach(string(data[dataPtr : dataPtr+dataLen]))
		list = append(list, item)
	}
	return list, nil
}

func (d DbFormatLinsoul2014v2) FormatAttach(attach dbformat.IpRangeAttach) (value string, err error) {
	return strings.Join([]string{
		attach.Country,
		"0",
		attach.Province,
		attach.City,
		attach.ISP,
	}, "|"), nil
}

func (obj DbFormatLinsoul2014v2) WriteData(list []dbformat.IpRangeItem) (data []byte, err error) {
	maker := &Maker{
		indexPolicy: obj.GetIndexPolicy(),
		vectorIndex: make([]byte, VectorIndexLength),
	}
	for _, one := range list {
		maker.segments = append(maker.segments, &Segment{
			StartIP: one.LowU32,
			EndIP:   one.HighU32,
			Region:  one.Attach,
		})
	}

	err = maker.Start()
	if err != nil {
		return nil, fmt.Errorf("failed Start: %s", err)
	}
	return maker.dstBuffer, nil
}

func (obj DbFormatLinsoul2014v2) GetIndexPolicy() IndexPolicy {
	switch obj.indexPolicy {
	case VectorIndexPolicy, BTreeIndexPolicy:
		return obj.indexPolicy
	default:
		return VectorIndexPolicy
	}
}

func (this *DbFormatLinsoul2014v2) SetIndexPolicyFromString(indexPolicyS string) (err error) {
	indexPolicy, err := IndexPolicyFromString(indexPolicyS)
	if err != nil {
		return err
	}
	this.indexPolicy = indexPolicy
	return nil
}

type TxtToXdbReq struct {
	SrcFile      string
	DstFile      string
	IndexPolicyS string
}

func init() {
	dbformat.RegisterDbFormat(DbFormatLinsoul2014v2{})
}
