package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type AIModelConnector struct {
	Client *http.Client
}

type Inputs struct {
	Table map[string][]string `json:"table"`
	Query string              `json:"query"`
}

type GPT2Inputs struct {
	Inputs string `json:"inputs"`
}

type Response struct {
	Answer      string   `json:"answer"`
	Coordinates [][]int  `json:"coordinates"`
	Cells       []string `json:"cells"`
	Aggregator  string   `json:"aggregator"`
}

type GeminiAIModelConnector struct {
	Client *genai.Client
}

type GeminiResponse struct {
	Answer string `json:"answer"`
}

type User struct {
	Username string
	Password string
}

var allowedUsers = map[string]string{
	"user1": "password1",
	"user2": "password2",
	// tambahkan pengguna lain sesuai kebutuhan
}

type VerificationQuestion struct {
	Question string
	Answer   int
}

func CsvToSlice(data string) (map[string][]string, error) {
	reader := csv.NewReader(strings.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]string)
	if len(records) == 0 {
		return result, nil
	}

	headers := records[0]

	for _, header := range headers {
		result[header] = []string{}
	}

	for _, row := range records[1:] {
		for j, cell := range row {
			result[headers[j]] = append(result[headers[j]], cell)
		}
	}

	return result, nil
}

func (c *AIModelConnector) ConnectAIModel(payload Inputs, token string) (Response, error) {
	url := "https://api-inference.huggingface.co/models/google/tapas-base-finetuned-wtq"

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return Response{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return Response{}, fmt.Errorf("unexpected status code: %d, response body: %s", resp.StatusCode, bodyString)
	}

	var result struct {
		Answer      string      `json:"answer"`
		Coordinates [][]float64 `json:"coordinates"`
		Cells       []string    `json:"cells"`
		Aggregator  string      `json:"aggregator"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Response{}, err
	}

	coordinates := make([][]int, len(result.Coordinates))
	for i, coord := range result.Coordinates {
		intCoords := make([]int, len(coord))
		for j, c := range coord {
			intCoords[j] = int(c)
		}
		coordinates[i] = intCoords
	}

	// fmt.Printf("result.Answer: %v\n", result.Answer)
	// fmt.Printf("result.Aggregator: %v\n", result.Aggregator)

	return Response{
		Answer:      result.Answer,
		Coordinates: coordinates,
		Cells:       result.Cells,
		Aggregator:  result.Aggregator,
	}, nil
}

func DisplayCsvInfo(data string) (string, error) {
	reader := csv.NewReader(strings.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	if len(records) == 0 {
		return "CSV data is empty", nil
	}

	columnHeaders := records[0]
	maxRows := 2
	var info strings.Builder
	info.WriteString("CSV Columns:\n")
	for _, header := range columnHeaders {
		info.WriteString(header)
		info.WriteString(", ")
	}
	info.WriteString("\nExample rows:\n")
	for i := 1; i <= maxRows && i < len(records); i++ {
		for _, cell := range records[i] {
			info.WriteString(cell)
			info.WriteString(", ")
		}
		info.WriteString("\n")
	}

	return info.String(), nil
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	token := os.Getenv("HUGGINGFACE_TOKEN")

	csvData, err := os.ReadFile("data-series.csv")
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	table, err := CsvToSlice(string(csvData))
	if err != nil {
		fmt.Println("Error parsing CSV:", err)
		return
	}

	rand.Seed(time.Now().UnixNano())

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	// r.Static("/static", "./static")

	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24, // 24 hours
		HttpOnly: true,
	})
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "welcome.html", nil)
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		storedPassword, ok := allowedUsers[username]
		if !ok || storedPassword != password {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"error": "Username or password incorrect, please try again.",
			})
			return
		}

		session := sessions.Default(c)
		session.Set("user", username)
		session.Save()

		c.Redirect(http.StatusFound, "/verify")
	})

	r.GET("/verify", func(c *gin.Context) {
		question := VerificationQuestion{
			Question: "What is 2 + 6?",
			Answer:   8,
		}

		session := sessions.Default(c)
		session.Set("correctAnswer", question.Answer)
		session.Save()

		c.HTML(http.StatusOK, "verify.html", gin.H{
			"question": question.Question,
		})
	})

	r.POST("/verify", func(c *gin.Context) {
		answer := c.PostForm("answer")
		correctAnswer := "8" // Sesuaikan dengan jawaban yang Anda atur
		if answer != correctAnswer {
			c.HTML(http.StatusUnauthorized, "verify.html", gin.H{
				"error": "Jawaban salah, silakan coba lagi.",
			})
			return
		}
		c.Redirect(http.StatusFound, "/home")
	})

	r.GET("/home", AuthRequired(), func(c *gin.Context) {
		info, err := DisplayCsvInfo(string(csvData))
		if err != nil {
			c.String(http.StatusInternalServerError, "Error displaying CSV info: %v", err)
			return
		}
		c.HTML(http.StatusOK, "home.html", gin.H{
			"info": info,
		})
	})

	r.POST("/ask", AuthRequired(), func(c *gin.Context) {
		var input struct {
			Question string `form:"question"`
		}

		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot parse form"})
			return
		}

		tapasConnector := &AIModelConnector{
			Client: &http.Client{},
		}

		response, err := tapasConnector.ConnectAIModel(Inputs{Table: table, Query: input.Question}, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "AI model error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"answer":      response.Answer,
			"coordinates": response.Coordinates,
			"cells":       response.Cells,
			"aggregator":  response.Aggregator,
		})
	})

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_TOKEN")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	cs := model.StartChat()
	cs.History = []*genai.Content{
		{
			Parts: []genai.Part{
				genai.Text("Hello, I have 2 dogs in my house."),
			},
			Role: "user",
		},
		{
			Parts: []genai.Part{
				genai.Text("Great to meet you. What would you like to know?"),
			},
			Role: "model",
		},
	}

	rec := csv.NewReader(strings.NewReader(string(csvData)))
	for {
		record, err := rec.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		cs.History = append(cs.History, &genai.Content{
			Parts: []genai.Part{
				genai.Text(strings.Join(record, ", ")),
			},
			Role: "user",
		})
	}

	cs.History = append(cs.History, &genai.Content{
		Parts: []genai.Part{
			genai.Text("ayo kita bicara bahasa indonesia!"),
		},
		Role: "user",
	})

	r.POST("/recommend", func(c *gin.Context) {
		var input struct {
			Text string `json:"text"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot parse request"})
			return
		}

		resp, err := cs.SendMessage(ctx, genai.Text(input.Text))
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Printf("input: %v\n", input)

		for _, ct := range resp.Candidates {
			if ct.Content != nil {
				// fmt.Println(*ct.Content)
				c.JSON(http.StatusOK, gin.H{
					"input":          input.Text,
					"recommendation": *ct.Content,
				})
			}
		}
	})

	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusFound, "/")
	})

	port := "8080"

	log.Printf("Server is running on port %v\n\n`http://localhost:%v`", port, port)
	if err := r.Run(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
