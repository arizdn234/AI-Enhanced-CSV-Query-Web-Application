# AI-Enhanced CSV Query Web Application

This web application enables users to interact with CSV data using AI models. Users can log in, view the parsed CSV data, and ask questions that are answered by AI models, integrating Google’s Gemini and Hugging Face's Tapas for data-driven responses.

## Features

- **User Authentication and Verification:** Secure login system with a verification step.
- **CSV Data Parsing and Display:** Reads CSV files, parses them into a map format, and displays column data with example rows.
- **AI Model Integration:**
  - **Hugging Face Tapas Model**: Answers natural language queries directly from CSV data.
  - **Google Gemini AI Model**: Provides conversational recommendations based on user inputs and contextual CSV data.
- **Dynamic Interaction:** Users can query the AI for insights or recommendations based on CSV data.

## Prerequisites

- Go 1.18 or higher
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Google Generative AI Go SDK](https://pkg.go.dev/github.com/google/generative-ai-go)
- [Hugging Face API](https://huggingface.co) for Tapas model
- Hugging Face and Gemini API Keys
- `.env` file to securely store API keys

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-repository.git
   cd your-repository
   ```

2. Install the dependencies:

   ```bash
   go mod tidy
   ```

3. Set up environment variables:

   Create a `.env` file in the project root and add your Hugging Face and Gemini API tokens:

   ```env
   HUGGINGFACE_TOKEN=your_huggingface_token
   GEMINI_TOKEN=your_gemini_token
   ```

4. Prepare your CSV data:

   Ensure that your CSV file (`data-series.csv`) is available in the project root. This file will be used to answer queries.

## Running the Application

1. Start the server:
   
   ```bash
   go run main.go
   ```

2. Access the application in your browser at [http://localhost:8080](http://localhost:8080).

## Option to Run the Application with Air (live-reloading)
If you want to use Air to automatically reload the server when code changes, follow these steps:

**Install Air (if not installed yet):**

```bash
go install github.com/cosmtrek/air@latest
```
Make sure the `$GOPATH/bin` directory is in your  `PATH`. Add this to your `.bashrc` or `.zshrc` if necessary:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```
Generate the Air configuration file by running this command in the root directory of your project:

```bash
air init
```
This will create a `.air.toml` file, which you can modify if needed.

Run the application with Air using the following command:

```bash
air
```

Air will watch for changes in your code and automatically restart the server when changes are detected.

You can still access the application via your browser at:

`http://localhost:8080`

## Usage

### Routes

- **`/` (Welcome Page):** The initial page of the application.
- **`/login` (Login Page):** User login with predefined credentials.
- **`/verify` (Verification Page):** Users answer a security question to proceed further.
- **`/home` (Home Page):** Displays parsed CSV data and provides a form to ask questions.
- **`/ask` (Ask AI):** Sends a user's query to Hugging Face Tapas for CSV-based answers.
- **`/recommend` (Get Recommendations):** Uses Google Gemini AI to provide conversational recommendations based on the user's input.
- **`/logout` (Logout):** Logs out the current user and clears the session.

### Authentication

The application uses predefined user credentials stored in the `allowedUsers` map. Example credentials include:

```go
var allowedUsers = map[string]string{
    "user1": "pass1",
    "user2": "pass2",
}
```

These can be modified or expanded as needed.

## Template Files

Ensure the following HTML files are placed in the `templates` directory:

- **`login.html`** - Login form for user authentication.
- **`verify.html`** - Verification question for access.
- **`home.html`** - Main page allowing user interaction with AI.

## AI Model Details

### Hugging Face Tapas Model

- **Endpoint**: `https://api-inference.huggingface.co/models/google/tapas-base-finetuned-wtq`
- **Functionality**: Provides answers to natural language queries directly from tabular CSV data.
- **Integration**: The app sends table data and a query to the Tapas model and displays the model’s response.

### Google Gemini AI Model

- **Functionality**: A conversational AI model that interacts with users to provide personalized recommendations based on prior chat history and inputs.
- **Integration**: Initializes a chat session with Gemini AI, feeding contextual CSV data and user messages to provide recommendations.

## Troubleshooting

- **API Errors**: Verify your API tokens are correctly set in the `.env` file and have sufficient permissions.
- **CSV Parsing Issues**: Ensure the CSV file format is correct, with headers and rows properly aligned.
- **Server Errors**: Check console logs for detailed error messages and ensure all dependencies are properly installed.

## Important Files

- **`main.go`**: Main application code.
- **`data-series.csv`**: Example CSV file used for querying.
- **`.env`**: Environment variables file for API tokens.
- **`templates/`**: Contains all HTML files required for the frontend.

## License

This project is licensed under the `MIT License`.

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin) for providing a lightweight web server framework.
- [Hugging Face](https://huggingface.co) for their Tapas model.
- [Google Generative AI](https://github.com/google/generative-ai-go) for the Gemini model.
