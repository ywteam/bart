package main

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

func main() {
    var rootCmd = &cobra.Command{
        Use:   "bart",
        Short: "Basic Application Repository Template",
        Long:  `A CLI tool for creating and managing repository templates`,
    }

    var initCmd = &cobra.Command{
        Use:   "init [name]",
        Short: "Initialize a new repository",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            repoName := args[0]
            fmt.Println("Initializing repository: " + repoName)
            return
        },
    }

    initCmd.Flags().String("license", "MIT", "License type")
    initCmd.Flags().String("owner", "", "Repository owner")
    
    rootCmd.AddCommand(initCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}