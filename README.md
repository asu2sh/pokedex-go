# Pokedex CLI

A command-line interface (CLI) application written in Go that allows users to explore Pokémon data, catch Pokémon, and manage their personal Pokédex. The application interacts with the Pokémon API and includes caching for improved performance.

## Features

- **Explore Pokémon Locations**: View Pokémon map locations and explore specific areas.
- **Catch Pokémon**: Attempt to catch Pokémon and add them to your personal Pokédex.
- **Inspect Pokémon**: View detailed stats and information about Pokémon you've caught.
- **Caching**: Frequently accessed data is cached to reduce API calls and improve performance.
- **Interactive CLI**: User-friendly command-line interface with multiple commands.

---

## Setup

### Prerequisites

- Go 1.24 or later installed on your system.
- Internet connection to access the Pokémon API.

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/asu2sh/pokedex-go.git
   cd pokedex-go
   ```

2. Build the project:
   ```bash
   go build -o pokedex-cli ./cmd/main.go
   ```

3. Run the application:
   ```bash
   ./pokedex-cli
   ```

---

## Usage

Once the application is running, you can use the following commands:

### General Commands

- **`help`**: Displays a list of available commands and their descriptions.
- **`clear`**: Clears the CLI screen.
- **`exit`**: Exits the application.

### Map Exploration

- **`map`**: Fetches and displays the next 20 Pokémon map locations.
- **`mapb`**: Fetches and displays the previous 20 Pokémon map locations.
- **`explore <map_name>`**: Explores a specific map location and lists the Pokémon found there.

### Pokémon Management

- **`catch <pokemon_name>`**: Attempts to catch a Pokémon by name and adds it to your Pokédex if successful.
- **`pokedex`**: Displays all Pokémon currently in your personal Pokédex.
- **`inspect <pokemon_name>`**: Displays detailed stats and information about a Pokémon in your Pokédex.

---

## Example Commands

1. Start the application:
   ```bash
   ./pokedex-cli
   ```

2. View available commands:
   ```
   Pokedex > help
   ```

3. Explore Pokémon map locations:
   ```
   Pokedex > map
   ```

4. Catch a Pokémon:
   ```
   Pokedex > catch pikachu
   ```

5. View your Pokédex:
   ```
   Pokedex > pokedex
   ```

6. Inspect a Pokémon:
   ```
   Pokedex > inspect pikachu
   ```

---

## Development

### Project Structure

- **`cmd/main.go`**: Entry point for the CLI application.
- **`internal/poke/poke_api.go`**: Handles API interactions with the Pokémon API.
- **`internal/poke/poke_cache.go`**: Implements caching for API responses.
- **`internal/poke/poke_cache_test.go`**: Unit tests for the caching functionality.

### Running Tests

To run the tests for the caching functionality:
```bash
go test ./internal/poke
```

---

## Acknowledgments

- [PokéAPI](https://pokeapi.co/) for providing the Pokémon data.
- [Boot.dev](https://boot.dev/) for providing the foundation and resources that contributed to the development of this project.
