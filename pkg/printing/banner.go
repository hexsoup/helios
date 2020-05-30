package printing

// Banner is Styx banner function
func Banner() {
	HeliosLog("", `
          _______  _       _________ _______  _______ 
|\     /|(  ____ \( \      \__   __/(  ___  )(  ____ \
| )   ( || (    \/| (         ) (   | (   ) || (    \/
| (___) || (__    | |         | |   | |   | || (_____ 
|  ___  ||  __)   | |         | |   | |   | |(_____  )
| (   ) || (      | |         | |   | |   | |      ) |
| )   ( || (____/\| (____/\___) (___| (___) |/\____) |
|/     \|(_______/(_______/\_______/(_______)\_______)
														  	
`)
	HeliosLog("", "Port Scanning tool based on golang")
	HeliosLog("", "Made with <3 by @hexsoup")
}
