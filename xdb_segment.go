package Ip2regionTool

import (
	"fmt"
	"strings"
)

type Segment struct {
	StartIP uint32
	EndIP   uint32
	Region  string
}

func SegmentFrom(seg string) (*Segment, error) {
	var ps = strings.SplitN(seg, "|", 3)
	if len(ps) != 3 {
		return nil, fmt.Errorf("invalid ip segment `%s`", seg)
	}

	sip, err := CheckIP(ps[0])
	if err != nil {
		return nil, fmt.Errorf("check start ip `%s`: %s", ps[0], err)
	}

	eip, err := CheckIP(ps[1])
	if err != nil {
		return nil, fmt.Errorf("check end ip `%s`: %s", ps[1], err)
	}

	if sip > eip {
		return nil, fmt.Errorf("start ip(%s) should not be greater than end ip(%s)", ps[0], ps[1])
	}

	return &Segment{
		StartIP: sip,
		EndIP:   eip,
		Region:  ps[2],
	}, nil
}

// Split the segment based on the pre-two bytes
func (s *Segment) Split() []*Segment {
	// 1, split the segment with the first byte
	var tList []*Segment
	var sByte1, eByte1 = (s.StartIP >> 24) & 0xFF, (s.EndIP >> 24) & 0xFF
	var nSip = s.StartIP
	for i := sByte1; i <= eByte1; i++ {
		sip := (i << 24) | (nSip & 0xFFFFFF)
		eip := (i << 24) | 0xFFFFFF
		if eip < s.EndIP {
			nSip = (i + 1) << 24
		} else {
			eip = s.EndIP
		}

		// append the new segment (maybe)
		tList = append(tList, &Segment{
			StartIP: sip,
			EndIP:   eip,
			// @Note: don't bother to copy the region
			/// Region: s.Region,
		})
	}

	// 2, split the segments with the second byte
	var segList []*Segment
	for _, seg := range tList {
		base := seg.StartIP & 0xFF000000
		nSip := seg.StartIP
		sb2, eb2 := (seg.StartIP>>16)&0xFF, (seg.EndIP>>16)&0xFF
		for i := sb2; i <= eb2; i++ {
			sip := base | (i << 16) | (nSip & 0xFFFF)
			eip := base | (i << 16) | 0xFFFF
			if eip < seg.EndIP {
				nSip = 0
			} else {
				eip = seg.EndIP
			}

			segList = append(segList, &Segment{
				StartIP: sip,
				EndIP:   eip,
				Region:  s.Region,
			})
		}
	}

	return segList
}

func (s *Segment) String() string {
	return Long2IP(s.StartIP) + "|" + Long2IP(s.EndIP) + "|" + s.Region
}
