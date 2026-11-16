package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/eiannone/keyboard" //Awesone utility!!!!
	"github.com/josepuga/goini"
)

// Git TAG
var version string = "<unknow>"


// Options with command line
type Options struct {
	Interactive bool
	ConfigFile  string
}

const DefaultConfigFile = "config.ini"

// Ensure chanel quit is closed only once and avoid panic. (server &/| console)
// The app can ended up with console (q, CTRL-C, ...) and/or with a server.Start() error. 
var quitOnce sync.Once 

// This "hacks" if for exit to the OS with an error code. os.Exit does not
// execute any `defer`. So we need to create another "main" function
// TODO: Create error codes.
func main() { os.Exit(mainWithExit()) }

func mainWithExit() int {
	opt, err := parseOptions()
	if err != nil {
		return 1
	}

	opt.ConfigFile = configFileName(opt)

	server := NewServer() // Global!
	if err := loadConfigFromFile(opt.ConfigFile, server); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	quit := make(chan struct{})
	result := 0

	// Run the HTTP Server, through Server struct.
	go func() {
		if err := server.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "error launching server: %v\n", err)
			result = 1
			quitOnce.Do(func() {
				close(quit)
			})
		}
	}()

	if opt.Interactive {
		if err := runInteractiveConsole(quit); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}
	<-quit
	return result
}

// configFileName Gets the file name by priority "Os Arg > Env > Default"
func configFileName(opt Options) string {
	if opt.ConfigFile != "" {
		return opt.ConfigFile
	}

	if env := os.Getenv("TINYRESTAPICONFIG"); env != "" {
		return env
	}

	return DefaultConfigFile
}

// loadConfigFromFile Load all configuration from the ini file and sets the
// server params.
func loadConfigFromFile(fileName string, server *Server) error {
	ini := goini.NewIni()
	if err := ini.LoadFromFile(fileName); err != nil {
		return err
	}

	// Set server settings. See comments on `config.ini` about all key=value
	//
	server.Json = ini.GetString("", "json", "")
	if !fileExists(server.Json) {
		return fmt.Errorf("Cannot access to JSON data: %s", server.Json)
	}

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
		handler.ReadOnlyMode = ini.GetBool(section, "readonly", false)
        
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
			handler.AddAlias(pairArray[0], pairArray[1])
		}

		// The JSON file. We keep it in memory. Every handler has it own data
		data, err := os.ReadFile(server.Json)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR [%s]: unable to load JSON file (%s).\n", section, server.Json)
			return err
		}
		if err = handler.SetData(data); err != nil {
			fmt.Printf("Warning [%s]: unable to set JSON data.\n", section)
			continue
		}
		server.AddHandler(handler)
	}
	return nil
}

// parseOptions set the options with the command line arguments, error if --help or bad parameter/s
func parseOptions() (Options, error) {
	var op Options
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	f.SetOutput(io.Discard) // Avoid print error command i.e: "flag provided but not defined: --noexist"

	f.BoolVar(&op.Interactive, "i", false, "Run in interactive mode.")
	f.BoolVar(&op.Interactive, "interactive", false, "Run in interactive mode.")
	f.StringVar(&op.ConfigFile, "c", "", "Configuration file. Overrides TINYRESTAPICONFIG and ./config.ini.")
	f.StringVar(&op.ConfigFile, "config", "", "Configuration file. Overrides TINYRESTAPICONFIG and ./config.ini.")

	if err := f.Parse(os.Args[1:]); err != nil {
		f.SetOutput(os.Stdout)
		f.Usage()
		return op, err
	}
	return op, nil
}

func runInteractiveConsole(quit chan<- struct{}) error {
	// Preparing keyboard for read keys...
	keyEvents, err := keyboard.GetKeys(10) // 10 = Buffer size
	if err != nil {
		return fmt.Errorf("error reading keyboard: %w", err)
	}

	displayInfo := `
Tiny REST API %s
By José Puga 2025. GPL3 License.
Press q or ESC to quit
`
	fmt.Printf(displayInfo, version)

	go func() {
		defer keyboard.Close()

		// Pool for read keys pressed
		//
		for {
			event := <-keyEvents
			// Muy improbable que suceda este error.
			if event.Err != nil {
				fmt.Fprintf(os.Stderr, "error in keyboard event: %v\n", event.Err)

				quitOnce.Do(func() {
					close(quit)
				})
				return
			}
			// 'q', ESC, CTRL-C to exit
			if event.Rune == 'q' || event.Key == keyboard.KeyEsc ||
				event.Key == keyboard.KeyCtrlC {
				quitOnce.Do(func() {
					close(quit)
				})
				return
			}
		}
	}()
	return nil
}
