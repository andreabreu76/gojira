# GoJira - Project Status: In Development

## Description

GoJira is a command-line tool designed to enhance productivity by providing streamlined functionality for generating commit messages and README files. It leverages AI to assist in maintaining high-quality documentation and version control practices within software projects.

## Technologies Used

- **Go**: The primary programming language used for development.
- **OpenAI API**: Powers the AI-driven functionalities such as message generation.

## Project Structure

Below is the hierarchical representation of the project files:

```
.
├── README.md
├── functions
│   ├── commitMessage.go      # Handles commit message generation
│   └── generateReadme.go     # Facilitates README file creation
├── go.mod                    # Module definition file for Go
├── go.sum                    # Dependencies checksum file
├── gojira                    # Main binary file after build
├── install.sh                # Script for setting up the project locally
├── main.go                   # Entry point of the application
├── releases
│   └── latest                # Directory for storing the latest release
├── services
│   └── openai.go             # Integration with OpenAI API
└── utils
    ├── commons               # Common utility functions
    └── git                   # Git-related utilities
```

## Installation

To set up GoJira locally, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/gojira.git
   cd gojira
   ```

2. Run the installation script:
   ```bash
   chmod +x install.sh
   ./install.sh
   ```

3. Ensure Go is installed and the `$GOPATH` is correctly set.

4. Build the project:
   ```bash
   go build -o gojira main.go
   ```

## Usage

After installation, you can use GoJira from the command line:

- To generate a commit message:
  ```bash
  ./gojira commit <your-change-description>
  ```

- To generate a README file:
  ```bash
  ./gojira readme <project-path>
  ```

## API Documentation

The project integrates with the OpenAI API to provide its AI functionalities. You must configure your API keys in the `services/openai.go` file before using the tool.

1. Obtain an API key from OpenAI and add it to your environment variables or directly in the code file.
2. Ensure network connectivity to enable API calls.

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/new-feature`).
3. Make your changes and commit them (`git commit -m "Add new feature"`).
4. Push to the branch (`git push origin feature/new-feature`).
5. Open a Pull Request, describing what you have done.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.