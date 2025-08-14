#!/bin/bash

# Скрипт для наполнения базы данных тестовыми данными
# Использование: ./seed_database.sh

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Начинаю наполнение базы данных тестовыми данными...${NC}"

# Проверяем, запущен ли Docker Compose
if ! docker-compose ps | grep -q "postgres.*Up"; then
    echo -e "${RED}Ошибка: PostgreSQL не запущен. Запустите docker-compose up -d${NC}"
    exit 1
fi

# Функция для выполнения SQL команд
execute_sql() {
    local sql="$1"
    docker-compose exec -T postgres psql -U dalnoboy -d dalnoboy -c "$sql"
}

# Очищаем существующие данные
echo -e "${YELLOW}Очищаю существующие данные...${NC}"
execute_sql "DELETE FROM orders;"
execute_sql "DELETE FROM customers;"

# Добавляем тестовых клиентов
echo -e "${YELLOW}Добавляю тестовых клиентов...${NC}"

execute_sql "
INSERT INTO customers (name, phone, telegram_id, telegram_tag) VALUES
('Иванов Иван Иванович', '+7-900-123-45-67', 123456789, '@ivanov_ivan'),
('Петров Петр Петрович', '+7-900-234-56-78', 234567890, '@petrov_petr'),
('Сидорова Анна Владимировна', '+7-900-345-67-89', 345678901, '@sidorova_anna'),
('Козлов Дмитрий Сергеевич', '+7-900-456-78-90', 456789012, '@kozlov_dmitry'),
('Морозова Елена Александровна', '+7-900-567-89-01', 567890123, '@morozova_elena'),
('Волков Андрей Николаевич', '+7-900-678-90-12', 678901234, '@volkov_andrey'),
('Соколова Мария Игоревна', '+7-900-789-01-23', 789012345, '@sokolova_maria'),
('Лебедев Сергей Викторович', '+7-900-890-12-34', 890123456, '@lebedev_sergey'),
('Новикова Ольга Дмитриевна', '+7-900-901-23-45', 901234567, '@novikova_olga'),
('Медведев Александр Павлович', '+7-900-012-34-56', 12345678, '@medvedev_alex');
"

# Добавляем тестовые заказы
echo -e "${YELLOW}Добавляю тестовые заказы...${NC}"

execute_sql "
INSERT INTO orders (customer_uuid, title, description, weight_kg, length_cm, width_cm, height_cm, from_location, to_location, tags, price, available_from, status) 
SELECT 
    c.uuid,
    'Доставка документов',
    'Срочная доставка важных документов в офис',
    0.5,
    30,
    20,
    5,
    'Москва, ул. Тверская, д. 1',
    'Москва, ул. Арбат, д. 15',
    ARRAY['документы', 'срочно', 'офис'],
    500,
    CURRENT_DATE,
    'active'
FROM customers c WHERE c.name = 'Иванов Иван Иванович'
UNION ALL
SELECT 
    c.uuid,
    'Перевозка мебели',
    'Перевозка дивана и кресла в новую квартиру',
    80.0,
    200,
    80,
    90,
    'Москва, ул. Ленина, д. 10',
    'Москва, ул. Пушкина, д. 25',
    ARRAY['мебель', 'квартира', 'переезд'],
    3000,
    CURRENT_DATE + INTERVAL '1 day',
    'active'
FROM customers c WHERE c.name = 'Петров Петр Петрович'
UNION ALL
SELECT 
    c.uuid,
    'Доставка продуктов',
    'Доставка продуктов из магазина домой',
    15.0,
    50,
    40,
    30,
    'Москва, ТЦ Мега',
    'Москва, ул. Гагарина, д. 5',
    ARRAY['продукты', 'магазин', 'домой'],
    800,
    CURRENT_DATE,
    'active'
FROM customers c WHERE c.name = 'Сидорова Анна Владимировна'
UNION ALL
SELECT 
    c.uuid,
    'Перевозка техники',
    'Перевозка холодильника и стиральной машины',
    120.0,
    180,
    70,
    70,
    'Москва, ул. Мира, д. 20',
    'Москва, ул. Солнечная, д. 12',
    ARRAY['техника', 'холодильник', 'стиральная машина'],
    4500,
    CURRENT_DATE + INTERVAL '2 days',
    'active'
