package sub

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/fatedier/frp/pkg/config"
)

func init() {
	rootCmd.AddCommand(reloadCmd)
}

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Hot-Reload frpc configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _, _, err := config.ParseClientConfig(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = reload(cfg)
		if err != nil {
			fmt.Printf("frpc reload error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("reload success\n")
		return nil
	},
}

func reload(clientCfg config.ClientCommonConf) error {
	if clientCfg.AdminPort == 0 {
		return fmt.Errorf("admin_port shoud be set if you want to use reload feature")
	}

	req, err := http.NewRequest("GET", "http://"+
		clientCfg.AdminAddr+":"+fmt.Sprintf("%d", clientCfg.AdminPort)+"/api/reload", nil)
	if err != nil {
		return err
	}

	authStr := "Basic " + base64.StdEncoding.EncodeToString([]byte(clientCfg.AdminUser+":"+
		clientCfg.AdminPwd))

	req.Header.Add("Authorization", authStr)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("code [%d], %s", resp.StatusCode, strings.TrimSpace(string(body)))
}
