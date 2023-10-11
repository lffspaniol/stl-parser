/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"stl-parser/internal/stl"

	"github.com/spf13/cobra"
)

// countCmd represents the count command
var countCmd = &cobra.Command{
	Use:   "count",
	Short: "This command counts the number of triangles in a STL file.",
	Long:  `This command counts the number of triangles in a STL file.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, count, err := stl.CountTriangles(args[0])
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("%s has %d triangles \n", name, count)
	},
}

func init() {
	rootCmd.AddCommand(countCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// countCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// countCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
