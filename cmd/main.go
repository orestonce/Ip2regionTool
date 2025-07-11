package main

import (
	"fmt"
	"github.com/orestonce/Ip2regionTool"
	"github.com/orestonce/Ip2regionTool/dbformat"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use: "Ip2regionTool",
}

func main() {
	initRootCmd()

	rootCmd.Execute()
}

func initRootCmd() {

	var (
		FromName         string
		FromType         string
		ToName           string
		ToType           string
		VerifyFullUint32 bool
		FillFullUint32   bool
		MergeIpRange     bool
	)

	convertCmd := &cobra.Command{
		Use: "ConvertDb",
		Run: func(cmd *cobra.Command, args []string) {
			errMsg := Ip2regionTool.ConvertDb(Ip2regionTool.ConvertDbReq{
				FromName:         FromName,
				FromType:         FromType,
				ToName:           ToName,
				ToType:           ToType,
				VerifyFullUint32: VerifyFullUint32,
				FillFullUint32:   FillFullUint32,
				MergeIpRange:     MergeIpRange,
			})
			if errMsg != `` {
				fmt.Println(errMsg)
				os.Exit(-1)
			}
		},
	}
	fromTypeList, toTypeList := GetTypeListForCmd()

	convertCmd.Flags().StringVarP(&FromType, `FromType`, ``, "", strings.Join(fromTypeList, ","))
	convertCmd.Flags().StringVarP(&FromName, `FromName`, ``, "", ``)
	convertCmd.Flags().StringVarP(&ToType, "ToType", "", "", strings.Join(toTypeList, ","))
	convertCmd.Flags().StringVarP(&ToName, "ToName", "", "", "")
	convertCmd.Flags().BoolVarP(&VerifyFullUint32, "VerifyFullUint32", "", false, "")
	convertCmd.Flags().BoolVarP(&FillFullUint32, "FillFullUint32", "", false, "")
	convertCmd.Flags().BoolVarP(&MergeIpRange, "MergeIpRange", "", false, "")

	rootCmd.AddCommand(convertCmd)
}

func GetTypeListForCmd() (fromTypeList []string, toTypeList []string) {
	list := dbformat.GetDbTypeList()

	for _, one := range list {
		fromTypeList = append(fromTypeList, one.NameForCmd)
		if one.SupportWrite {
			toTypeList = append(toTypeList, one.NameForCmd)
		}
	}
	return fromTypeList, toTypeList
}
