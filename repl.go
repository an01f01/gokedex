package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"math/rand"
    "strconv"

	"github.com/anegri01f01/pokegocli/internal/pokecache"
)

type Config struct {
	Next     	string
	Previous 	string
	PokemonInfo string
	Cache    	pokecache.Cache
	Args     	string
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config) error
}

type (
	Locations struct {
		Count    int
		Next     string
		Previous string
		Results  []Result
	}

	Result struct {
		Name string
		Url  string
	}
)

type LocationArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int           `json:"chance"`
				ConditionValues []interface{} `json:"condition_values"`
				MaxLevel        int           `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order        interface{} `json:"order"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []struct {
		Abilities []struct {
			Ability  interface{} `json:"ability"`
			IsHidden bool        `json:"is_hidden"`
			Slot     int         `json:"slot"`
		} `json:"abilities"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"past_abilities"`
	PastTypes []interface{} `json:"past_types"`
	Species   struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string      `json:"back_default"`
		BackFemale       interface{} `json:"back_female"`
		BackShiny        string      `json:"back_shiny"`
		BackShinyFemale  interface{} `json:"back_shiny_female"`
		FrontDefault     string      `json:"front_default"`
		FrontFemale      interface{} `json:"front_female"`
		FrontShiny       string      `json:"front_shiny"`
		FrontShinyFemale interface{} `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string      `json:"front_default"`
				FrontFemale  interface{} `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string      `json:"front_default"`
				FrontFemale      interface{} `json:"front_female"`
				FrontShiny       string      `json:"front_shiny"`
				FrontShinyFemale interface{} `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string      `json:"back_default"`
				BackFemale       interface{} `json:"back_female"`
				BackShiny        string      `json:"back_shiny"`
				BackShinyFemale  interface{} `json:"back_shiny_female"`
				FrontDefault     string      `json:"front_default"`
				FrontFemale      interface{} `json:"front_female"`
				FrontShiny       string      `json:"front_shiny"`
				FrontShinyFemale interface{} `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string      `json:"back_default"`
					BackFemale       interface{} `json:"back_female"`
					BackShiny        string      `json:"back_shiny"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string      `json:"back_default"`
					BackFemale       interface{} `json:"back_female"`
					BackShiny        string      `json:"back_shiny"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string      `json:"back_default"`
					BackFemale       interface{} `json:"back_female"`
					BackShiny        string      `json:"back_shiny"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationIx struct {
				ScarletViolet struct {
					FrontDefault string      `json:"front_default"`
					FrontFemale  interface{} `json:"front_female"`
				} `json:"scarlet-violet"`
			} `json:"generation-ix"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string      `json:"back_default"`
						BackFemale       interface{} `json:"back_female"`
						BackShiny        string      `json:"back_shiny"`
						BackShinyFemale  interface{} `json:"back_shiny_female"`
						FrontDefault     string      `json:"front_default"`
						FrontFemale      interface{} `json:"front_female"`
						FrontShiny       string      `json:"front_shiny"`
						FrontShinyFemale interface{} `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string      `json:"back_default"`
					BackFemale       interface{} `json:"back_female"`
					BackShiny        string      `json:"back_shiny"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string      `json:"front_default"`
					FrontFemale  interface{} `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				BrilliantDiamondShiningPearl struct {
					FrontDefault string      `json:"front_default"`
					FrontFemale  interface{} `json:"front_female"`
				} `json:"brilliant-diamond-shining-pearl"`
				Icons struct {
					FrontDefault string      `json:"front_default"`
					FrontFemale  interface{} `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

var cmdRegistry map[string]cliCommand

var pokemonRegistry map[string]Pokemon

func cleanInput(text string) []string {
	str := strings.ToLower(text)
	str = strings.ReplaceAll(str, "  ", " ")
	str = strings.TrimSpace(str)
	return strings.Split(str, " ")
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return errors.New("Could not exit the application")
}

func commandHelp(config *Config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for key, cmd := range cmdRegistry {
		fmt.Printf("%s: %s\n", key, cmd.description)
	}
	return nil
}

func commandMap(config *Config) error {

	var err error
	var locations Locations

	if val, ok := config.Cache.Get(config.Next); ok {
		if err := json.Unmarshal(val, &locations); err != nil {
			return err
		}
	} else {
		res, err := http.Get(config.Next)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(res.Body)

		if err := json.Unmarshal(body, &locations); err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode > 299 {
			return nil
		}
		config.Cache.Add(config.Next, body)
	}

	for i := 0; i < len(locations.Results); i++ {
		fmt.Println(locations.Results[i].Name)
	}

	config.Next = locations.Next
	config.Previous = locations.Previous

	return err
}

func commandMapb(config *Config) error {

	var err error
	var locations Locations

	if val, ok := config.Cache.Get(config.Previous); ok {
		if err := json.Unmarshal(val, &locations); err != nil {
			return err
		}
	} else {
		res, err := http.Get(config.Previous)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(res.Body)

		if err := json.Unmarshal(body, &locations); err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode > 299 {
			return nil
		}
		config.Cache.Add(config.Previous, body)
	}

	for i := 0; i < len(locations.Results); i++ {
		fmt.Println(locations.Results[i].Name)
	}

	config.Next = locations.Next
	config.Previous = locations.Previous

	return err
}

func commandExplore(config *Config) error {
	var err error
	var locationArea LocationArea

	var query string = config.Next + "/" + config.Args

	if val, ok := config.Cache.Get(query); ok {
		if err := json.Unmarshal(val, &locationArea); err != nil {
			return err
		}
	} else {
		res, err := http.Get(query)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(res.Body)

		if err := json.Unmarshal(body, &locationArea); err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode > 299 {
			return nil
		}
		config.Cache.Add(config.Next, body)
	}

	fmt.Println("Exploring " + config.Args + "...")
	fmt.Println("Found Pokemon:")
	for i := 0; i < len(locationArea.PokemonEncounters); i++ {
		fmt.Println(" - " + locationArea.PokemonEncounters[i].Pokemon.Name)
	}

	return err
}

func commandCatch(config *Config) error {
	var err error
	var pokemon Pokemon

	var query string = config.PokemonInfo + "/" + config.Args

	if val, ok := config.Cache.Get(query); ok {
		if err := json.Unmarshal(val, &pokemon); err != nil {
			return err
		}
	} else {
		res, err := http.Get(query)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(res.Body)

		if err := json.Unmarshal(body, &pokemon); err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode > 299 {
			return nil
		}
		config.Cache.Add(config.Next, body)
	}

	fmt.Println("Throwing a Pokeball at " + config.Args + "...")
	
	var catchTry int = rand.Intn(pokemon.BaseExperience)

	if (catchTry <= 20) {
		fmt.Println(config.Args + " was caught!")
		pokemonRegistry[config.Args] = pokemon

	} else {
		fmt.Println(config.Args + " escaped!")
	}
	return err
}

func commandInspect(config *Config) error {
	var err error

	pokemon, ok := pokemonRegistry[config.Args]
	
	if ok {
		fmt.Println("Name: " + config.Args)
		fmt.Println("Height: " + strconv.Itoa(pokemon.Height))
		fmt.Println("Weight: " + strconv.Itoa(pokemon.Weight))
		fmt.Println("Stats:")
		for i :=0; i < len(pokemon.Stats); i++ {
			fmt.Println("  -" + pokemon.Stats[i].Stat.Name + ": " + strconv.Itoa(pokemon.Stats[i].BaseStat))
		}
		fmt.Println("Types:")
		for i :=0; i < len(pokemon.Types); i++ {
			fmt.Println("  - " + pokemon.Types[i].Type.Name)
		}
	}
	return err
}

func commandPokedex(config *Config) error {
	var err error
	for name := range pokemonRegistry {
		fmt.Println("  - ", name)
	}	
	return err
}

func repl() {

	pokemonRegistry = make(map[string]Pokemon)

	cmdRegistry = map[string]cliCommand{
		"pokedex": {
			name:        "pokedex",
			description: "List all names of the pokemon the user has caught",
			callback:    commandPokedex,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects a pokemon and displays its name, weight, stats, and type(s)",
			callback:    commandInspect,
		},
		"catch": {
			name:        "catch",
			description: "Trys to catch a pokemon given the name",
			callback:    commandCatch,
		},
		"explore": {
			name:        "explore",
			description: "Displays all pokeman in a given area",
			callback:    commandExplore,
		},
		"map": {
			name:        "map",
			description: "Displays all areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays all areas",
			callback:    commandMapb,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Pokedex > ")

	conf := Config{}
	conf.PokemonInfo = "https://pokeapi.co/api/v2/pokemon/"
	conf.Next = "https://pokeapi.co/api/v2/location-area"
	conf.Cache = pokecache.NewCache(time.Minute * 5)

	for scanner.Scan() {
		line := scanner.Text()

		args := cleanInput(line)

		cmd, ok := cmdRegistry[args[0]]
		if len(args) > 1 {
			conf.Args = args[1]
		}
		if ok {
			cmd.callback(&conf)
		} else {
			fmt.Println("Command does not exists")
		}

		fmt.Printf("Pokedex > ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error during scanning: %v\n", err)
	}
}
