package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/robot/client"
	"go.viam.com/rdk/services/datamanager"
	"go.viam.com/utils"
	"go.viam.com/utils/rpc"
)

type Config struct {
	Address string `json:"address"`
	Id      string `json:"id"`
	Secret  string `json:"secret"`
}

func main() { utils.ContextualMain(run, logging.NewDebugLogger("data-manager-client")) }

func parseArgs(args []string) (string, *string, error) {
	var dur *string
	var configFilePath string
	if len(args) == 2 {
		configFilePath = args[1]
	} else if len(args) == 3 {
		configFilePath = args[1]
		dur = &args[2]
	} else {
		return "", nil, fmt.Errorf("usage %s config_file_path <loop_interval>", args[0])
	}
	return configFilePath, dur, nil
}

func run(ctx context.Context, args []string, logger logging.Logger) error {
	configFilePath, durStr, err := parseArgs(args)
	if err != nil {
		return err
	}

	config, err := parseConfig(configFilePath)
	if err != nil {
		return err
	}

	machine, err := client.New(
		context.Background(),
		config.Address,
		logger,
		client.WithDialOptions(rpc.WithEntityCredentials(
			config.Id,
			rpc.Credentials{
				Type:    rpc.CredentialsTypeAPIKey,
				Payload: config.Secret,
			})),
	)
	if err != nil {
		return err
	}

	defer machine.Close(context.Background())

	dm, err := datamanager.FromRobot(machine, "data_manager-1")
	if err != nil {
		return err
	}

	if durStr != nil {
		dur, err := time.ParseDuration(*durStr)
		if err != nil {
			return err
		}
		logger.Infof("calling Sync every %s", dur)
		for utils.SelectContextOrWait(ctx, dur) {
			if err := dm.Sync(ctx, map[string]interface{}{}); err != nil {
				return err
			}
		}
		return nil
	}

	logger.Info("calling Sync once")
	return dm.Sync(ctx, map[string]interface{}{})

}

func parseConfig(configFilePath string) (Config, error) {
	f, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("error opening config file %s, err: %v", configFilePath, err)
	}
	defer f.Close()

	jd := json.NewDecoder(f)
	config := &Config{}
	if err := jd.Decode(config); err != nil {
		return Config{}, fmt.Errorf("error parsing config file %s, err: %v", configFilePath, err)
	}
	return *config, nil
}
