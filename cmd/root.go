/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	cidr4File   string
	outFile     string
	disablePush bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cidr2route",
	Short: "Convert CIDR to openvpn route",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		readFile, err := os.Open(cidr4File)
		if err != nil {
			fmt.Println(err)
		}
		defer readFile.Close()
		fileScanner := bufio.NewScanner(readFile)
		fileScanner.Split(bufio.ScanLines)

		fh, err := os.Create(outFile)
		if err != nil {
			panic(err)
		}
		defer fh.Close()
		w := bufio.NewWriter(fh)
		defer w.Flush()

		var z string
		for fileScanner.Scan() {
			ipT := fileScanner.Text()
			ip, netmask, err := net.ParseCIDR(ipT)
			if err != nil {
				fmt.Println("Error: ", err)
				continue
			}
			m := netmask.Mask
			if disablePush {
				z = fmt.Sprintf("%s %d.%d.%d.%d\n", ip, m[0], m[1], m[2], m[3])
			} else {
				z = fmt.Sprintf(`push "route %s %d.%d.%d.%d"%s`, ip, m[0], m[1], m[2], m[3], "\n")
			}
			_, err = w.WriteString(z)
			if err != nil {
				panic(err)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cidr2route.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().StringVar(&cidr4File, "cidr4File", "cidr4.txt", "CIDR List from https://raw.githubusercontent.com/touhidurrr/iplist-youtube/main/cidr4.txt")
	rootCmd.PersistentFlags().StringVar(&outFile, "outFile", "DEFAULT", `openvpn server config file for options like this: 'push "route 216.18.168.124 255.255.255.255"'`)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cidr2route" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cidr2route")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
