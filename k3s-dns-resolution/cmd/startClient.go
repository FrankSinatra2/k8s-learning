/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// startClientCmd represents the startClient command
var startClientCmd = &cobra.Command{
	Use:   "startClient",
	Short: "",
	Long:  ``,
	Run:   startClient,
}

func init() {
	rootCmd.AddCommand(startClientCmd)
}

func startClient(cmd *cobra.Command, args []string) {

	api := os.Getenv("API_URL")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resBuf, err := http.Get(api)

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		body := ApiResponse{}
		err = json.NewDecoder(resBuf.Body).Decode(&body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmpl := template.New("index.html")
		tmpl, err = tmpl.ParseFiles("index.html")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, body)
	})

	http.ListenAndServe(":3001", nil)
}
