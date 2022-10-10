package sub

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/fatedier/frp/pkg/config"
	"github.com/fatedier/frp/pkg/consts"
)

func init() {
	RegisterCommonFlags(xtcpCmd)

	xtcpCmd.PersistentFlags().StringVarP(&proxyName, "proxy_name", "n", "", "proxy name")
	xtcpCmd.PersistentFlags().StringVarP(&role, "role", "", "server", "role")
	xtcpCmd.PersistentFlags().StringVarP(&sk, "sk", "", "", "secret key")
	xtcpCmd.PersistentFlags().StringVarP(&serverName, "server_name", "", "", "server name")
	xtcpCmd.PersistentFlags().StringVarP(&localIP, "local_ip", "i", "127.0.0.1", "local ip")
	xtcpCmd.PersistentFlags().IntVarP(&localPort, "local_port", "l", 0, "local port")
	xtcpCmd.PersistentFlags().StringVarP(&bindAddr, "bind_addr", "", "", "bind addr")
	xtcpCmd.PersistentFlags().IntVarP(&bindPort, "bind_port", "", 0, "bind port")
	xtcpCmd.PersistentFlags().BoolVarP(&useEncryption, "ue", "", false, "use encryption")
	xtcpCmd.PersistentFlags().BoolVarP(&useCompression, "uc", "", false, "use compression")

	rootCmd.AddCommand(xtcpCmd)
}

var xtcpCmd = &cobra.Command{
	Use:   "xtcp",
	Short: "Run frpc with a single xtcp proxy",
	RunE: func(cmd *cobra.Command, args []string) error {
		clientCfg, err := parseClientCommonCfgFromCmd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		proxyConfs := make(map[string]config.ProxyConf)
		visitorConfs := make(map[string]config.VisitorConf)

		var prefix string
		if user != "" {
			prefix = user + "."
		}

		if role == "server" {
			cfg := &config.XTCPProxyConf{}
			cfg.ProxyName = prefix + proxyName
			cfg.ProxyType = consts.XTCPProxy
			cfg.UseEncryption = useEncryption
			cfg.UseCompression = useCompression
			cfg.Role = role
			cfg.Sk = sk
			cfg.LocalIP = localIP
			cfg.LocalPort = localPort
			err = cfg.CheckForCli()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			proxyConfs[cfg.ProxyName] = cfg
		} else if role == "visitor" {
			cfg := &config.XTCPVisitorConf{}
			cfg.ProxyName = prefix + proxyName
			cfg.ProxyType = consts.XTCPProxy
			cfg.UseEncryption = useEncryption
			cfg.UseCompression = useCompression
			cfg.Role = role
			cfg.Sk = sk
			cfg.ServerName = serverName
			cfg.BindAddr = bindAddr
			cfg.BindPort = bindPort
			err = cfg.Check()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			visitorConfs[cfg.ProxyName] = cfg
		} else {
			fmt.Println("invalid role")
			os.Exit(1)
		}

		err = startService(clientCfg, proxyConfs, visitorConfs, "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}
