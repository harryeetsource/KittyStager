package cli

import (
	"KittyStager/cmd/config"
	"KittyStager/cmd/http"
	"KittyStager/cmd/interact"
	"KittyStager/cmd/util"
	"fmt"
	i "github.com/JoaoDanielRufino/go-input-autocomplete"
	"github.com/c-bata/go-prompt"
	color "github.com/logrusorgru/aurora"
	"strconv"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "exit", Description: "Exit the program"},
		{Text: "config", Description: "Show config"},
		{Text: "target", Description: "Show targets"},
		{Text: "interact", Description: "Interact with a target"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func Cli(conf config.General) {
	for {
		t := prompt.Input("KittyStager 🐈❯ ", completer)
		switch t {
		case "exit":
			return
		case "config":
			printConfig(conf)
		case "target":
			printTarget()
		case "interact":
			printTarget()
			fmt.Printf("%s", color.Yellow("[*] Please enter the id of the target"))
			id, err := i.Read("id: ")
			if err != nil {
				util.ErrPrint(err)
				break
			}
			s, err := strconv.Atoi(id)
			if err != nil {
				util.ErrPrint(fmt.Errorf("invalid input"))
				break
			}
			kittenName, err := findId(s)
			if err != nil {
				util.ErrPrint(err)
				break
			}
			if _, ok := http.Targets[kittenName]; !ok {
				util.ErrPrint(fmt.Errorf("invalid id"))
				break
			}
			fmt.Println()
			err = interact.Interact(http.Targets[kittenName].GetTarget())
			if err != nil {
				util.ErrPrint(err)
				break
			}
		}
	}
}

func printTarget() {
	fmt.Printf("%s\n", color.Green("[*] Targets:"))
	fmt.Printf("%s\n", color.Green("Id:\tKitten name:\tIp:\t\tHostname:"))
	for name, x := range http.Targets {
		fmt.Printf("%d\t%s\t%s\t%s\n", x.GetId(), color.Yellow(name), color.Yellow(x.InitChecks.GetIp()), color.Yellow(x.InitChecks.GetHostname()))
	}
	fmt.Println()
}

func findId(id int) (string, error) {
	for _, x := range http.Targets {
		if x.GetId() == id {
			return x.GetTarget(), nil
		}
	}
	return "", fmt.Errorf("invalid id")
}

func printConfig(conf config.General) {
	fmt.Printf("\n%s\t\t%s\n", color.Green("Host:"), color.Yellow(conf.GetHost()))
	fmt.Printf("%s\t\t%d\n", color.Green("Port:"), color.Yellow(conf.GetPort()))
	fmt.Printf("%s\t%s\n", color.Green("Endpoint:"), color.Yellow(conf.GetEndpoint()))
	fmt.Printf("%s\t%s\n", color.Green("UserAgent:"), color.Yellow(conf.GetUserAgent()))
	fmt.Printf("%s\t\t%d\n", color.Green("Sleep:"), color.Yellow(conf.GetSleep()))
	for _, v := range conf.GetMalPath() {
		fmt.Printf("%s\t%s\n", color.Green("Malware path:"), color.Yellow(v))
	}
	fmt.Println()
}
