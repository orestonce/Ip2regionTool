package main

import (
	"fmt"
	"github.com/orestonce/Ip2regionTool"
	"github.com/orestonce/Ip2regionTool/dbformat"
	"github.com/orestonce/go2cpp"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	ctx := go2cpp.NewGo2cppContext(go2cpp.NewGo2cppContext_Req{
		CppBaseName:                 "ip2region",
		EnableQtClass_RunOnUiThread: true,
		EnableQtClass_Toast:         true,
	})
	ctx.Generate1(Ip2regionTool.ConvertDb)
	ctx.Generate1(dbformat.GetDbTypeList)

	if os.Getenv("GITHUB_ACTIONS") == "" { // 正常编译
		ctx.MustCreateAmd64LibraryInDir("Ip2regionTool-qt")
	} else { // github actions 编译
		version := strings.TrimPrefix(os.Getenv("GITHUB_REF_NAME"), "v")
		urlStr := "https://github.com/" + os.Getenv("GITHUB_REPOSITORY")
		ctx.MustCreateAmd64LibraryInDir("Ip2regionTool-qt")
		WriteVersionDotRc(WriteVersionDotRc_Req{
			Version:      version,
			ProductName:  "Ip2region数据转换工具",
			CopyRightUrl: urlStr,
			OutputRcName: "Ip2regionTool-qt/version.rc",
		})
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		type buildCfg struct {
			GOOS   string
			GOARCH string
			Ext    string
		}
		var list = []buildCfg{
			{
				GOOS:   "linux",
				GOARCH: "386",
			},
			{
				GOOS:   "darwin",
				GOARCH: "amd64",
			},
			{
				GOOS:   "windows",
				GOARCH: "386",
				Ext:    ".exe",
			},
		}
		for _, cfg := range list {
			name := "Ip2regionTool_cli_" + cfg.GOOS + "_" + cfg.GOARCH + cfg.Ext
			cmd := exec.Command("go", "build", "-ldflags", "-s -w", "-o", filepath.Join(wd, "bin", name)) // "-trimpath"
			cmd.Dir = filepath.Join(wd, "cmd")
			cmd.Env = append(os.Environ(), "GOOS="+cfg.GOOS)
			cmd.Env = append(cmd.Env, "GOARCH="+cfg.GOARCH)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				fmt.Println(cmd.Dir)
				panic(err)
			}
			fmt.Println("done", name)
		}
	}

}

type WriteVersionDotRc_Req struct {
	Version      string
	ProductName  string
	CopyRightUrl string
	OutputRcName string
}

func WriteVersionDotRc(req WriteVersionDotRc_Req) {
	tmp := strings.Split(req.Version, ".")
	ok := len(tmp) == 3
	for _, v := range tmp {
		vi, err := strconv.Atoi(v)
		if err != nil {
			ok = false
			break
		}
		if vi < 0 {
			ok = false
			break
		}
	}
	if ok == false {
		panic("version invalid: " + strconv.Quote(req.Version))
	}
	tmp = append(tmp, "0")
	v1 := strings.Join(tmp, ",")
	data := []byte(`IDI_ICON1 ICON "favicon.ico"

#if defined(UNDER_CE)
#include <winbase.h>
#else
#include <winver.h>
#endif

VS_VERSION_INFO VERSIONINFO
    FILEVERSION ` + v1 + `
    PRODUCTVERSION ` + v1 + `
    FILEFLAGSMASK 0x3fL
#ifdef _DEBUG
    FILEFLAGS VS_FF_DEBUG
#else
    FILEFLAGS 0x0L
#endif
    FILEOS VOS__WINDOWS32
    FILETYPE VFT_DLL
    FILESUBTYPE 0x0L
    BEGIN
        BLOCK "StringFileInfo"
        BEGIN
            BLOCK "080404b0"
            BEGIN
                VALUE "ProductVersion", "` + req.Version + `.0\0"
                VALUE "ProductName", "` + req.ProductName + `\0"
                VALUE "LegalCopyright", "https://github.com/orestonce/m3u8d\0"
                VALUE "FileDescription", "` + req.ProductName + `\0"
           END
        END

        BLOCK "VarFileInfo"
        BEGIN
            VALUE "Translation", 0x804, 1200
        END
    END
`)
	data, err := simplifiedchinese.GBK.NewEncoder().Bytes(data)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(req.OutputRcName, data, 0777)
	if err != nil {
		panic(err)
	}
}
