#!/bin/bash

# Скрипт для заполнения базы данных начальными данными

# Проверяем, что PostgreSQL запущен
if ! pg_isready -h localhost -p 5432 > /dev/null 2>&1; then
    echo "PostgreSQL не запущен. Запустите базу данных сначала."
    exit 1
fi

# Подключаемся к базе данных и выполняем вставку данных
psql -h localhost -U postgres -d dalnoboy << EOF

-- Вставляем города
INSERT INTO cities (name) VALUES 
    ('Москва'),
    ('Санкт-Петербург'),
    ('Новосибирск'),
    ('Екатеринбург'),
    ('Казань'),
    ('Нижний Новгород'),
    ('Челябинск'),
    ('Самара'),
    ('Ростов-на-Дону'),
    ('Уфа')
ON CONFLICT (name) DO NOTHING;

-- Вставляем тестовых водителей
INSERT INTO drivers (name, telegram_id, telegram_tag, city_uuid) VALUES 
    ('Иван Петров', 123456789, '@ivan_petrov', (SELECT uuid FROM cities WHERE name = 'Москва')),
    ('Алексей Сидоров', 987654321, '@alex_sidorov', (SELECT uuid FROM cities WHERE name = 'Санкт-Петербург')),
    ('Дмитрий Козлов', 555666777, '@dmitry_kozlov', (SELECT uuid FROM cities WHERE name = 'Москва'))
ON CONFLICT (telegram_id) DO NOTHING;

-- Вставляем тестовых клиентов
INSERT INTO customers (name, phone, telegram_id, telegram_tag) VALUES 
    ('ООО "Грузовик"', '+7(495)123-45-67', 111222333, '@gruzovik_company'),
    ('ИП Сидоров', '+7(812)987-65-43', 444555666, '@ip_sidorov'),
    ('АО "Транспорт"', '+7(343)555-44-33', 777888999, '@transport_ao')
ON CONFLICT DO NOTHING;

-- Вставляем тестовые заказы
INSERT INTO orders (customer_uuid, title, description, weight_kg, length_cm, width_cm, height_cm, from_location, to_location, tags, price) VALUES 
    ((SELECT uuid FROM customers WHERE name = 'ООО "Грузовик"'), 'Доставка мебели', 'Диван и кресло, аккуратная перевозка', 50.5, 200, 80, 60, 'Москва, ул. Ленина, 1', 'Москва, ул. Пушкина, 10', ARRAY['мебель', 'хрупкий'], 2500),
    ((SELECT uuid FROM customers WHERE name = 'ИП Сидоров'), 'Перевозка техники', 'Стиральная машина и холодильник', 120.0, 180, 70, 70, 'Санкт-Петербург, Невский пр., 100', 'Санкт-Петербург, Московский пр., 200', ARRAY['техника', 'крупногабарит'], 3500),
    ((SELECT uuid FROM customers WHERE name = 'АО "Транспорт"'), 'Доставка документов', 'Коробка с документами', 5.0, 30, 20, 15, 'Екатеринбург, ул. Мира, 50', 'Екатеринбург, ул. Свободы, 75', ARRAY['документы', 'срочно'], 800)
ON CONFLICT DO NOTHING;

EOF

echo "База данных успешно заполнена начальными данными!" 