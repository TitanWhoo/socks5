/*
Copyright © 2024 urCoffee <ez2t@outlook.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/TitanWhoo/socks5"
	"log"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the socks5 server",
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		host, _ := cmd.Flags().GetString("host")
		user, _ := cmd.Flags().GetString("user")
		pass, _ := cmd.Flags().GetString("pass")
		bind, _ := cmd.Flags().GetStringSlice("bind")
		tcpTimeout, _ := cmd.Flags().GetInt("tcp-timeout")
		udpTimeout, _ := cmd.Flags().GetInt("udp-timeout")
		server, err := socks5.NewServer(addr, host, user, pass, bind, tcpTimeout, udpTimeout)
		if err != nil {
			log.Fatalln(err)
		}
		err = server.ListenAndServe(nil)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.Flags().BoolVar(&socks5.Debug, "debug", false, "Enable debug mode")
	runCmd.Flags().StringP("addr", "a", "127.0.0.1:1080", "The address to listen on")
	runCmd.Flags().StringP("host", "p", "127.0.0.1", "The Host for UDP connections")
	runCmd.Flags().StringP("user", "u", "", "The username for authentication")
	runCmd.Flags().StringP("pass", "P", "", "The password for authentication")
	runCmd.Flags().StringSliceP("bind", "b", nil, "The address to bind to (CIDR list)")
	runCmd.Flags().Int("tcp-timeout", 600, "The timeout for TCP connections")
	runCmd.Flags().Int("udp-timeout", 600, "The timeout for UDP connections")
}
