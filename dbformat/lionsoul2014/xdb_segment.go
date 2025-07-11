package lionsoul2014

import "fmt"

type Segment struct {
	StartIP uint32
	EndIP   uint32
	Region  string
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

func Long2IP(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", (ip>>24)&0xFF, (ip>>16)&0xFF, (ip>>8)&0xFF, (ip>>0)&0xFF)
}
