package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wangfeiping/wallet/wallet/adapter"
	"github.com/wangfeiping/wallet/wallet/adapter/types"
	"github.com/wangfeiping/wallet/wallet/version"
)

func main() {
	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:   "wallet",
		Short: "Command line interface for cosmos&eth wallet",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			viper.BindPFlags(cmd.Flags())
			return nil
		},
	}
	rootCmd.PersistentFlags().String(types.FlagHome, "$HOME/.coscli/", "home dir")

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: "Create an account",
		RunE:  doCreate,
	}
	cmdCreate.Flags().Int16P(types.FlagBitSize, "b", 0, "bit size of the mnemonic")

	cmdRecover := &cobra.Command{
		Use:   "recover",
		Short: "Recover an account from mnemonic",
		RunE:  doRecover,
	}
	cmdRecover.Flags().StringP(types.FlagName, "n", "", "name")
	cmdRecover.Flags().StringP(types.FlagPasswd, "p", "", "password")
	cmdRecover.Flags().StringP(types.FlagMnemonic, "m", "", "mnemonic")

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: "Update the password",
		RunE:  doUpdate,
	}

	cmdQuery := &cobra.Command{
		Use:   "query",
		Short: "Query account info",
		RunE:  doQuery,
	}

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List account",
		RunE:  doList,
	}

	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "Show version info",
		RunE:  doVersion,
	}
	rootCmd.AddCommand(cmdCreate, cmdRecover, cmdUpdate, cmdQuery,
		cmdList, cmdVersion)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func doCreate(_ *cobra.Command, _ []string) error {
	home := os.ExpandEnv(viper.GetString(types.FlagHome))
	viper.Set(types.FlagHome, home)

	fmt.Println("do create...")
	var seed types.SeedOutput
	ret := adapter.CreateSeed()
	if err := json.Unmarshal([]byte(ret), &seed); err != nil {
		fmt.Println("error: ", err)
		return err
	}

	showJSONString(ret)
	ret = adapter.CreateAccount("./wallet_root/", "test", "12345678", seed.Seed)
	fmt.Println("create account: ", ret)
	showJSONString(ret)
	return nil
}

func doRecover(_ *cobra.Command, _ []string) error {
	home := os.ExpandEnv(viper.GetString(types.FlagHome))
	viper.Set(types.FlagHome, home)
	name := viper.GetString(types.FlagName)
	if name == "" {
		return fmt.Errorf("please input the name")
	}
	passwd := viper.GetString(types.FlagPasswd)
	if passwd == "" {
		return fmt.Errorf("please input the password")
	}
	mnem := viper.GetString(types.FlagMnemonic)
	if mnem == "" {
		return fmt.Errorf("please input the mnemonic")
	}

	ret := adapter.RecoverKey(home, name, passwd, mnem)
	showJSONString(ret)
	return nil
}

func doUpdate(_ *cobra.Command, _ []string) error {
	home := os.ExpandEnv(viper.GetString(types.FlagHome))
	viper.Set(types.FlagHome, home)

	fmt.Println("do update...")
	// showJSONString(ret)
	return nil
}

func doQuery(_ *cobra.Command, _ []string) error {
	home := os.ExpandEnv(viper.GetString(types.FlagHome))
	viper.Set(types.FlagHome, home)

	fmt.Println("do query...")
	// showJSONString(ret)
	return nil
}

func doList(_ *cobra.Command, _ []string) error {
	home := os.ExpandEnv(viper.GetString(types.FlagHome))
	viper.Set(types.FlagHome, home)

	fmt.Println("listing...")
	// showJSONString(ret)
	return nil
}

func doVersion(_ *cobra.Command, _ []string) error {
	version.ShowVersion()
	return nil
}

func showJSONString(js string) (err error) {
	var out bytes.Buffer
	if err = json.Indent(&out, []byte(js), "", "  "); err == nil {
		fmt.Println(out.String())
	}
	return
}
