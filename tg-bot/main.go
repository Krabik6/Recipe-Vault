package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type RegistrationState int

// Метод MarshalBinary реализует интерфейс encoding.BinaryMarshaler
func (rs RegistrationState) MarshalBinary() ([]byte, error) {
	// Преобразуем RegistrationState в строку
	str := strconv.Itoa(int(rs))
	// Преобразуем строку в байтовый массив
	data := []byte(str)
	return data, nil
}

// Метод UnmarshalBinary реализует интерфейс encoding.BinaryUnmarshaler
func (rs *RegistrationState) UnmarshalBinary(data []byte) error {
	// Преобразуем байтовый массив в строку
	str := string(data)
	// Преобразуем строку в число типа int
	val, err := strconv.Atoi(str)
	if err != nil {
		return err
	}
	// Присваиваем полученное значение типу RegistrationState
	*rs = RegistrationState(val)
	return nil
}

const (
	StateIdle RegistrationState = iota
	StateWaitingForLogin
	StateWaitingForNickname
	StateWaitingForPassword
	StateConfirmation
)

type RecipeState int

const (
	StateRecipeIdle RecipeState = iota
	StateWaitingForRecipeName
	StateWaitingForRecipeDescription
	StateRecipeConfirmation
)

type User struct {
	UserID    int64
	Login     string
	Nickname  string
	Password  string
	RegStatus RegistrationState
}

type Recipe struct {
	UserID      int64
	RecipeName  string
	Description string
	RecipeState RecipeState
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	apiToken, ok := os.LookupEnv("TELEGRAM_TOKEN")
	if !ok {
		log.Fatal("TELEGRAM_TOKEN not found")
	}

	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	redisClient.FlushDB(context.Background()).Err() //todo its delete all

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(bot, redisClient, update)
		} else if update.CallbackQuery != nil {
			handleCallbackQuery(bot, redisClient, update)
		}
	}
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, redisClient *redis.Client, update tgbotapi.Update) {
	callbackQuery := update.CallbackQuery
	if callbackQuery != nil {
		message := callbackQuery.Message
		userID := message.Chat.ID

		user, err := getUserFromRedis(redisClient, userID)
		if err != nil {
			log.Println("Ошибка при получении пользователя из Redis:", err)
			return
		}

		recipe := &Recipe{
			UserID:      userID,
			RecipeState: StateRecipeIdle,
		}

		// Обработка нажатия на кнопку
		switch callbackQuery.Data {
		case "confirm":
			handleConfirmationState(bot, redisClient, user, recipe, update)
		case "restart":
			handleRestartState(bot, redisClient, user, recipe, update)
		case "recipe_confirm":
			handleRecipeConfirmationState(bot, redisClient, user, recipe, update)
		case "recipe_restart":
			handleRecipeRestartState(bot, redisClient, user, recipe, update)
		case "cancel":
			handleCancelState(bot, redisClient, user, recipe, update)
		}
	}
}

func handleRestartState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
	// Perform the necessary actions when the "restart" button is pressed
	// For example, reset the registration process or clear user data
	reply := "Процесс будет начат заново."
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
	bot.Send(msg)

	// Reset user registration status and remove user data from Redis
	user.RegStatus = StateIdle
	err := redisClient.Del(context.Background(), fmt.Sprintf("user:%d", user.UserID)).Err()
	if err != nil {
		log.Println("Ошибка при удалении данных пользователя из Redis:", err)
	}
}

func handleRecipeRestartState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
	// Perform the necessary actions when the "recipe_restart" button is pressed
	// For example, reset the recipe creation process or clear recipe data
	reply := "Процесс создания рецепта будет начат заново."
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
	bot.Send(msg)

	// Reset recipe creation status and remove recipe data from Redis
	recipe.RecipeState = StateRecipeIdle
	err := redisClient.Del(context.Background(), fmt.Sprintf("recipe:%d", recipe.UserID)).Err()
	if err != nil {
		log.Println("Ошибка при удалении данных рецепта из Redis:", err)
	}
}
func handleCancelState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
	// Perform the necessary actions when the "cancel" button is pressed
	// For example, cancel the current process or clear user/recipe data

	// Reset user registration status and remove user data from Redis
	user.RegStatus = StateIdle
	err := redisClient.Del(context.Background(), fmt.Sprintf("user:%d", user.UserID)).Err()
	if err != nil {
		log.Println("Ошибка при удалении данных пользователя из Redis:", err)
	}

	// Reset recipe creation status and remove recipe data from Redis
	recipe.RecipeState = StateRecipeIdle
	err = redisClient.Del(context.Background(), fmt.Sprintf("recipe:%d", recipe.UserID)).Err()
	if err != nil {
		log.Println("Ошибка при удалении данных рецепта из Redis:", err)
	}

	var chatID int64
	if update.Message != nil {
		chatID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
	}

	reply := "Процесс отменен. Вы можете начать заново или выбрать другое действие."
	msg := tgbotapi.NewMessage(chatID, reply)
	bot.Send(msg)
}

