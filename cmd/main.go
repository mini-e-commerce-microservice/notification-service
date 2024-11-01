package main

import (
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{}

	rootCmd.AddCommand(consumerPushNotificationEmailCmd)

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
