package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// adminMainMenuKeyboard возвращает главное меню админского бота с кнопками "Заказы", "Заказчики" и "Водители"
func adminMainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				{Text: "📋 Заказы"},
				{Text: "👥 Заказчики"},
			},
			{
				{Text: "🚚 Водители"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}

// driverMainMenuKeyboard возвращает главное меню водительского бота с кнопками заказов и уведомлений
func driverMainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				{Text: "📋 Заказы"},
			},
			{
				{Text: "🔔 Включить уведомления"},
				{Text: "🔕 Выключить уведомления"},
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
				{Text: "➕ Создать заказ"},
			},
			{
				{Text: "🟢 Активные заказы"},
				{Text: "🔴 Архивные заказы"},
			},
			// Закомментировано - убираем фильтры
			// {
			// 	{Text: "⚙️ Фильтр"},
			// 	{Text: "⬅️ Назад"},
			// },
			{
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
			// Закомментировано - убираем фильтры
			// {
			// 	{Text: "⚙️ Фильтр"},
			// 	{Text: "⬅️ Назад"},
			// },
			{
				{Text: "⬅️ Назад"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}

// filterMenuKeyboard возвращает меню настройки фильтров
// Закомментировано - убираем функционал фильтров
/*
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
*/

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

// driversMenuKeyboard возвращает меню для раздела водителей
func driversMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
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
