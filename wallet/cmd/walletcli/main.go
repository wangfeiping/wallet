package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wangfeiping/wallet/wallet/version"
)

// no-lint
const (
	FlagHome     = "home"
	FlagMnemonic = "mnemonic"
	FlagPass     = "passwd"
)

func main() {
	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:   "wallet",
		Short: "Command line interface for cosmos&eth wallet",
	}
	rootCmd.PersistentFlags().String(FlagHome, "$HOME/.coscli/", "home dir")
	viper.BindPFlag(FlagHome, rootCmd.PersistentFlags().Lookup(FlagHome))

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: "Create an account",
		RunE:  doCreate,
	}

	cmdRecover := &cobra.Command{
		Use:   "recover",
		Short: "Recover an account from mnemonic",
		RunE:  doRecover,
	}
	cmdRecover.Flags().StringP(FlagMnemonic, "m", "", "mnemonic")
	cmdRecover.Flags().StringP(FlagPass, "p", "", "passworld")

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
	home := os.ExpandEnv(viper.GetString(FlagHome))
	viper.Set(FlagHome, home)

	fmt.Println("do create...")
	// showJSONString(ret)
	return nil
}

func doRecover(_ *cobra.Command, _ []string) error {
	home := os.ExpandEnv(viper.GetString(FlagHome))
	viper.Set(FlagHome, home)
	// qosacc19ee0dmedngya6akhyyc2fllqq8hmrgkn24n62g
	// m := "wage maximum acid car catalog aisle attend rookie outdoor unusual donkey script maximum weather tiger expire negative wine evidence grass lemon forget concert planet"
	// regular trumpet envelope oak jar loop comic turkey forest frozen divide pond identify increase magnet power alarm develop depart manual dry gap coin bubble

	m := viper.GetString(FlagMnemonic)
	if m == "" {
		return fmt.Errorf("please input the mnemonic")
	}
	p := viper.GetString(FlagPass)
	if p == "" {
		return fmt.Errorf("please input the passwd")
	}

	// showJSONString(ret)
	return nil
}

func doUpdate(_ *cobra.Command, _ []string) error {
	home := os.ExpandEnv(viper.GetString(FlagHome))
	viper.Set(FlagHome, home)

	fmt.Println("do update...")
	// showJSONString(ret)
	return nil
}

func doQuery(_ *cobra.Command, _ []string) error {
	home := os.ExpandEnv(viper.GetString(FlagHome))
	viper.Set(FlagHome, home)

	fmt.Println("do query...")
	// showJSONString(ret)
	return nil
}

func doList(_ *cobra.Command, _ []string) error {
	home := os.ExpandEnv(viper.GetString(FlagHome))
	viper.Set(FlagHome, home)

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