func handleMessage(bot *tgbotapi.BotAPI, redisClient *redis.Client, update tgbotapi.Update) {
	message := update.Message
	userID := message.Chat.ID

	exists, err := checkUserExistsInRedis(redisClient, userID)
	if err != nil {
		log.Println("Ошибка при проверке существования пользователя в Redis:", err)
		return
	}

	if !exists {
		// Если пользователь не найден, создаем нового пользователя
		user := &User{
			UserID:    userID,
			RegStatus: StateIdle,
		}

		// Сохраняем нового пользователя в Redis
		err := saveUserToRedis(redisClient, user)
		if err != nil {
			log.Println("Ошибка при сохранении пользователя в Redis:", err)
			return
		}
	}
	user, err := getUserFromRedis(redisClient, userID)

	recipe := &Recipe{
		UserID:      userID,
		RecipeState: StateRecipeIdle,
	}

	// Обработка текстовых сообщений
	if message != nil && message.Text != "" {
		switch user.RegStatus {
		case StateIdle:
			handleIdleState(bot, redisClient, user, recipe, update)
		case StateWaitingForLogin:
			handleLoginState(bot, redisClient, user, recipe, update)
		case StateWaitingForNickname:
			handleNicknameState(bot, redisClient, user, recipe, update)
		case StateWaitingForPassword:
			handlePasswordState(bot, redisClient, user, recipe, update)
		}

		switch recipe.RecipeState {
		case StateRecipeIdle:
			handleRecipeIdleState(bot, redisClient, user, recipe, update)
		case StateWaitingForRecipeName:
			handleRecipeNameState(bot, redisClient, user, recipe, update)
		case StateWaitingForRecipeDescription:
			handleRecipeDescriptionState(bot, redisClient, user, recipe, update)
		}
	}
}

