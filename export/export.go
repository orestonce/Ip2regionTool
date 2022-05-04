package main

import (
	"github.com/orestonce/Ip2regionTool"
	"github.com/orestonce/go2cpp"
)

func main() {
	ctx := go2cpp.NewGo2cppContext(go2cpp.NewGo2cppContext_Req{
		CppBaseName:                 "ip2region",
		EnableQtClass_RunOnUiThread: false,
	})
	ctx.Generate1(Ip2regionTool.ConvertDbToTxt)
	ctx.Generate1(Ip2regionTool.ConvertTxtToDb)
	ctx.MustCreateAmd64LibraryInDir("Ip2regionTool-qt")
}
