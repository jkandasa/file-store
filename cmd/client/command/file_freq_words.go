package command

import (
	"fmt"
	"sort"

	api "github.com/jkandasa/file-store/cmd/client/api"
	printer "github.com/jkandasa/file-store/cmd/client/printer"
	"github.com/jkandasa/file-store/pkg/types"
	"github.com/spf13/cobra"
)

var (
	freqWordLimit  uint
	freqWordSortBy string
)

type TableDataFreqWords struct {
	Word  string `json:"word" yaml:"word" structs:"word"`
	Count string `json:"count" yaml:"count" structs:"count"`
}

func init() {
	rootCmd.AddCommand(freqWordsCmd)

	freqWordsCmd.PersistentFlags().UintVarP(&freqWordLimit, "limit", "n", 10, "limits the number of entries to display")
	freqWordsCmd.PersistentFlags().StringVar(&freqWordSortBy, "order", types.OrderByDsc, "order the result. options: dsc, asc")
}

var freqWordsCmd = &cobra.Command{
	Use:   "freq-words",
	Short: "Prints the frequent words count from the available text files",
	Example: `  # frequent words
  store freq-words
  
  # frequent words with custom options
  sore freq-words --limit 5 --order dsc
  `,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(serverAddress, insecure)

		if freqWordLimit == 0 {
			freqWordLimit = 10
		}

		freqWords, err := client.FreqWords(types.FreqWordsRequest{
			Limit:   freqWordLimit,
			OrderBy: freqWordSortBy,
		})
		if err != nil {
			fmt.Fprintf(ioStreams.ErrOut, "error: %s\n", err.Error())
			return
		}

		if len(freqWords) == 0 && outputFormat == printer.OutputConsole {
			fmt.Fprintln(ioStreams.Out, "No words found")
			return
		}

		// do order
		// needs to reorder again on the client side. JSON or golang map not keeping the key in order
		keys := make([]string, 0, len(freqWords))

		for key := range freqWords {
			keys = append(keys, key)
		}

		// sort the keys
		if freqWordSortBy == types.OrderByASC {
			sort.Strings(keys)
		} else {
			sort.Sort(sort.Reverse(sort.StringSlice(keys)))
		}

		sort.SliceStable(keys, func(i, j int) bool {
			if freqWordSortBy == types.OrderByASC {
				return freqWords[keys[i]] < freqWords[keys[j]]
			} else {
				return freqWords[keys[i]] > freqWords[keys[j]]
			}
		})

		headers := []string{"word", "count"}
		rows := make([]interface{}, 0)

		for _, k := range keys {
			rows = append(rows, TableDataFreqWords{Word: k, Count: fmt.Sprintf("%d", freqWords[k])})
		}

		printer.Print(ioStreams.Out, headers, rows, hideHeader, outputFormat, pretty)
	},
}
