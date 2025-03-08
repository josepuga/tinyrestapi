package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/eiannone/keyboard" //Awesone utility!!!!
	"github.com/josepuga/goini"
)

const ConfigFile = "config.ini"
const JSONFile = "data.json"

// Git TAG
var version string = "<unknow>"

var server *Server

// This "hacks" if for exit to the OS with an error code. os.Exit does not
// execute any `defer`. So we need to create another "main" function
// TODO: Create error codes.

func main() { os.Exit(mainWithExit()) }

func mainWithExit() int {
	server = NewServer() // Global!
	if err := loadConfigFromFile(ConfigFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	// Preparing keyboard for read keys...
	keyEvents, err := keyboard.GetKeys(10) // 10 = Buffer size
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading keyboard. %v\n", err)
		return 1
	}
	defer keyboard.Close()

    displayInfo := `
Tiny REST API %s
By José Puga 2025. GPL3 License.
Press q or ESC to quit
`
    fmt.Printf(displayInfo, version)


	quit := make(chan struct{})
	result := 0

	// Run the HTTP Server, through Server struct.
	go func() {
		if err := server.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "error launching server: %v\n", err)
			result = 1
			close(quit)
		}
	}()

	go func() {

		// Pool for read keys pressed
		//
		for {
			event := <-keyEvents
			if event.Err != nil {
				fmt.Fprintf(os.Stderr, "error in keyboard event: %v\n", event.Err)
				result = 1
				close(quit)
				return
			}
			// 'q', ESC, CTRL-C to exit
			if event.Rune == 'q' || event.Key == keyboard.KeyEsc ||
				event.Key == keyboard.KeyCtrlC {
				close(quit)
				return
			}
		}
	}()

	<-quit
	return result
}

// loadConfigFromFile Load all configuration from the ini file and sets the
// server params.
func loadConfigFromFile(fileName string) error {
	ini := goini.NewIni()
	if err := ini.LoadFromFile(fileName); err != nil {
		return err
	}

	// Set server settings. See comments on `config.ini` about all key=value
	//
	server.Port = ini.GetInt("", "port", DefaultPort)
	if server.Port <= 0 || server.Port > 65535 {
		return fmt.Errorf("invalid listening port: %d", server.Port)
	}
    server.Host = ini.GetString("", "host", DefaultHost)

	// Adding all ini sections as "Handlers" or "virtual servers"
	//
	for _, section := range ini.GetSectionValues() {
		if section == "" { // Empty section is for server configuration, not handlers
			continue
		}
		handler := NewHandler()
		handler.ReadOnlyMode = ini.GetBool(section, "safe mode", false)

		// Path for URL request
		//
		for _, req := range ini.GetStringSlice(section, "paths", "", ",") {
			if req == "" {
				fmt.Printf("Warning [%s]: No paths.\n", section)
				continue
			}
			handler.AddPath(req)
		}

		// Allowed Methods
		//
		for _, m := range ini.GetStringSlice(section, "methods", "", ",") {
			if m == "" {
				fmt.Printf("Warning [%s]: Empty method %s.\n", section, m)
				continue
			}
			// Only set to true selected methods under the `methods=`` key
			if _, exists := handler.Methods[m]; exists {
				debugPrint("in [%s], adding method %s", section, m)
				handler.Methods[m] = true
			} else {
				fmt.Printf("Warning [%s]: Unknown method %s.\n", section, m)
			}
		}

		// Fields aliases to replace JSON fields.
		//
		for _, pairString := range ini.GetStringSlice(section, "field aliases", "", "|") {
			if pairString == "" {
				continue
			}
			// Every pair is a "key.alias" string
			pairArray := strings.Split(pairString, ",")
			if (len(pairArray)) != 2 || pairArray[0] == "" || pairArray[1] == "" {
				fmt.Printf("Warning [%s]: Malformed field alias (%s).\n", section, pairString)
				continue
			}
			handler.AddAlias(pairArray[0],pairArray[1])
		}

		// The JSON file. We keep it in memory. Every handler has it own data
		data, err := os.ReadFile(JSONFile)
		if err != nil {
			fmt.Printf("ERROR [%s]: unable to load JSON file (%s).\n", section, JSONFile)
			continue
		}
		if err = handler.SetData(data); err != nil {
			fmt.Printf("ERROR [%s]: unable to set JSON data.\n", section)
			continue
		}
		server.AddHandler(handler)
	}
	return nil
}