FROM customers c WHERE c.name = 'Козлов Дмитрий Сергеевич'
UNION ALL
SELECT 
    c.uuid,
    'Доставка цветов',
    'Доставка букета роз на день рождения',
    2.0,
    60,
    40,
    80,
    'Москва, Цветочный рынок',
    'Москва, ул. Романтичная, д. 8',
    ARRAY['цветы', 'подарок', 'день рождения'],
    1200,
    CURRENT_DATE,
    'active'
FROM customers c WHERE c.name = 'Морозова Елена Александровна'
UNION ALL
SELECT 
    c.uuid,
    'Перевозка строительных материалов',
    'Перевозка мешков с цементом и песком',
    500.0,
    100,
    80,
    60,
    'Москва, Строительный рынок',
    'Москва, ул. Строителей, д. 30',
    ARRAY['стройматериалы', 'цемент', 'песок'],
    8000,
    CURRENT_DATE + INTERVAL '3 days',
    'active'
FROM customers c WHERE c.name = 'Волков Андрей Николаевич'
UNION ALL
SELECT 
    c.uuid,
    'Доставка одежды',
    'Доставка заказанной одежды из интернет-магазина',
    3.0,
    40,
    30,
    20,
    'Москва, Склад интернет-магазина',
    'Москва, ул. Модная, д. 18',
    ARRAY['одежда', 'интернет-магазин', 'доставка'],
    600,
    CURRENT_DATE,
    'active'
FROM customers c WHERE c.name = 'Соколова Мария Игоревна'
UNION ALL
SELECT 
    c.uuid,
    'Перевозка спортивного инвентаря',
    'Перевозка тренажеров в спортзал',
    200.0,
    250,
    120,
    150,
    'Москва, ул. Спортивная, д. 7',
    'Москва, ул. Физкультурная, д. 22',
    ARRAY['спорт', 'тренажеры', 'спортзал'],
    6000,
    CURRENT_DATE + INTERVAL '1 day',
    'active'
FROM customers c WHERE c.name = 'Лебедев Сергей Викторович'
UNION ALL
SELECT 
    c.uuid,
    'Доставка книг',
    'Доставка заказанных книг в библиотеку',
    25.0,
    80,
    60,
    40,
    'Москва, Книжный склад',
    'Москва, ул. Книжная, д. 33',
    ARRAY['книги', 'библиотека', 'образование'],
    1500,
    CURRENT_DATE,
    'active'
FROM customers c WHERE c.name = 'Новикова Ольга Дмитриевна'
UNION ALL
SELECT 
    c.uuid,
    'Перевозка музыкальных инструментов',
    'Перевозка пианино в музыкальную школу',
    300.0,
    160,
    140,
    130,
    'Москва, ул. Музыкальная, д. 11',
    'Москва, ул. Школьная, д. 44',
    ARRAY['музыка', 'пианино', 'школа'],
    7500,
    CURRENT_DATE + INTERVAL '4 days',
    'active'
FROM customers c WHERE c.name = 'Медведев Александр Павлович';
"

# Проверяем количество добавленных записей
echo -e "${YELLOW}Проверяю количество добавленных записей...${NC}"

customers_count=$(docker-compose exec -T postgres psql -U dalnoboy -d dalnoboy -t -c "SELECT COUNT(*) FROM customers;" | tr -d ' ')
orders_count=$(docker-compose exec -T postgres psql -U dalnoboy -d dalnoboy -t -c "SELECT COUNT(*) FROM orders;" | tr -d ' ')

echo -e "${GREEN}✓ Успешно добавлено клиентов: $customers_count${NC}"
echo -e "${GREEN}✓ Успешно добавлено заказов: $orders_count${NC}"

# Показываем примеры данных
echo -e "${YELLOW}Примеры добавленных данных:${NC}"
echo -e "${YELLOW}Клиенты:${NC}"
docker-compose exec -T postgres psql -U dalnoboy -d dalnoboy -c "SELECT name, phone, telegram_tag FROM customers LIMIT 5;"

echo -e "${YELLOW}Заказы:${NC}"
docker-compose exec -T postgres psql -U dalnoboy -d dalnoboy -c "SELECT title, weight_kg, price, tags FROM orders LIMIT 5;"

echo -e "${GREEN}✓ База данных успешно наполнена тестовыми данными!${NC}" 