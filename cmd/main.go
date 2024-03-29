package main

import (
	"fmt"
	"github.com/orestonce/Ip2regionTool"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "Ip2regionTool",
}

func main() {
	rootCmd.Execute()
}

func init() {
	var dbFileName string
	var txtFileName string
	var merge bool
	var dbVersion int

	DbToTxtCmd := &cobra.Command{
		Use: "DbToTxt",
		Run: func(cmd *cobra.Command, args []string) {
			errMsg := Ip2regionTool.ConvertDbToTxt(Ip2regionTool.ConvertDbToTxt_Req{
				DbFileName:  dbFileName,
				TxtFileName: txtFileName,
				Merge:       merge,
				DbVersion:   dbVersion,
			})
			if errMsg != `` {
				fmt.Println(errMsg)
				os.Exit(-1)
			}
		},
	}
	DbToTxtCmd.Flags().StringVarP(&txtFileName, `txt`, ``, "", ``)
	DbToTxtCmd.Flags().StringVarP(&dbFileName, `db`, ``, "", ``)
	DbToTxtCmd.Flags().IntVarP(&dbVersion, "dbversion", "", 0, "0 -> auto detect, 1 -> v1, 2 -> v2")
	DbToTxtCmd.Flags().BoolVarP(&merge, "merge", "", true, "")
	rootCmd.AddCommand(DbToTxtCmd)

	var regionCsvvFileName string

	TxtToDbCmd := &cobra.Command{
		Use: "TxtToDb",
		Run: func(cmd *cobra.Command, args []string) {
			errMsg := Ip2regionTool.ConvertTxtToDb(Ip2regionTool.ConvertTxtToDb_Req{
				TxtFileName:       txtFileName,
				DbFileName:        dbFileName,
				RegionCsvFileName: regionCsvvFileName,
				Merge:             merge,
			})
			if errMsg != `` {
				fmt.Println(errMsg)
				os.Exit(-1)
			}
		},
	}
	TxtToDbCmd.Flags().StringVarP(&txtFileName, `txt`, ``, "", ``)
	TxtToDbCmd.Flags().StringVarP(&dbFileName, `db`, ``, "", ``)
	TxtToDbCmd.Flags().StringVarP(&regionCsvvFileName, "region", "", "", "")
	TxtToDbCmd.Flags().BoolVarP(&merge, "merge", "", true, "")
	rootCmd.AddCommand(TxtToDbCmd)

	var srcFile string
	var dstFile string
	var indexPolicy string

	TxtToXdbCmd := &cobra.Command{
		Use: "TxtToXdb",
		Run: func(cmd *cobra.Command, args []string) {
			errMsg := Ip2regionTool.TxtToXdb(Ip2regionTool.TxtToXdbReq{
				SrcFile:      srcFile,
				DstFile:      dstFile,
				IndexPolicyS: indexPolicy,
			})
			if errMsg != "" {
				fmt.Println(errMsg)
				os.Exit(-1)
			}
		},
	}
	TxtToXdbCmd.Flags().StringVarP(&srcFile, "srcFile", "", "", "")
	TxtToXdbCmd.Flags().StringVarP(&dstFile, "dstFile", "", "", "")
	TxtToXdbCmd.Flags().StringVarP(&indexPolicy, "indexPolicy", "", "", "[vector/btree]")
	rootCmd.AddCommand(TxtToXdbCmd)
}
