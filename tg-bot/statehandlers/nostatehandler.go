package statehandlers

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strconv"
)

type NoStateHandler struct {
	Client       *redis.Client
	Bot          *tgbotapi.BotAPI
	StateHandler *StateHandler
}

// constants with response messages
const (
	startMessage = "Привет! Я бот для создания рецептов.\n Список комманд: \n /start - начать работу с ботом \n /registration - зарегистрироваться \n /login - войти в аккаунт \n /create_recipe - создать рецепт  \n /logout - выйти из аккаунта"
)

const (
	helpCommand         = "/help"
	registrationCommand = "/signup"
	createRecipeCommand = "/create_recipe"
	logInCommand        = "/login"
	logOutCommand       = "/logout"
	startCommand        = "/start"
	cancelCommand       = "/cancel"
	recipesListCommand  = "/recipes_list"
	createMealCommand   = "/create_meal"
	mealsListCommand    = "/meals_list"
)

// HandleMessage функция для обработки команды в состоянии без состояния
func (nsh *NoStateHandler) HandleMessage(ctx context.Context, userID int64, command string) error {
	switch command {
	case startCommand:
		// Вывод сообщения о том, что пользователь уже зарегистрирован
		err := nsh.Start(ctx, userID)
		if err != nil {
			return err
		}
	case mealsListCommand:
		err := nsh.StateHandler.MealPlansList(ctx, userID)
		if err != nil {
			return err
		}
	case helpCommand:
		err := nsh.Help(ctx, userID)
		if err != nil {
			return err
		}
	case recipesListCommand:
		err := nsh.StateHandler.RecipesList(ctx, userID)
		if err != nil {
			return err
		}
	case logOutCommand:
		err := nsh.LogOut(ctx, userID)
		if err != nil {
			return err
		}
	case createMealCommand:
		nsh.StateHandler.State = CreateMealState

	case registrationCommand:
		nsh.StateHandler.State = RegistrationState
	case createRecipeCommand:
		nsh.StateHandler.State = RecipeCreationState
	case logInCommand:
		nsh.StateHandler.State = LogInState
	default:
		// Обработка неизвестной команды
		err := nsh.UnknownCommand(ctx, userID)
		if err != nil {
			return err
		}
	}
	err := nsh.StateHandler.setUserState(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

// UnknownCommand функция для обработки неизвестной команды в состоянии без состояния
func (nsh *NoStateHandler) UnknownCommand(ctx context.Context, userID int64) error {
	// Вывод сообщения о том, что команда неизвестна
	_, err := nsh.Bot.Send(tgbotapi.NewMessage(userID, "Неизвестная команда"))
	if err != nil {
		return err
	}
	return nil

}

// LogOut функция для обработки команды /logout без состояния
func (nsh *NoStateHandler) LogOut(ctx context.Context, userID int64) error {
	//if user has already logged in - log out, else - send message that user is not logged in
	// Получение состояния пользователя
	token, err := nsh.StateHandler.GetUserJWTToken(ctx, userID)
	if err != nil {
		return err
	}
	if token != "" {
		// Удаление токена из базы данных
		err = nsh.StateHandler.DeleteUserJWTToken(ctx, userID)
		if err != nil {
			return err
		}
		msg := tgbotapi.NewMessage(userID, "You logout successfully.")

		msg.ReplyMarkup = nsh.StateHandler.createMainMenu(ctx, userID)
		// Вывод сообщения о том, что пользователь вышел из аккаунта
		_, err := nsh.Bot.Send(msg)
		if err != nil {
			return err
		}
		return nil
	} else {
		// Вывод сообщения о том, что пользователь не зарегистрирован
		_, err := nsh.Bot.Send(tgbotapi.NewMessage(userID, "You are not logged in."))
		if err != nil {
			return err
		}
		return nil
	}

}

// RecipesList функция для обработки команды /recipes_list в состоянии без состояния
func (sh *StateHandler) RecipesList(ctx context.Context, userID int64) error {
	client := &http.Client{}
	token, err := sh.GetUserJWTToken(ctx, userID)
	if err != nil {
		return err
	}
	recipes, err := api.GetRecipes(client, token)
	if err != nil {
		return err
	}

	// Вывод списка рецептов в виде сообщения красиво оформленного списка
	/*
			 Recipe struct {
				Id            int     `json:"id,omitempty" db:"id"`
				Title         string  `json:"title,omitempty" db:"title"`
				Description   string  `json:"description,omitempty" db:"description"`
				IsPublic      bool    `json:"public,omitempty" db:"public"`
				Cost          float64 `json:"cost,omitempty" db:"cost"`
				TimeToPrepare int64   `json:"timeToPrepare,omitempty" db:"timeToPrepare"`
				Healthy       int     `json:"healthy,omitempty" db:"healthy"`
			}

		msg.Text += fmt.Sprintf("\n*Title*: %s", recipe.Title)
				msg.Text += fmt.Sprintf("\n*Description*: %s", recipe.Description)
				msg.Text += fmt.Sprintf("\n*Cost*: %.2f", recipe.Cost)
				msg.Text += fmt.Sprintf("\n*Time to prepare*: %d", recipe.TimeToPrepare)
				msg.Text += fmt.Sprintf("\n*Healthy(1-3)*: %d", recipe.Healthy)
	*/

	//msg := tgbotapi.NewMessage(userID, "Recipes list:")
	for _, recipe := range recipes {
		msg := tgbotapi.NewMessage(userID, "")
		msg.Text += fmt.Sprintf("\n*Title*: %s", recipe.Title)
		msg.Text += fmt.Sprintf("\n*Description*: %s", recipe.Description)
		msg.Text += fmt.Sprintf("\n*Cost*: %.2f", recipe.Cost)
		msg.Text += fmt.Sprintf("\n*Time to prepare*: %d", recipe.TimeToPrepare)
		msg.Text += fmt.Sprintf("\n*Healthy(1-3)*: %d", recipe.Healthy)

		// Создаем CallbackData с ID рецепта
		callbackData := fmt.Sprintf(strconv.Itoa(recipe.Id))

		// Создаем инлайн-кнопку с текстом и CallbackData
		button := tgbotapi.NewInlineKeyboardButtonData(recipe.Title, callbackData)

		// Создаем клавиатуру с одной кнопкой и привязываем ее к сообщению
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(button),
		)
		msg.ReplyMarkup = inlineKeyboard
		msg.ParseMode = "Markdown"

		// Отправляем сообщение с кнопкой
		_, err := sh.Bot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

// MealPlansList функция для обработки команды /meal_plans_list в состоянии без состояния
func (sh *StateHandler) MealPlansList(ctx context.Context, userID int64) error {
	client := &http.Client{}
	token, err := sh.GetUserJWTToken(ctx, userID)
	if err != nil {
		return err
	}
	mealPlans, err := api.GetMealPlans(client, token)
	if err != nil {
		return err
	}

	// Вывод списка расписаний еды в виде сообщения красиво оформленного списка
	/*
		MealPlan struct {
			ID         int          `json:"id"`
			Name       string       `json:"name"`
			StartDate  time.Time    `json:"startDate"`
			EndDate    time.Time    `json:"endDate"`
			MealTimes  []MealTime   `json:"mealTimes"`
			RecipeIDs  []int        `json:"recipeIDs"`
			RecipeData []RecipeData `json:"recipeData"`
		}

		msg.Text += fmt.Sprintf("\n*Name*: %s", mealPlan.Name)
		msg.Text += fmt.Sprintf("\n*Start Date*: %s", mealPlan.StartDate.Format("2006-01-02"))
		msg.Text += fmt.Sprintf("\n*End Date*: %s", mealPlan.EndDate.Format("2006-01-02"))
		msg.Text += fmt.Sprintf("\n*Meal Times*: %s", mealTimesString)
	*/

	for _, mealPlan := range mealPlans {
		// Создаем строку с временем приемов пищи в формате "HH:MM"

		msg := tgbotapi.NewMessage(userID, "")
		msg.Text += fmt.Sprintf("\n*Name*: %s", mealPlan.Name)
		msg.Text += fmt.Sprintf("\n*At Date*: %s", mealPlan.AtTime)
		//recipes
		msg.Text += fmt.Sprintf("\n*Id*: %s", mealPlan.Id)
		msg.Text += fmt.Sprintf("\n*Title*: %s", mealPlan.Title)
		msg.Text += fmt.Sprintf("\n*Description*: %s", mealPlan.Description)
		msg.Text += fmt.Sprintf("\n*Public*: %s", mealPlan.Public)
		msg.Text += fmt.Sprintf("\n*Cost*: %s", mealPlan.Cost)
		msg.Text += fmt.Sprintf("\n*Time to prepare*: %s", mealPlan.TimeToPrepare)
		msg.Text += fmt.Sprintf("\n*Healthy(1-3)*: %s", mealPlan.Healthy)

		/*
		   Id            int     `json:"id,omitempty" db:"id"`
		   Name          string  `json:"name,omitempty" db:"name"`
		   AtTime        string  `json:"at_time,omitempty" db:"at_time"`
		   Title         string  `json:"title,omitempty" db:"title"`
		   Description   string  `json:"description,omitempty" db:"description"`
		   Public        bool    `json:"public,omitempty" db:"public"`
		   Cost          float64 `json:"cost,omitempty" db:"cost"`
		   TimeToPrepare int     `json:"timeToPrepare,omitempty" db:"timeToPrepare"`
		   Healthy       int     `json:"healthy,omitempty" db:"healthy"`
		*/
		// Создаем CallbackData с ID расписания еды
		callbackData := fmt.Sprintf("view_meal_plan:%d", mealPlan)

		// Создаем инлайн-кнопку с текстом и CallbackData
		button := tgbotapi.NewInlineKeyboardButtonData(mealPlan.Name, callbackData)

		// Создаем клавиатуру с одной кнопкой и привязываем ее к сообщению
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(button),
		)
		msg.ReplyMarkup = inlineKeyboard
		msg.ParseMode = "Markdown"

		// Отправляем сообщение с кнопкой
		_, err := sh.Bot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

// Start функция для обработки команды /start в состоянии без состояния, что будет отображаться при входе в бота и также отображает кнопки на боте
func (nsh *NoStateHandler) Start(ctx context.Context, userID int64) error {
	msg := tgbotapi.NewMessage(userID, startMessage)

	msg.ReplyMarkup = nsh.StateHandler.createMainMenu(ctx, userID)
	_, err := nsh.Bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
	return nil
}

// Help функция для обработки команды /start в состоянии без состояния
func (nsh *NoStateHandler) Help(ctx context.Context, userID int64) error {
	// Вывод сообщения о том, что пользователь не зарегистрирован
	msg := tgbotapi.NewMessage(userID, startMessage)
	_, err := nsh.Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

// Cancel функция для обработки команды /cancel в состоянии без состояния
func (nsh *NoStateHandler) Cancel(ctx context.Context, userID int64) error {
	// Обработка команды /cancel
	// вывод сообщения о том что пользователь уже на начальном экране
	return nil
}

// HandleCallback функция для обработки callback-кнопок в состоянии без состояния
func (nsh *NoStateHandler) HandleCallback(ctx context.Context, userID int64, callbackData string) error {
	// Обработка callback-кнопок
	switch callbackData {
	case "cancel":
		// Обработка callback-кнопки "Отмена"
		log.Println("cancel")
	default:
		// Обработка неизвестной callback-кнопки
		return fmt.Errorf("unknown callback data")
	}
	return nil
}

func (sh *StateHandler) createMainMenu(ctx context.Context, userID int64) tgbotapi.ReplyKeyboardMarkup {
	loggedIn, err := sh.CheckLoggedIn(ctx, userID)
	log.Println("current login status:", loggedIn)
	if err != nil {
		return tgbotapi.ReplyKeyboardMarkup{}
	}

	if loggedIn {
		return jwtMenu()
	} else {
		return noJwtMenu()
	}
}

func jwtMenu() tgbotapi.ReplyKeyboardMarkup {
	row1 := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(helpCommand),
		tgbotapi.NewKeyboardButton(logOutCommand),
	}

	row2 := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(createRecipeCommand),
		tgbotapi.NewKeyboardButton(recipesListCommand),
		tgbotapi.NewKeyboardButton(createMealCommand),
		tgbotapi.NewKeyboardButton(mealsListCommand),
	}

	keyboard := tgbotapi.NewReplyKeyboard(row1, row2)
	keyboard.OneTimeKeyboard = false // Здесь изменено значение на false
	return keyboard
}

func noJwtMenu() tgbotapi.ReplyKeyboardMarkup {
	row1 := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(helpCommand),
	}

	row2 := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(registrationCommand),
		tgbotapi.NewKeyboardButton(logInCommand),
	}

	keyboard := tgbotapi.NewReplyKeyboard(row1, row2)
	keyboard.OneTimeKeyboard = false // Здесь изменено значение на false
	return keyboard
}

func (sh *StateHandler) createCancelKeyboard() tgbotapi.ReplyKeyboardMarkup {
	row1 := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(cancelCommand),
	}

	keyboard := tgbotapi.NewReplyKeyboard(row1)
	keyboard.OneTimeKeyboard = true
	return keyboard
}
