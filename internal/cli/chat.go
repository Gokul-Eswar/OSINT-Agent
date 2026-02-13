package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spectre/spectre/internal/agent"
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/cobra"
)

var chatCaseID string

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start an interactive AI investigation session",
	Long: `Open an interactive terminal session with the SPECTRE Agent. 
The agent can run collectors, search entities, and provide summaries of your case using AI.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := storage.InitDB(); err != nil {
			return err
		}

		if chatCaseID == "" {
			// List recent cases or ask to create one
			cases, err := storage.ListCases()
			if err != nil || len(cases) == 0 {
				return fmt.Errorf("no cases found. Create one first with 'spectre case create'")
			}
			chatCaseID = cases[0].ID
			fmt.Printf("[*] No case ID provided. Using most recent case: %s (%s)\n", cases[0].Name, chatCaseID)
		} else {
			// Verify case exists
			c, err := storage.GetCase(chatCaseID)
			if err != nil {
				return err
			}
			if c == nil {
				return fmt.Errorf("case with ID '%s' not found", chatCaseID)
			}
			fmt.Printf("[*] Using case: %s (%s)\n", c.Name, chatCaseID)
		}

		engine := agent.NewEngine(chatCaseID)
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("--- SPECTRE AGENT SESSION ---")
		fmt.Println("Type 'exit' or 'quit' to end session.")
		fmt.Println()

		for {
			fmt.Print("You > ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "exit" || input == "quit" {
				break
			}
			if input == "" {
				continue
			}

			fmt.Print("Agent is thinking...")
			response, err := engine.Execute(input)
			fmt.Print("\r                      \r") // Clear thinking line

			if err != nil {
				fmt.Printf("\nError: %v\n", err)
				continue
			}

			fmt.Printf("Agent > %s\n\n", response)
		}

		fmt.Println("Session ended.")
		return nil
	},
}

func init() {
	chatCmd.Flags().StringVarP(&chatCaseID, "case", "c", "", "Case ID to use for the session")
	rootCmd.AddCommand(chatCmd)
}
