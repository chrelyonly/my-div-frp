package sub

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/fatedier/frp/pkg/config"
	"github.com/fatedier/frp/pkg/consts"
)

func init() {
	RegisterCommonFlags(tcpCmd)

	tcpCmd.PersistentFlags().StringVarP(&proxyName, "proxy_name", "n", "", "proxy name")
	tcpCmd.PersistentFlags().StringVarP(&localIP, "local_ip", "i", "127.0.0.1", "local ip")
	tcpCmd.PersistentFlags().IntVarP(&localPort, "local_port", "l", 0, "local port")
	tcpCmd.PersistentFlags().IntVarP(&remotePort, "remote_port", "r", 0, "remote port")
	tcpCmd.PersistentFlags().BoolVarP(&useEncryption, "ue", "", false, "use encryption")
	tcpCmd.PersistentFlags().BoolVarP(&useCompression, "uc", "", false, "use compression")

	rootCmd.AddCommand(tcpCmd)
}

var tcpCmd = &cobra.Command{
	Use:   "tcp",
	Short: "Run frpc with a single tcp proxy",
	RunE: func(cmd *cobra.Command, args []string) error {
		clientCfg, err := parseClientCommonCfgFromCmd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cfg := &config.TCPProxyConf{}
		var prefix string
		if user != "" {
			prefix = user + "."
		}
		cfg.ProxyName = prefix + proxyName
		cfg.ProxyType = consts.TCPProxy
		cfg.LocalIP = localIP
		cfg.LocalPort = localPort
		cfg.RemotePort = remotePort
		cfg.UseEncryption = useEncryption
		cfg.UseCompression = useCompression

		err = cfg.CheckForCli()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		proxyConfs := map[string]config.ProxyConf{
			cfg.ProxyName: cfg,
		}
		err = startService(clientCfg, proxyConfs, nil, "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}
