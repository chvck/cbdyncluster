package cmd

import (
	"fmt"
	"os"

	"github.com/couchbaselabs/cbdynclusterd/daemon"
	"github.com/spf13/cobra"
)

var allocateCmd = &cobra.Command{
	Use:   "allocate",
	Short: "Allocates a new cluster",
	//Long:  `Allocates a new cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		checkConfigInitialized()

		numNodes, _ := cmd.Flags().GetInt("num-nodes")
		serverVersion, _ := cmd.Flags().GetString("server-version")
		clusterID, _ := cmd.Flags().GetString("cluster-id")

		if numNodes < 0 || numNodes > 24 {
			fmt.Printf("Must allocate between 1 and 24 nodes\n")
			os.Exit(1)
		}

		reqData := daemon.CreateClusterJSON{
			ClusterID: clusterID,
		}
		for i := 0; i < numNodes; i++ {
			reqData.Nodes = append(reqData.Nodes, daemon.CreateClusterNodeJSON{
				ServerVersion: serverVersion,
			})
		}

		var respData daemon.NewClusterJSON
		err := serverRestCall("POST", "/clusters", reqData, &respData, false)
		if err != nil {
			fmt.Printf("Failed to allocate cluster: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", respData.ID)
	},
}

func init() {
	rootCmd.AddCommand(allocateCmd)

	allocateCmd.Flags().Int("num-nodes", 3, "The number of nodes to initialize")
	allocateCmd.Flags().String("server-version", "5.5.0", "The server version to use when allocating the nodes.")
	allocateCmd.Flags().String("cluster-id", "", "The, optional, id to assign to this cluster.")
}
