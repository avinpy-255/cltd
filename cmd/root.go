/*
Copyright © 2025 Avinash Chowdhury <avinashx36@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/anacrolix/torrent"
	"github.com/spf13/cobra"

	cc "github.com/ivanpirog/coloredcobra"
)

var torrentLink string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cltd",
	Short: "A brief description of your application",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if torrentLink == "" {
			fmt.Println("Error: Please provide a torrent link with -L flag")
			os.Exit(1)
		}
		downloadTorrent(torrentLink)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.

func Execute() {
	cc.Init(&cc.Config{
		RootCmd:  rootCmd,
		Headings: cc.HiGreen + cc.Bold + cc.Underline,
		Commands: cc.HiYellow + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.Bold,
	})

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cltd.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&torrentLink, "link", "L", "", "URL of the torrent to download")
}

func downloadTorrent(link string) {
	clientConfig := torrent.NewDefaultClientConfig()
	client, err := torrent.NewClient(clientConfig)
	if err != nil {
		fmt.Printf("Error creating torrent client: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	// Add torent from magnet URL
	t, err := client.AddMagnet(link)
	if err != nil {
		fmt.Printf("Error adding torrent: %v\n", err)

	}

	fmt.Println("Fetching Metadata ...")
	<-t.GotInfo()

	//SEE the metadata
	info := t.Info()
	fmt.Printf("Torrent Name: %s\n", info.Name)
	fmt.Printf("Total Files: %d\n", len(info.Files))
	fmt.Printf("Total Size: %d bytes\n", info.TotalLength())

	fmt.Println("Files:")
	for _, file := range info.Files {
		fmt.Printf("- %s (%d bytes)\n", file.Path, file.Length)
	}

	t.DownloadAll()

	fmt.Println("Metadata added successfully! Downloading files....")
	for t.Info().Name != "" && t.BytesMissing() > 0 {
		fmt.Printf("Downloading: %.2f%%\n", 100*(1-float64(t.BytesMissing())/float64(t.Info().TotalLength())))
	}

	fmt.Printf(("Downlaod complete!"))

}
