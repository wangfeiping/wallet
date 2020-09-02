package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/input"
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
		Short: "Command line interface for cosmos & ethereum wallet",
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

	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "Show version info",
		RunE:  doVersion,
	}
	rootCmd.AddCommand(cmdCreate, cmdRecover, cmdVersion)

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
	ret = adapter.CosmosCreateAccount("./wallet_root/", "test", "12345678", seed.Seed)
	fmt.Println("create account: ", ret)
	showJSONString(ret)
	return nil
}

func doRecover(cmd *cobra.Command, _ []string) error {
	home := os.ExpandEnv(viper.GetString(types.FlagHome))
	viper.Set(types.FlagHome, home)
	name := viper.GetString(types.FlagName)
	if name == "" {
		return fmt.Errorf("please input the name")
	}
	// passwd := viper.GetString(types.FlagPasswd)
	// if passwd == "" {
	// 	return fmt.Errorf("please input the password")
	// }
	// mnem := viper.GetString(types.FlagMnemonic)
	// if mnem == "" {
	// 	return fmt.Errorf("please input the mnemonic")
	// }
	buf := bufio.NewReader(cmd.InOrStdin())
	passwd, err := input.GetString(types.FlagPasswd, buf)
	if err != nil {
		return err
	}
	mnem, err := input.GetString(types.FlagMnemonic, buf)
	if err != nil {
		return err
	}

	// var ko types.KeyOutput
	fmt.Println("--- cosmos")
	ret := adapter.CosmosRecoverKey(home, name, passwd, mnem)
	showJSONString(ret)
	// if err := json.Unmarshal([]byte(ret), &ko); err != nil {
	// 	fmt.Println("error: ", err)
	// 	return err
	// }
	// fmt.Println(ko.PrivKeyArmor)

	fmt.Println("--- ethereum")
	ret = adapter.EthRecoverAccount(home, name, passwd, mnem)
	showJSONString(ret)
	// if err := json.Unmarshal([]byte(ret), &ko); err != nil {
	// 	fmt.Println("error: ", err)
	// 	return err
	// }
	// fmt.Println(ko.PrivKeyArmor)
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