func checkUserExistsInRedis(redisClient *redis.Client, userID int64) (bool, error) {
	hashKey := fmt.Sprintf("user:%d", userID)
	exists, err := redisClient.Exists(context.Background(), hashKey).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func getUserFromRedis(client *redis.Client, userID int64) (*User, error) {
	hashKey := fmt.Sprintf("user:%d", userID)
	log.Println(hashKey)
	result, err := client.HGetAll(context.Background(), hashKey).Result()
	if err != nil {
		return nil, err
	}
	log.Println(result)

	//userIDStr := result["UserID"]
	//if userIDStr == "" {
	//	return nil, fmt.Errorf("UserID is empty for user with ID %d", userID)
	//}

	regStatusStr := result["RegStatus"]
	regStatus, err := strconv.Atoi(regStatusStr)
	if err != nil {
		return nil, err
	}

	user := &User{
		UserID:    userID,
		Login:     result["Login"],
		Nickname:  result["Nickname"],
		Password:  result["Password"],
		RegStatus: RegistrationState(regStatus),
	}

	return user, nil
}

func saveUserToRedis(client *redis.Client, user *User) error {
	hashKey := fmt.Sprintf("user:%d", user.UserID)

	data := map[string]interface{}{
		"UserID":    user.UserID,
		"Login":     user.Login,
		"Nickname":  user.Nickname,
		"Password":  user.Password,
		"RegStatus": user.RegStatus,
	}

	err := client.HMSet(context.Background(), hashKey, data).Err()
	if err != nil {
		return err
	}

	return nil
}

func getRecipeFromRedis(client *redis.Client, recipeID int64) (*Recipe, error) {
	hashKey := fmt.Sprintf("recipe:%d", recipeID)

	result, err := client.HGetAll(context.Background(), hashKey).Result()
	if err != nil {
		return nil, err
	}

	userIDStr := result["UserID"]
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	recipeStateStr := result["RecipeState"]
	recipeState, _ := strconv.Atoi(recipeStateStr)

	recipe := &Recipe{
		UserID:      userID,
		RecipeName:  result["RecipeName"],
		Description: result["Description"],
		RecipeState: RecipeState(recipeState),
	}

	return recipe, nil
}

func saveRecipeToRedis(client *redis.Client, recipe *Recipe) error {
	hashKey := fmt.Sprintf("recipe:%d", recipe.UserID)

	data := map[string]interface{}{
		"UserID":      recipe.UserID,
		"RecipeName":  recipe.RecipeName,
		"Description": recipe.Description,
		"RecipeState": recipe.RecipeState,
	}

	err := client.HMSet(context.Background(), hashKey, data).Err()
	if err != nil {
		return err
	}

	return nil
}

func handleIdleState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
	if update.Message.Text == "/signUp" {
		reply := "Добро пожаловать в процесс регистрации!\nПожалуйста, введите свой логин:"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
		user.RegStatus = StateWaitingForLogin
		user.UserID = update.Message.Chat.ID // Сохраняем UserID пользователя
		err := saveUserToRedis(redisClient, user)
		if err != nil {
			log.Println("Ошибка при сохранении пользователя в Redis:", err)
		}
	} else if update.Message.Text == "/createRecipe" {
		reply := "Добро пожаловать в процесс создания рецепта!\nПожалуйста, введите название рецепта:"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)

		recipe.RecipeState = StateWaitingForRecipeName
		err := saveRecipeToRedis(redisClient, recipe)
		if err != nil {
			log.Println("Ошибка при сохранении рецепта в Redis:", err)
		}
	} else if update.Message.Text == "/cancel" {
		handleCancelState(bot, redisClient, user, recipe, update)
	}
}

func handleLoginState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
	login := strings.TrimSpace(update.Message.Text)
	if login == "" {
		reply := "Логин не может быть пустым. Пожалуйста, введите логин:"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
		return
	}

	user.Login = login

	reply := "Отлично! Теперь введите свой никнейм:"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	bot.Send(msg)

	user.RegStatus = StateWaitingForNickname
	err := saveUserToRedis(redisClient, user)
	if err != nil {
		log.Println("Ошибка при сохранении пользователя в Redis:", err)
	}
}

func handleNicknameState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
	nickname := strings.TrimSpace(update.Message.Text)
	if nickname == "" {
		reply := "Никнейм не может быть пустым. Пожалуйста, введите никнейм:"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
		return
	}

	user.Nickname = nickname

	reply := "Отлично! Теперь введите свой пароль:"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	bot.Send(msg)

	user.RegStatus = StateWaitingForPassword
	err := saveUserToRedis(redisClient, user)
	if err != nil {
		log.Println("Ошибка при сохранении пользователя в Redis:", err)
	}
}

func handlePasswordState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
	password := strings.TrimSpace(update.Message.Text)
	if password == "" {
		reply := "Пароль не может быть пустым. Пожалуйста, введите пароль:"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
		return
	}

	user.Password = password

	reply := fmt.Sprintf("Вы ввели следующие данные:\nЛогин: %s\nНикнейм: %s\nПароль: %s\n\nПожалуйста, подтвердите, что все верно.",
		user.Login, user.Nickname, user.Password)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да", "confirm"),
			tgbotapi.NewInlineKeyboardButtonData("Нет", "restart"),
		),
	)
	bot.Send(msg)

	user.RegStatus = StateConfirmation
	err := saveUserToRedis(redisClient, user)
	if err != nil {
		log.Println("Ошибка при сохранении пользователя в Redis:", err)
	}
}

