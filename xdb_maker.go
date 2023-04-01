package Ip2regionTool

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
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
	srcHandle *os.File
	dstHandle *os.File

	indexPolicy IndexPolicy
	segments    []*Segment
	regionPool  map[string]uint32
	vectorIndex []byte
}

func NewMaker(policy IndexPolicy, srcFile string, dstFile string) (*Maker, error) {
	// open the source file with READONLY mode
	srcHandle, err := os.OpenFile(srcFile, os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("open source file `%s`: %w", srcFile, err)
	}

	// open the destination file with Read/Write mode
	dstHandle, err := os.OpenFile(dstFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, fmt.Errorf("open target file `%s`: %w", dstFile, err)
	}

	return &Maker{
		srcHandle: srcHandle,
		dstHandle: dstHandle,

		indexPolicy: policy,
		segments:    []*Segment{},
		regionPool:  map[string]uint32{},
		vectorIndex: make([]byte, VectorIndexLength),
	}, nil
}

func (m *Maker) initDbHeader() error {

	_, err := m.dstHandle.Seek(0, 0)
	if err != nil {
		return err
	}

	// make and write the header space
	var header = make([]byte, 256)

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

	_, err = m.dstHandle.Write(header)
	if err != nil {
		return err
	}

	return nil
}

func (m *Maker) loadSegments() error {
	var last *Segment = nil

	var scanner = bufio.NewScanner(m.srcHandle)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		var l = strings.TrimSpace(strings.TrimSuffix(scanner.Text(), "\n"))

		var ps = strings.SplitN(l, "|", 3)
		if len(ps) != 3 {
			return fmt.Errorf("invalid ip segment line `%s`", l)
		}

		sip, err := CheckIP(ps[0])
		if err != nil {
			return fmt.Errorf("check start ip `%s`: %s", ps[0], err)
		}

		eip, err := CheckIP(ps[1])
		if err != nil {
			return fmt.Errorf("check end ip `%s`: %s", ps[1], err)
		}

		if sip > eip {
			return fmt.Errorf("start ip(%s) should not be greater than end ip(%s)", ps[0], ps[1])
		}

		if len(ps[2]) < 1 {
			return fmt.Errorf("empty region info in segment line `%s`", l)
		}

		var seg = &Segment{
			StartIP: sip,
			EndIP:   eip,
			Region:  ps[2],
		}

		// check the continuity of the data segment
		if last != nil {
			if last.EndIP+1 != seg.StartIP {
				return fmt.Errorf("discontinuous data segment: last.eip+1(%d) != seg.sip(%d, %s)", sip, eip, ps[0])
			}
		}

		m.segments = append(m.segments, seg)
		last = seg
	}

	return nil
}

// Init the db binary file
func (m *Maker) Init() error {
	// init the db header
	err := m.initDbHeader()
	if err != nil {
		return fmt.Errorf("init db header: %w", err)
	}

	// load all the segments
	err = m.loadSegments()
	if err != nil {
		return fmt.Errorf("load segments: %w", err)
	}

	return nil
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

// Start to make the binary file
func (m *Maker) Start() error {
	if len(m.segments) < 1 {
		return fmt.Errorf("empty segment list")
	}

	// 1, write all the region/data to the binary file
	_, err := m.dstHandle.Seek(int64(HeaderInfoLength+VectorIndexLength), 0)
	if err != nil {
		return fmt.Errorf("seek to data first ptr: %w", err)
	}

	for _, seg := range m.segments {
		_, has := m.regionPool[seg.Region]
		if has {
			continue
		}

		var region = []byte(seg.Region)
		if len(region) > 0xFFFF {
			return fmt.Errorf("too long region info `%s`: should be less than %d bytes", seg.Region, 0xFFFF)
		}

		// get the first ptr of the next region
		pos, err := m.dstHandle.Seek(0, 1)
		if err != nil {
			return fmt.Errorf("seek to current ptr: %w", err)
		}

		_, err = m.dstHandle.Write(region)
		if err != nil {
			return fmt.Errorf("write region '%s': %w", seg.Region, err)
		}

		m.regionPool[seg.Region] = uint32(pos)
	}

	// 2, write the index block and cache the super index block
	var indexBuff = make([]byte, SegmentIndexSize)
	var counter, startIndexPtr, endIndexPtr = 0, int64(-1), int64(-1)
	for _, seg := range m.segments {
		dataPtr, has := m.regionPool[seg.Region]
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
			pos, err := m.dstHandle.Seek(0, io.SeekCurrent)
			if err != nil {
				return fmt.Errorf("seek to segment index block: %w", err)
			}

			// encode the segment index
			binary.LittleEndian.PutUint32(indexBuff, s.StartIP)
			binary.LittleEndian.PutUint32(indexBuff[4:], s.EndIP)
			binary.LittleEndian.PutUint16(indexBuff[8:], uint16(dataLen))
			binary.LittleEndian.PutUint32(indexBuff[10:], dataPtr)
			_, err = m.dstHandle.Write(indexBuff)
			if err != nil {
				return fmt.Errorf("write segment index for '%s': %w", s.String(), err)
			}

			m.setVectorIndex(s.StartIP, uint32(pos))
			counter++

			// check and record the start index ptr
			if startIndexPtr == -1 {
				startIndexPtr = pos
			}

			endIndexPtr = pos
		}
	}

	// synchronized the vector index block
	_, err = m.dstHandle.Seek(int64(HeaderInfoLength), 0)
	if err != nil {
		return fmt.Errorf("seek vector index first ptr: %w", err)
	}
	_, err = m.dstHandle.Write(m.vectorIndex)
	if err != nil {
		return fmt.Errorf("write vector index: %w", err)
	}

	// synchronized the segment index info
	binary.LittleEndian.PutUint32(indexBuff, uint32(startIndexPtr))
	binary.LittleEndian.PutUint32(indexBuff[4:], uint32(endIndexPtr))
	_, err = m.dstHandle.Seek(8, io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek segment index ptr: %w", err)
	}

	_, err = m.dstHandle.Write(indexBuff[:8])
	if err != nil {
		return fmt.Errorf("write segment index ptr: %w", err)
	}

	return nil
}

func (m *Maker) End() error {
	err := m.dstHandle.Close()
	if err != nil {
		return err
	}

	err = m.srcHandle.Close()
	if err != nil {
		return err
	}

	return nil
}
