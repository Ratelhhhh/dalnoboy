package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// adminMainMenuKeyboard возвращает главное меню админского бота с кнопками "Заказы" и "Пользователи"
func adminMainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				{Text: "📋 Заказы"},
				{Text: "👥 Пользователи"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}

// driverMainMenuKeyboard возвращает главное меню водительского бота только с кнопкой "Заказы"
func driverMainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				{Text: "📋 Заказы"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}

// ordersMenuKeyboard возвращает меню для раздела заказов
func ordersMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				{Text: "🟢 Активные заказы"},
				{Text: "🔴 Архивные заказы"},
			},
			{
				{Text: "⚙️ Фильтр"},
				{Text: "⬅️ Назад"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}

// driverOrdersMenuKeyboard возвращает меню для раздела заказов водителей (без статусов)
func driverOrdersMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				{Text: "⚙️ Фильтр"},
				{Text: "⬅️ Назад"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}

// filterMenuKeyboard возвращает меню настройки фильтров
func filterMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				{Text: "📍 Маршрут"},
				{Text: "💰 Цена"},
			},
			{
				{Text: "📅 Дата"},
				{Text: "📦 Тип груза"},
			},
			{
				{Text: "♻️ Сбросить"},
				{Text: "⬅️ Назад"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}

// usersMenuKeyboard возвращает меню для раздела пользователей
func usersMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				{Text: "⬅️ Назад"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}
