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

	DbToTxtCmd := &cobra.Command{
		Use: "DbToTxt",
		Run: func(cmd *cobra.Command, args []string) {
			errMsg := Ip2regionTool.ConvertDbToTxt(dbFileName, txtFileName)
			if errMsg != `` {
				fmt.Println(errMsg)
				os.Exit(-1)
			}
		},
	}
	DbToTxtCmd.Flags().StringVarP(&txtFileName, `txt`, ``, "", ``)
	DbToTxtCmd.Flags().StringVarP(&dbFileName, `db`, ``, "", ``)
	rootCmd.AddCommand(DbToTxtCmd)

	TxtToDbCmd := &cobra.Command{
		Use: "TxtToDb",
		Run: func(cmd *cobra.Command, args []string) {
			errMsg := Ip2regionTool.ConvertTxtToDb(txtFileName, dbFileName)
			if errMsg != `` {
				fmt.Println(errMsg)
				os.Exit(-1)
			}
		},
	}
	TxtToDbCmd.Flags().StringVarP(&txtFileName, `txt`, ``, "", ``)
	TxtToDbCmd.Flags().StringVarP(&dbFileName, `db`, ``, "", ``)
	rootCmd.AddCommand(TxtToDbCmd)
}
