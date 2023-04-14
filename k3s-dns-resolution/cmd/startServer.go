/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// startServerCmd represents the startServer command
var startServerCmd = &cobra.Command{
	Use:   "startServer",
	Short: "",
	Long:  ``,
	Run:   startServer,
}

func init() {
	rootCmd.AddCommand(startServerCmd)
}

type ApiResponse struct {
	Message string `json:"message"`
}

func startServer(cmd *cobra.Command, args []string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")

		hostname, err := os.Hostname()

		if err != nil {
			response := ApiResponse{
				Message: "Failed to retrieve hostname.",
			}

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)

			return
		}

		response := ApiResponse{
			Message: fmt.Sprintf("Hello from %s", hostname),
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})

	http.ListenAndServe(":3000", nil)
}
