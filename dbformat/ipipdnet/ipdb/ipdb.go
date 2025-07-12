package ipdb

import (
	"errors"
	"github.com/orestonce/Ip2regionTool/dbformat"
)

type DbFormatIpipdnetIpdb struct {
}

func (d DbFormatIpipdnetIpdb) GetType() dbformat.DbFormatType {
	return dbformat.DbFormatType{
		ShowPriority: dbformat.ShowPriority_IpipdnetIpdb,
		NameForCmd:   "IpipdnetIpdb",
		Desc:         "IpipdnetIpdb, ipip.net的ipdb格式",
		ExtName:      "*.ipdb",
		SupportWrite: false,
	}
}

func (d DbFormatIpipdnetIpdb) ReadData(data []byte) (list []dbformat.IpRangeItem, err error) {
	reader, err := InitBytes(data, len(data), nil)
	if err != nil {
		return nil, err
	}
	if reader.IsIPv4Support() == false {
		return nil, errors.New("no data")
	}
	nodeIdToAttachMap := map[int]dbformat.IpRangeAttach{}

	for _, one := range reader.ListV4Node() {
		item, ok := nodeIdToAttachMap[one.NodeId]
		if ok == false {
			var info map[string]string
			info, err = reader.resolveMap("CN", one.NodeId)
			if err != nil {
				return nil, err
			}
			item = dbformat.IpRangeAttach{
				Country:  info["country_name"],
				Province: info["region_name"],
				City:     info["city_name"],
				ISP:      "",
			}
			nodeIdToAttachMap[one.NodeId] = item
		}
		list = append(list, dbformat.IpRangeItem{
			LowU32:    one.LowUint32,
			HighU32:   one.HighUint32,
			AttachObj: item,
		})
	}
	return list, nil
}

func (d DbFormatIpipdnetIpdb) FormatAttach(attach dbformat.IpRangeAttach) (value string, err error) {
	return "", errors.New("implement me")
}

func (d DbFormatIpipdnetIpdb) WriteData(list []dbformat.IpRangeItem) (data []byte, err error) {
	return nil, errors.New("implement me")
}

func init() {
	dbformat.RegisterDbFormat(DbFormatIpipdnetIpdb{})
}
