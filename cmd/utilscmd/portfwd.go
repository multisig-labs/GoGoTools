package utilscmd

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/spf13/cobra"
)

func newPortFwdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "portfwd [port] [url]",
		Short: "Listen to http on [port] and fwd any path to [url]",
		Long: `Useful for handling tools that expect the evm to be listening to http://localhost:8545 for example Hardhat.
		So you can say
		    ggt utils portfwd 8545 http://localhost:9650/ext/bc/C/rpc
		and now your Avalanche node is reachable at the default hardhat location
		`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			port := args[0]
			url, err := url.Parse(args[1])
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(url)

			f := func(res http.ResponseWriter, req *http.Request) {
				// proxy := httputil.NewSingleHostReverseProxy(url)

				proxy := &httputil.ReverseProxy{
					Rewrite: func(r *httputil.ProxyRequest) {
						r.SetURL(url)
						r.Out.URL.Path = url.Path
					},
				}
				proxy.ServeHTTP(res, req)
			}

			http.HandleFunc("/", f)
			log.Fatal(http.ListenAndServe(":"+port, nil))
		},
	}
	return cmd
}
