package Ip2regionTool

import "fmt"

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
