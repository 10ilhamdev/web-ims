package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"crypto/tls"
	"strings"
	"time"

	goravelhttp "github.com/goravel/framework/contracts/http"

	"ims/app/facades"
	"ims/app/models"
)

type ChatController struct{}

func NewChatController() *ChatController {
	return &ChatController{}
}

func (r *ChatController) InitSupport(ctx goravelhttp.Context) goravelhttp.Response {
	name := ctx.Request().Input("name")
	email := ctx.Request().Input("email")

	if name == "" || email == "" {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
			"error": "Name and email are required",
		})
	}

	token := generateRandomToken()
	chat := models.Chat{
		Token: token,
		Name:  name,
		Email: email,
	}

	// Associate user if logged in
	user := GetCurrentUser(ctx)
	if user != nil {
		chat.UserID = &user.ID
	}

	err := facades.Orm().Query().Create(&chat)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error": "Failed to create chat session: " + err.Error(),
		})
	}

	// Create initial system message
	welcomeMsg := models.ChatMessage{
		ChatID:     chat.ID,
		SenderType: "admin",
		Message:    "Hello! Welcome to IMS Support. An agent will be with you shortly. Please feel free to type your question below.",
	}
	_ = facades.Orm().Query().Create(&welcomeMsg)

	return ctx.Response().Json(http.StatusOK, map[string]any{
		"token":   token,
		"chat_id": chat.ID,
	})
}

func (r *ChatController) SendSupportMessage(ctx goravelhttp.Context) goravelhttp.Response {
	token := ctx.Request().Input("token")
	message := ctx.Request().Input("message")

	if token == "" || message == "" {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
			"error": "Token and message are required",
		})
	}

	var chat models.Chat
	err := facades.Orm().Query().Where("token = ?", token).First(&chat)
	if err != nil || chat.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"error": "Support chat session not found",
		})
	}

	msg := models.ChatMessage{
		ChatID:     chat.ID,
		SenderType: "user",
		Message:    message,
	}

	err = facades.Orm().Query().Create(&msg)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error": "Failed to save message",
		})
	}

	// Update chat updated_at timestamp
	_ = facades.Orm().Query().Save(&chat)

	return ctx.Response().Json(http.StatusOK, map[string]any{
		"success": true,
	})
}

func (r *ChatController) GetSupportMessages(ctx goravelhttp.Context) goravelhttp.Response {
	token := ctx.Request().Input("token")
	if token == "" {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
			"error": "Token is required",
		})
	}

	var chat models.Chat
	err := facades.Orm().Query().Where("token = ?", token).First(&chat)
	if err != nil || chat.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"error": "Support chat session not found",
		})
	}

	var messages []models.ChatMessage
	err = facades.Orm().Query().Where("chat_id = ?", chat.ID).Order("id asc").Get(&messages)
	if err != nil {
		messages = []models.ChatMessage{}
	}

	return ctx.Response().Json(http.StatusOK, messages)
}

func (r *ChatController) AskAI(ctx goravelhttp.Context) goravelhttp.Response {
	message := ctx.Request().Input("message")
	if message == "" {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
			"error": "Message is required",
		})
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		if val := facades.Config().Env("GEMINI_API_KEY", ""); val != nil {
			apiKey = fmt.Sprintf("%v", val)
		}
	}
	apiKey = strings.Trim(apiKey, "\"' ")

	if apiKey == "" {
		return ctx.Response().Json(http.StatusOK, map[string]any{
			"reply": "AI Assistant is offline. Please configure GEMINI_API_KEY in the .env file to enable dynamic AI responses, or click 'Chat with Admin' to contact us directly!",
		})
	}

	lang := ctx.Request().Input("lang", "en")

	// Call Gemini API (completely free tier)
	reply, err := callGeminiAPI(apiKey, message, lang)
	if err != nil {
		fmt.Printf("[Gemini API Error] %v\n", err)
		return ctx.Response().Json(http.StatusOK, map[string]any{
			"reply": fmt.Sprintf("Google Gemini API error: %v. Please verify your GEMINI_API_KEY in the .env file.", err),
		})
	}

	return ctx.Response().Json(http.StatusOK, map[string]any{
		"reply": reply,
	})
}

// Admin handlers
func (r *ChatController) AdminChatView(ctx goravelhttp.Context) goravelhttp.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	return ctx.Response().View().Make("admin/chat.tmpl", map[string]any{
		"User": user,
	})
}

func (r *ChatController) AdminGetSessions(ctx goravelhttp.Context) goravelhttp.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	var chats []models.Chat
	err := facades.Orm().Query().Order("updated_at desc").Get(&chats)
	if err != nil {
		chats = []models.Chat{}
	}

	// Fetch last message for each chat
	type ChatSessionResponse struct {
		models.Chat
		LastMessage string `json:"last_message"`
		LastTime    string `json:"last_time"`
	}

	var result []ChatSessionResponse
	for _, c := range chats {
		var lastMsg models.ChatMessage
		_ = facades.Orm().Query().Where("chat_id = ?", c.ID).Order("id desc").First(&lastMsg)
		
		lastTimeStr := ""
		if lastMsg.CreatedAt != nil {
			lastTimeStr = lastMsg.CreatedAt.ToDateTimeString()
		}
		
		result = append(result, ChatSessionResponse{
			Chat:        c,
			LastMessage: lastMsg.Message,
			LastTime:    lastTimeStr,
		})
	}

	return ctx.Response().Json(http.StatusOK, result)
}

func (r *ChatController) AdminSendReply(ctx goravelhttp.Context) goravelhttp.Response {
	user := GetCurrentUser(ctx)
	if user == nil || user.Role != "admin" {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	chatIDStr := ctx.Request().Input("chat_id")
	message := ctx.Request().Input("message")

	if chatIDStr == "" || message == "" {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
			"error": "Chat ID and message are required",
		})
	}

	var chat models.Chat
	err := facades.Orm().Query().Where("id = ?", chatIDStr).First(&chat)
	if err != nil || chat.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"error": "Chat session not found",
		})
	}

	msg := models.ChatMessage{
		ChatID:     chat.ID,
		SenderType: "admin",
		Message:    message,
	}

	err = facades.Orm().Query().Create(&msg)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error": "Failed to save message",
		})
	}

	// Update chat updated_at timestamp
	_ = facades.Orm().Query().Save(&chat)

	return ctx.Response().Json(http.StatusOK, map[string]any{
		"success": true,
	})
}

// Helpers
func generateRandomToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func callGeminiAPI(apiKey, prompt, lang string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-flash-latest:generateContent?key=%s", apiKey)

	systemInstruction := "You are a professional, helpful customer service virtual assistant for IMS (Innovation Massive Solutions), a professional web development agency. Keep your tone polite and friendly. You MUST answer the user in English. Do not mention that you are Google Gemini."
	if lang == "id" {
		systemInstruction = "Anda adalah asisten virtual layanan pelanggan yang profesional dan membantu untuk IMS (Innovation Massive Solutions), sebuah agensi pengembangan web profesional. Jaga nada bicara Anda tetap sopan dan ramah. Anda WAJIB menjawab pengguna dalam Bahasa Indonesia. Jangan menyebutkan bahwa Anda adalah Google Gemini."
	}

	payload := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]any{
					{
						"text": fmt.Sprintf("System Instruction: %s\n\nUser Prompt: %s", systemInstruction, prompt),
					},
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody bytes.Buffer
		_, _ = errBody.ReadFrom(resp.Body)
		return "", fmt.Errorf("gemini api returned status %d: %s", resp.StatusCode, errBody.String())
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
		return result.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("no reply from gemini api")
}
