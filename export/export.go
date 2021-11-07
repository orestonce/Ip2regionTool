package main

import (
	"github.com/orestonce/Ip2regionTool"
	"github.com/orestonce/go2cpp"
)

func main() {
	ctx := go2cpp.NewGo2cppContext(go2cpp.NewGo2cppContextReq{
		CppBaseName:       "ip2region",
		EnableQt:          false,
		QtExtendBaseClass: "",
		QtIncludeList:     nil,
	})
	ctx.Generate1(Ip2regionTool.ConvertDbToTxt)
	ctx.Generate1(Ip2regionTool.ConvertTxtToDb)
	ctx.MustCreate386LibraryInDir("Ip2regionTool-qt")
}
