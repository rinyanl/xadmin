package cmd

import (
	"fmt"
	"os/exec"
)

func RunstatusXary() bool {
	_, err := exec.Command("/bin/sh", "-c", "sudo pgrep -lo xray").Output()
	return err == nil
}

func RunstatusNginx() bool {
	_, err := exec.Command("/bin/sh", "-c", "sudo pgrep -lo nginx").Output()
	return err == nil
}

func RunstatusMongo() bool {
	_, err := exec.Command("/bin/sh", "-c", "sudo pgrep -lo mongo").Output()
	return err == nil
}

func RestartXary() error {
	_, err := exec.Command("/bin/sh", "-c", `sudo systemctl restart xray`).Output()

	if err != nil {
		return fmt.Errorf("xary 服务重启失败")
	}

	return nil
}

func RestartNginx() error {
	_, err := exec.Command("/bin/sh", "-c", `sudo systemctl restart nginx`).Output()

	if err != nil {
		return fmt.Errorf("nginx 服务重启失败")
	}

	return nil
}

func RestartMongo() error {
	_, err := exec.Command("/bin/sh", "-c", `sudo systemctl restart mongo`).Output()

	if err != nil {
		return fmt.Errorf("mongo 服务重启失败")
	}

	return nil
}
