package lionsoul2014

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

// Copyright 2022 The Ip2Region Authors. All rights reserved.
// Use of this source code is governed by a Apache2.0-style
// license that can be found in the LICENSE file.

// ----
// ip2region database v2.0 structure
//
// +----------------+-------------------+---------------+--------------+
// | header space   | speed up index    |  data payload | block index  |
// +----------------+-------------------+---------------+--------------+
// | 256 bytes      | 512 KiB (fixed)   | dynamic size  | dynamic size |
// +----------------+-------------------+---------------+--------------+
//
// 1. padding space : for header info like block index ptr, version, release date eg ... or any other temporary needs.
// -- 2bytes: version number, different version means structure update, it fixed to 2 for now
// -- 2bytes: index algorithm code.
// -- 4bytes: generate unix timestamp (version)
// -- 4bytes: index block start ptr
// -- 4bytes: index block end ptr
//
//
// 2. data block : region or whatever data info.
// 3. segment index block : binary index block.
// 4. vector index block  : fixed index info for block index search speed up.
// space structure table:
// -- 0   -> | 1rt super block | 2nd super block | 3rd super block | ... | 255th super block
// -- 1   -> | 1rt super block | 2nd super block | 3rd super block | ... | 255th super block
// -- 2   -> | 1rt super block | 2nd super block | 3rd super block | ... | 255th super block
// -- ...
// -- 255 -> | 1rt super block | 2nd super block | 3rd super block | ... | 255th super block
//
//
// super block structure:
// +-----------------------+----------------------+
// | first index block ptr | last index block ptr |
// +-----------------------+----------------------+
//
// data entry structure:
// +--------------------+-----------------------+
// | 2bytes (for desc)	| dynamic length		|
// +--------------------+-----------------------+
//  data length   whatever in bytes
//
// index entry structure
// +------------+-----------+---------------+------------+
// | 4bytes		| 4bytes	| 2bytes		| 4 bytes    |
// +------------+-----------+---------------+------------+
//  start ip 	  end ip	  data length     data ptr

const VersionNo = 2
const HeaderInfoLength = 256
const VectorIndexRows = 256
const VectorIndexCols = 256
const VectorIndexSize = 8
const SegmentIndexSize = 14
const VectorIndexLength = VectorIndexRows * VectorIndexCols * VectorIndexSize

type Maker struct {
	dstHandle *os.File
	dstBuffer []byte

	indexPolicy IndexPolicy
	segments    []*Segment
	vectorIndex []byte
}

func (m *Maker) initDbHeader() {

	// make and write the header space
	var header = make([]byte, HeaderInfoLength)

	// 1, version number
	binary.LittleEndian.PutUint16(header, uint16(VersionNo))

	// 2, index policy code
	binary.LittleEndian.PutUint16(header[2:], uint16(m.indexPolicy))

	// 3, generate unix timestamp
	binary.LittleEndian.PutUint32(header[4:], uint32(time.Now().Unix()))

	// 4, index block start ptr
	binary.LittleEndian.PutUint32(header[8:], uint32(0))

	// 5, index block end ptr
	binary.LittleEndian.PutUint32(header[12:], uint32(0))

	m.dstBuffer = header
}

// refresh the vector index of the specified ip
func (m *Maker) setVectorIndex(ip uint32, ptr uint32) {
	var il0 = (ip >> 24) & 0xFF
	var il1 = (ip >> 16) & 0xFF
	var idx = il0*VectorIndexCols*VectorIndexSize + il1*VectorIndexSize
	var sPtr = binary.LittleEndian.Uint32(m.vectorIndex[idx:])
	if sPtr == 0 {
		binary.LittleEndian.PutUint32(m.vectorIndex[idx:], ptr)
		binary.LittleEndian.PutUint32(m.vectorIndex[idx+4:], ptr+SegmentIndexSize)
	} else {
		binary.LittleEndian.PutUint32(m.vectorIndex[idx+4:], ptr+SegmentIndexSize)
	}
}

// to make the binary file
func (m *Maker) Start() error {
	if len(m.segments) < 1 {
		return fmt.Errorf("empty segment list")
	}

	m.initDbHeader()

	// 1, write all the region/data to the binary file
	m.dstBuffer = append(m.dstBuffer, make([]byte, VectorIndexLength)...)

	regionPool := map[string]uint32{}

	for _, seg := range m.segments {
		_, has := regionPool[seg.Region]
		if has {
			continue
		}

		var region = []byte(seg.Region)
		if len(region) > 0xFFFF {
			return fmt.Errorf("too long region info `%s`: should be less than %d bytes", seg.Region, 0xFFFF)
		}

		// get the first ptr of the next region
		pos := len(m.dstBuffer)
		m.dstBuffer = append(m.dstBuffer, region...)

		regionPool[seg.Region] = uint32(pos)
	}

	// 2, write the index block and cache the super index block
	var indexBuff = make([]byte, SegmentIndexSize)
	var counter, startIndexPtr, endIndexPtr = 0, int64(-1), int64(-1)
	for _, seg := range m.segments {
		dataPtr, has := regionPool[seg.Region]
		if !has {
			return fmt.Errorf("missing ptr cache for region `%s`", seg.Region)
		}

		// @Note: data length should be the length of bytes.
		// this works find cuz of the string feature (byte sequence) of golang.
		var dataLen = len(seg.Region)
		if dataLen < 1 {
			// @TODO: could this even be a case ?
			return fmt.Errorf("empty region info for segment '%s'", seg)
		}

		var segList = seg.Split()
		for _, s := range segList {
			pos := len(m.dstBuffer)

			// encode the segment index
			binary.LittleEndian.PutUint32(indexBuff, s.StartIP)
			binary.LittleEndian.PutUint32(indexBuff[4:], s.EndIP)
			binary.LittleEndian.PutUint16(indexBuff[8:], uint16(dataLen))
			binary.LittleEndian.PutUint32(indexBuff[10:], dataPtr)
			m.dstBuffer = append(m.dstBuffer, indexBuff...)

			m.setVectorIndex(s.StartIP, uint32(pos))
			counter++

			// check and record the start index ptr
			if startIndexPtr == -1 {
				startIndexPtr = int64(pos)
			}

			endIndexPtr = int64(pos)
		}
	}

	// synchronized the vector index block
	copy(m.dstBuffer[HeaderInfoLength:], m.vectorIndex)

	// synchronized the segment index info
	binary.LittleEndian.PutUint32(indexBuff, uint32(startIndexPtr))
	binary.LittleEndian.PutUint32(indexBuff[4:], uint32(endIndexPtr))
	copy(m.dstBuffer[8:], indexBuff[:8])

	return nil
}