func handleConfirmationState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "confirm":
			client := &http.Client{}
			SignUp(bot, update, client, *user)
			reply := "Регистрация успешно завершена!"
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
			bot.Send(msg)

			user.RegStatus = StateIdle
			err := redisClient.Del(context.Background(), fmt.Sprintf("user:%d", user.UserID)).Err()
			if err != nil {
				log.Println("Ошибка при удалении данных пользователя из Redis:", err)
			}

		case "restart":
			reply := "Процесс регистрации будет начат заново."
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
			bot.Send(msg)

			user.RegStatus = StateIdle
			err := redisClient.Del(context.Background(), fmt.Sprintf("user:%d", user.UserID)).Err()
			if err != nil {
				log.Println("Ошибка при удалении данных пользователя из Redis:", err)
			}
		case "cancel":
			handleCancelState(bot, redisClient, user, recipe, update)
		}
	}
}

func handleRecipeIdleState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
	if update.Message.Text == "/createRecipe" {
		reply := "Добро пожаловать в процесс создания рецепта!\nПожалуйста, введите название рецепта:"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)

		recipe.RecipeState = StateWaitingForRecipeName
		err := saveRecipeToRedis(redisClient, recipe)
		if err != nil {
			log.Println("Ошибка при сохранении рецепта в Redis:", err)
		}
	} else if update.Message.Text == "/cancel" {
		handleCancelState(bot, redisClient, user, recipe, update)
	}
}

func handleRecipeNameState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
	recipeName := strings.TrimSpace(update.Message.Text)
	if recipeName == "" {
		reply := "Название рецепта не может быть пустым. Пожалуйста, введите название рецепта:"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
		return
	}

	recipe.RecipeName = recipeName

	reply := "Отлично! Теперь введите описание рецепта:"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	bot.Send(msg)

	recipe.RecipeState = StateWaitingForRecipeDescription
	err := saveRecipeToRedis(redisClient, recipe)
	if err != nil {
		log.Println("Ошибка при сохранении рецепта в Redis:", err)
	}
}

func handleRecipeDescriptionState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
	description := strings.TrimSpace(update.Message.Text)
	if description == "" {
		reply := "Описание рецепта не может быть пустым. Пожалуйста, введите описание рецепта:"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
		return
	}

	recipe.Description = description

	reply := fmt.Sprintf("Вы ввели следующие данные для рецепта:\nНазвание: %s\nОписание: %s\n\nПожалуйста, подтвердите, что все верно.",
		recipe.RecipeName, recipe.Description)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да", "recipe_confirm"),
			tgbotapi.NewInlineKeyboardButtonData("Нет", "recipe_restart"),
		),
	)
	bot.Send(msg)

	recipe.RecipeState = StateRecipeConfirmation
	err := saveRecipeToRedis(redisClient, recipe)
	if err != nil {
		log.Println("Ошибка при сохранении рецепта в Redis:", err)
	}
}

func handleRecipeConfirmationState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "recipe_confirm":
			reply := "Рецепт успешно создан!"
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
			bot.Send(msg)

			recipe.RecipeState = StateRecipeIdle
			err := redisClient.Del(context.Background(), fmt.Sprintf("recipe:%d", recipe.UserID)).Err()
			if err != nil {
				log.Println("Ошибка при удалении данных рецепта из Redis:", err)
			}

		case "recipe_restart":
			reply := "Процесс создания рецепта будет начат заново."
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
			bot.Send(msg)

			recipe.RecipeState = StateRecipeIdle
			err := redisClient.Del(context.Background(), fmt.Sprintf("recipe:%d", recipe.UserID)).Err()
			if err != nil {
				log.Println("Ошибка при удалении данных рецепта из Redis:", err)
			}
		}
	}
}

func SignUp(bot *tgbotapi.BotAPI, update tgbotapi.Update, client *http.Client, user User) error {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Sign up")
	bot.Send(msg)

	signUpCredentials := model.SignUpCredentials{
		Username: user.Login,
		Password: user.Password,
		Name:     user.Nickname,
	}

	requestBody, err := json.Marshal(signUpCredentials)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Error while signing up ("+err.Error()+"). Please try again."))
		return err
	}

	resp, err := client.Post("http://localhost:8000/auth/sign-up", "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Error while signing up ("+err.Error()+"). Please try again. POST"))
		return err
	}
	defer resp.Body.Close()

	var authResponse model.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Error while signing up ("+err.Error()+"). Please try again."))
		return err
	}

	fmt.Println(resp.StatusCode, authResponse.Token)

	message := "You have successfully signed up. Now signIn" + authResponse.Token
	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, message))
	return nil
}
