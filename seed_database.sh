#!/bin/bash

# Скрипт для заполнения базы данных начальными данными

# Проверяем, что PostgreSQL запущен
if ! pg_isready -h localhost -p 5432 > /dev/null 2>&1; then
    echo "PostgreSQL не запущен. Запустите базу данных сначала."
    exit 1
fi

echo "Очищаем базу данных..."

# Подключаемся к базе данных и очищаем все таблицы
PGPASSWORD=dalnoboy_password psql -h localhost -U dalnoboy -d dalnoboy << EOF

-- Отключаем проверку внешних ключей для очистки
SET session_replication_role = replica;

-- Очищаем все таблицы в правильном порядке (сначала зависимые)
TRUNCATE TABLE orders CASCADE;
TRUNCATE TABLE drivers CASCADE;
TRUNCATE TABLE customers CASCADE;
TRUNCATE TABLE cities CASCADE;

-- Включаем проверку внешних ключей
SET session_replication_role = DEFAULT;

EOF

echo "База данных очищена. Заполняем начальными данными..."

# Подключаемся к базе данных и выполняем вставку данных
PGPASSWORD=dalnoboy_password psql -h localhost -U dalnoboy -d dalnoboy << EOF

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
    ('Уфа'),
    ('Волгоград'),
    ('Пермь'),
    ('Воронеж'),
    ('Краснодар'),
    ('Саратов'),
    ('Тюмень'),
    ('Тольятти'),
    ('Ижевск'),
    ('Барнаул'),
    ('Ульяновск'),
    ('Иркутск'),
    ('Хабаровск'),
    ('Ярославль'),
    ('Владивосток'),
    ('Махачкала'),
    ('Томск'),
    ('Оренбург'),
    ('Кемерово'),
    ('Рязань'),
    ('Астрахань');

-- Вставляем 30 тестовых водителей
INSERT INTO drivers (name, telegram_id, telegram_tag, city_uuid) VALUES 
    ('Иван Петров', 123456789, '@ivan_petrov', (SELECT uuid FROM cities WHERE name = 'Москва' LIMIT 1)),
    ('Алексей Сидоров', 987654321, '@alex_sidorov', (SELECT uuid FROM cities WHERE name = 'Санкт-Петербург' LIMIT 1)),
    ('Дмитрий Козлов', 555666777, '@dmitry_kozlov', (SELECT uuid FROM cities WHERE name = 'Москва' LIMIT 1)),
    ('Сергей Волков', 111222333, '@sergey_volkov', (SELECT uuid FROM cities WHERE name = 'Новосибирск' LIMIT 1)),
    ('Андрей Морозов', 444555666, '@andrey_morozov', (SELECT uuid FROM cities WHERE name = 'Екатеринбург' LIMIT 1)),
    ('Михаил Соколов', 777888999, '@mikhail_sokolov', (SELECT uuid FROM cities WHERE name = 'Казань' LIMIT 1)),
    ('Владимир Попов', 123123123, '@vladimir_popov', (SELECT uuid FROM cities WHERE name = 'Нижний Новгород' LIMIT 1)),
    ('Николай Лебедев', 456456456, '@nikolay_lebedev', (SELECT uuid FROM cities WHERE name = 'Челябинск' LIMIT 1)),
    ('Павел Козлов', 789789789, '@pavel_kozlov', (SELECT uuid FROM cities WHERE name = 'Самара' LIMIT 1)),
    ('Александр Новиков', 321321321, '@alexander_novikov', (SELECT uuid FROM cities WHERE name = 'Ростов-на-Дону' LIMIT 1)),
    ('Евгений Морозов', 654654654, '@evgeny_morozov', (SELECT uuid FROM cities WHERE name = 'Уфа' LIMIT 1)),
    ('Виктор Петров', 987987987, '@viktor_petrov', (SELECT uuid FROM cities WHERE name = 'Волгоград' LIMIT 1)),
    ('Геннадий Сидоров', 147147147, '@gennady_sidorov', (SELECT uuid FROM cities WHERE name = 'Пермь' LIMIT 1)),
    ('Борис Волков', 258258258, '@boris_volkov', (SELECT uuid FROM cities WHERE name = 'Воронеж' LIMIT 1)),
    ('Валерий Соколов', 369369369, '@valery_sokolov', (SELECT uuid FROM cities WHERE name = 'Краснодар' LIMIT 1)),
    ('Юрий Попов', 159159159, '@yury_popov', (SELECT uuid FROM cities WHERE name = 'Саратов' LIMIT 1)),
    ('Анатолий Лебедев', 357357357, '@anatoly_lebedev', (SELECT uuid FROM cities WHERE name = 'Тюмень' LIMIT 1)),
    ('Степан Козлов', 951951951, '@stepan_kozlov', (SELECT uuid FROM cities WHERE name = 'Тольятти' LIMIT 1)),
    ('Аркадий Новиков', 753753753, '@arkady_novikov', (SELECT uuid FROM cities WHERE name = 'Ижевск' LIMIT 1)),
    ('Тимофей Морозов', 159753159, '@timofey_morozov', (SELECT uuid FROM cities WHERE name = 'Барнаул' LIMIT 1)),
    ('Роман Петров', 753159753, '@roman_petrov', (SELECT uuid FROM cities WHERE name = 'Ульяновск' LIMIT 1)),
    ('Даниил Сидоров', 951357951, '@daniel_sidorov', (SELECT uuid FROM cities WHERE name = 'Иркутск' LIMIT 1)),
    ('Максим Волков', 357951357, '@maxim_volkov', (SELECT uuid FROM cities WHERE name = 'Хабаровск' LIMIT 1)),
    ('Артем Соколов', 159357159, '@artem_sokolov', (SELECT uuid FROM cities WHERE name = 'Ярославль' LIMIT 1)),
    ('Кирилл Попов', 753951753, '@kirill_popov', (SELECT uuid FROM cities WHERE name = 'Владивосток' LIMIT 1)),
    ('Игорь Лебедев', 951159951, '@igor_lebedev', (SELECT uuid FROM cities WHERE name = 'Махачкала' LIMIT 1)),
    ('Виталий Козлов', 357753357, '@vitaly_kozlov', (SELECT uuid FROM cities WHERE name = 'Томск' LIMIT 1)),
    ('Олег Новиков', 159951159, '@oleg_novikov', (SELECT uuid FROM cities WHERE name = 'Оренбург' LIMIT 1)),
    ('Семен Морозов', 753357753, '@semen_morozov', (SELECT uuid FROM cities WHERE name = 'Кемерово' LIMIT 1)),
    ('Федор Петров', 951753951, '@fedor_petrov', (SELECT uuid FROM cities WHERE name = 'Рязань' LIMIT 1)),
    ('Константин Сидоров', 357159357, '@konstantin_sidorov', (SELECT uuid FROM cities WHERE name = 'Астрахань' LIMIT 1));

-- Вставляем 30 тестовых клиентов
INSERT INTO customers (name, phone, telegram_id, telegram_tag) VALUES 
    ('ООО "Грузовик"', '+7(495)123-45-67', 111222333, '@gruzovik_company'),
    ('ИП Сидоров', '+7(812)987-65-43', 444555666, '@ip_sidorov'),
    ('АО "Транспорт"', '+7(343)555-44-33', 777888999, '@transport_ao'),
    ('ООО "Логистика Плюс"', '+7(383)111-22-33', 123456789, '@logistika_plus'),
    ('ИП Иванов', '+7(351)444-55-66', 987654321, '@ip_ivanov'),
    ('АО "Доставка Экспресс"', '+7(843)777-88-99', 555666777, '@delivery_express'),
    ('ООО "Перевозки"', '+7(831)222-33-44', 111333555, '@perevozki_company'),
    ('ИП Петров', '+7(351)666-77-88', 444777999, '@ip_petrov'),
    ('АО "Транс Сервис"', '+7(846)999-00-11', 777111333, '@trans_service'),
    ('ООО "Груз Доставка"', '+7(347)333-44-55', 222666888, '@gruz_delivery'),
    ('ИП Волков', '+7(844)555-66-77', 888111444, '@ip_volkov'),
    ('АО "Логистика"', '+7(342)777-88-99', 333555777, '@logistika_company'),
    ('ООО "Транспортные Решения"', '+7(345)111-22-33', 666999222, '@transport_solutions'),
    ('ИП Соколов', '+7(846)444-55-66', 111777444, '@ip_sokolov'),
    ('АО "Доставка"', '+7(347)777-88-99', 444111777, '@delivery_company'),
    ('ООО "Перевозчик"', '+7(844)222-33-44', 777444111, '@perevozchik'),
    ('ИП Попов', '+7(342)555-66-77', 222777555, '@ip_popov'),
    ('АО "Транс Логистика"', '+7(345)888-99-00', 555222888, '@trans_logistika'),
    ('ООО "Груз Сервис"', '+7(846)111-22-33', 888555222, '@gruz_service'),
    ('ИП Лебедев', '+7(347)444-55-66', 333888555, '@ip_lebedev'),
    ('АО "Логистика Сервис"', '+7(844)777-88-99', 666333999, '@logistika_service'),
    ('ООО "Транспорт Плюс"', '+7(342)000-11-22', 999666333, '@transport_plus'),
    ('ИП Козлов', '+7(345)333-44-55', 222999666, '@ip_kozlov'),
    ('АО "Доставка Сервис"', '+7(846)666-77-88', 555333999, '@delivery_service'),
    ('ООО "Перевозки Плюс"', '+7(347)999-00-11', 888666333, '@perevozki_plus'),
    ('ИП Новиков', '+7(844)222-33-44', 333999666, '@ip_novikov'),
    ('АО "Транс Доставка"', '+7(342)555-66-77', 666333999, '@trans_delivery'),
    ('ООО "Груз Логистика"', '+7(345)888-99-00', 999666333, '@gruz_logistika'),
    ('ИП Морозов', '+7(846)111-22-33', 222999666, '@ip_morozov'),
    ('АО "Логистика Доставка"', '+7(347)444-55-66', 555333999, '@logistika_delivery'),
    ('ООО "Транспорт Сервис"', '+7(844)777-88-99', 888666333, '@transport_service_company');

-- Вставляем 30 тестовых заказов
INSERT INTO orders (customer_uuid, title, description, weight_kg, length_cm, width_cm, height_cm, from_city_uuid, from_address, to_city_uuid, to_address, tags, price) VALUES 
    ((SELECT uuid FROM customers WHERE name = 'ООО "Грузовик"' LIMIT 1), 'Доставка мебели', 'Диван и кресло, аккуратная перевозка', 50.5, 200, 80, 60, (SELECT uuid FROM cities WHERE name = 'Москва' LIMIT 1), 'ул. Ленина, 1', (SELECT uuid FROM cities WHERE name = 'Санкт-Петербург' LIMIT 1), 'Невский пр., 100', ARRAY['мебель', 'хрупкий'], 2500),
    ((SELECT uuid FROM customers WHERE name = 'ИП Сидоров' LIMIT 1), 'Перевозка техники', 'Стиральная машина и холодильник', 120.0, 180, 70, 70, (SELECT uuid FROM cities WHERE name = 'Екатеринбург' LIMIT 1), 'ул. Мира, 50', (SELECT uuid FROM cities WHERE name = 'Москва' LIMIT 1), 'ул. Пушкина, 10', ARRAY['техника', 'крупногабарит'], 3500),
    ((SELECT uuid FROM customers WHERE name = 'АО "Транспорт"' LIMIT 1), 'Доставка документов', 'Коробка с документами', 5.0, 30, 20, 15, (SELECT uuid FROM cities WHERE name = 'Казань' LIMIT 1), 'ул. Баумана, 25', (SELECT uuid FROM cities WHERE name = 'Нижний Новгород' LIMIT 1), 'ул. Большая Покровская, 15', ARRAY['документы', 'срочно'], 800),
    ((SELECT uuid FROM customers WHERE name = 'ООО "Логистика Плюс"' LIMIT 1), 'Перевозка оборудования', 'Промышленное оборудование для завода', 500.0, 300, 150, 200, (SELECT uuid FROM cities WHERE name = 'Новосибирск' LIMIT 1), 'ул. Красная, 100', (SELECT uuid FROM cities WHERE name = 'Кемерово' LIMIT 1), 'ул. Шахтеров, 50', ARRAY['оборудование', 'промышленное'], 15000),
    ((SELECT uuid FROM customers WHERE name = 'ИП Иванов' LIMIT 1), 'Доставка продуктов', 'Свежие овощи и фрукты', 25.0, 100, 80, 60, (SELECT uuid FROM cities WHERE name = 'Челябинск' LIMIT 1), 'ул. Кирова, 75', (SELECT uuid FROM cities WHERE name = 'Екатеринбург' LIMIT 1), 'ул. Ленина, 200', ARRAY['продукты', 'скоропортящиеся'], 1200),
    ((SELECT uuid FROM customers WHERE name = 'АО "Доставка Экспресс"' LIMIT 1), 'Перевозка одежды', 'Гардероб для магазина', 80.0, 150, 100, 80, (SELECT uuid FROM cities WHERE name = 'Казань' LIMIT 1), 'ул. Баумана, 150', (SELECT uuid FROM cities WHERE name = 'Уфа' LIMIT 1), 'ул. Ленина, 300', ARRAY['одежда', 'магазин'], 2000),
    ((SELECT uuid FROM customers WHERE name = 'ООО "Перевозки"' LIMIT 1), 'Доставка стройматериалов', 'Цемент, кирпич, арматура', 2000.0, 400, 200, 150, (SELECT uuid FROM cities WHERE name = 'Самара' LIMIT 1), 'ул. Московская, 500', (SELECT uuid FROM cities WHERE name = 'Тольятти' LIMIT 1), 'ул. Юбилейная, 100', ARRAY['стройматериалы', 'тяжелый'], 8000),
    ((SELECT uuid FROM customers WHERE name = 'ИП Петров' LIMIT 1), 'Перевозка книг', 'Библиотека для школы', 150.0, 200, 150, 100, (SELECT uuid FROM cities WHERE name = 'Ростов-на-Дону' LIMIT 1), 'ул. Большая Садовая, 25', (SELECT uuid FROM cities WHERE name = 'Волгоград' LIMIT 1), 'ул. Мира, 75', ARRAY['книги', 'образование'], 1800),
    ((SELECT uuid FROM customers WHERE name = 'АО "Транс Сервис"' LIMIT 1), 'Доставка медикаментов', 'Лекарства и медицинское оборудование', 30.0, 80, 60, 40, (SELECT uuid FROM cities WHERE name = 'Пермь' LIMIT 1), 'ул. Ленина, 120', (SELECT uuid FROM cities WHERE name = 'Екатеринбург' LIMIT 1), 'ул. Мира, 80', ARRAY['медикаменты', 'хрупкий'], 2500),
    ((SELECT uuid FROM customers WHERE name = 'ООО "Груз Доставка"' LIMIT 1), 'Перевозка автомобиля', 'Легковой автомобиль на эвакуаторе', 1500.0, 500, 200, 150, (SELECT uuid FROM cities WHERE name = 'Уфа' LIMIT 1), 'ул. Ленина, 400', (SELECT uuid FROM cities WHERE name = 'Челябинск' LIMIT 1), 'ул. Кирова, 250', ARRAY['автомобиль', 'эвакуатор'], 12000),
    ((SELECT uuid FROM customers WHERE name = 'ИП Волков' LIMIT 1), 'Доставка цветов', 'Букеты и горшечные растения', 15.0, 60, 40, 50, (SELECT uuid FROM cities WHERE name = 'Волгоград' LIMIT 1), 'ул. Мира, 90', (SELECT uuid FROM cities WHERE name = 'Ростов-на-Дону' LIMIT 1), 'ул. Большая Садовая, 150', ARRAY['цветы', 'хрупкий'], 900),
    ((SELECT uuid FROM customers WHERE name = 'АО "Логистика"' LIMIT 1), 'Перевозка компьютеров', 'Офисная техника', 100.0, 120, 80, 60, (SELECT uuid FROM cities WHERE name = 'Пермь' LIMIT 1), 'ул. Ленина, 300', (SELECT uuid FROM cities WHERE name = 'Ижевск' LIMIT 1), 'ул. Советская, 100', ARRAY['компьютеры', 'техника'], 3000),
    ((SELECT uuid FROM customers WHERE name = 'ООО "Транспортные Решения"' LIMIT 1), 'Доставка спортивного инвентаря', 'Тренажеры и спортивное оборудование', 300.0, 250, 120, 100, (SELECT uuid FROM cities WHERE name = 'Воронеж' LIMIT 1), 'ул. Мира, 180', (SELECT uuid FROM cities WHERE name = 'Краснодар' LIMIT 1), 'ул. Красная, 200', ARRAY['спорт', 'оборудование'], 4500),
    ((SELECT uuid FROM customers WHERE name = 'ИП Соколов' LIMIT 1), 'Перевозка музыкальных инструментов', 'Пианино и гитары', 200.0, 180, 100, 120, (SELECT uuid FROM cities WHERE name = 'Самара' LIMIT 1), 'ул. Московская, 350', (SELECT uuid FROM cities WHERE name = 'Тольятти' LIMIT 1), 'ул. Юбилейная, 75', ARRAY['музыка', 'хрупкий'], 2800),
    ((SELECT uuid FROM customers WHERE name = 'АО "Доставка"' LIMIT 1), 'Доставка продуктов питания', 'Мясо, рыба, молочные продукты', 80.0, 120, 80, 60, (SELECT uuid FROM cities WHERE name = 'Ижевск' LIMIT 1), 'ул. Советская, 200', (SELECT uuid FROM cities WHERE name = 'Пермь' LIMIT 1), 'ул. Ленина, 450', ARRAY['продукты', 'скоропортящиеся'], 2200),
    ((SELECT uuid FROM customers WHERE name = 'ООО "Перевозчик"' LIMIT 1), 'Перевозка мебели', 'Кухонный гарнитур', 120.0, 250, 150, 80, (SELECT uuid FROM cities WHERE name = 'Барнаул' LIMIT 1), 'ул. Ленина, 100', (SELECT uuid FROM cities WHERE name = 'Новосибирск' LIMIT 1), 'ул. Красная, 300', ARRAY['мебель', 'кухня'], 3200),
    ((SELECT uuid FROM customers WHERE name = 'ИП Попов' LIMIT 1), 'Доставка электроники', 'Телевизоры и аудиосистемы', 60.0, 100, 80, 50, (SELECT uuid FROM cities WHERE name = 'Ульяновск' LIMIT 1), 'ул. Гончарова, 150', (SELECT uuid FROM cities WHERE name = 'Самара' LIMIT 1), 'ул. Московская, 400', ARRAY['электроника', 'хрупкий'], 1800),
    ((SELECT uuid FROM customers WHERE name = 'АО "Транс Логистика"' LIMIT 1), 'Перевозка промышленного оборудования', 'Станки для металлообработки', 800.0, 400, 200, 250, (SELECT uuid FROM cities WHERE name = 'Иркутск' LIMIT 1), 'ул. Ленина, 500', (SELECT uuid FROM cities WHERE name = 'Хабаровск' LIMIT 1), 'ул. Ленина, 100', ARRAY['оборудование', 'промышленное'], 25000),
    ((SELECT uuid FROM customers WHERE name = 'ООО "Груз Сервис"' LIMIT 1), 'Доставка одежды', 'Зимняя одежда для магазина', 100.0, 150, 100, 80, (SELECT uuid FROM cities WHERE name = 'Ярославль' LIMIT 1), 'ул. Советская, 200', (SELECT uuid FROM cities WHERE name = 'Владивосток' LIMIT 1), 'ул. Светланская, 150', ARRAY['одежда', 'зимняя'], 3500),
    ((SELECT uuid FROM customers WHERE name = 'ИП Лебедев' LIMIT 1), 'Перевозка продуктов', 'Консервы и крупы', 200.0, 150, 100, 80, (SELECT uuid FROM cities WHERE name = 'Махачкала' LIMIT 1), 'ул. Ленина, 300', (SELECT uuid FROM cities WHERE name = 'Ростов-на-Дону' LIMIT 1), 'ул. Большая Садовая, 400', ARRAY['продукты', 'консервы'], 2800),
    ((SELECT uuid FROM customers WHERE name = 'АО "Логистика Сервис"' LIMIT 1), 'Доставка строительных материалов', 'Доски, брус, фанера', 500.0, 300, 150, 100, (SELECT uuid FROM cities WHERE name = 'Томск' LIMIT 1), 'ул. Ленина, 250', (SELECT uuid FROM cities WHERE name = 'Новосибирск' LIMIT 1), 'ул. Красная, 500', ARRAY['стройматериалы', 'дерево'], 6000),
    ((SELECT uuid FROM customers WHERE name = 'ООО "Транспорт Плюс"' LIMIT 1), 'Перевозка бытовой техники', 'Холодильники, стиральные машины', 150.0, 180, 80, 80, (SELECT uuid FROM cities WHERE name = 'Оренбург' LIMIT 1), 'ул. Советская, 100', (SELECT uuid FROM cities WHERE name = 'Самара' LIMIT 1), 'ул. Московская, 600', ARRAY['техника', 'бытовая'], 4200),
    ((SELECT uuid FROM customers WHERE name = 'ИП Козлов' LIMIT 1), 'Доставка документов', 'Архивные документы', 10.0, 40, 30, 20, (SELECT uuid FROM cities WHERE name = 'Кемерово' LIMIT 1), 'ул. Шахтеров, 150', (SELECT uuid FROM cities WHERE name = 'Новосибирск' LIMIT 1), 'ул. Красная, 700', ARRAY['документы', 'архив'], 600),
    ((SELECT uuid FROM customers WHERE name = 'АО "Доставка Сервис"' LIMIT 1), 'Перевозка мебели', 'Офисная мебель', 180.0, 200, 120, 100, (SELECT uuid FROM cities WHERE name = 'Рязань' LIMIT 1), 'ул. Ленина, 300', (SELECT uuid FROM cities WHERE name = 'Москва' LIMIT 1), 'ул. Тверская, 50', ARRAY['мебель', 'офисная'], 4800),
    ((SELECT uuid FROM customers WHERE name = 'ООО "Перевозки Плюс"' LIMIT 1), 'Доставка продуктов', 'Свежие овощи и фрукты', 40.0, 100, 80, 60, (SELECT uuid FROM cities WHERE name = 'Астрахань' LIMIT 1), 'ул. Ленина, 200', (SELECT uuid FROM cities WHERE name = 'Волгоград' LIMIT 1), 'ул. Мира, 300', ARRAY['продукты', 'свежие'], 1500),
    ((SELECT uuid FROM customers WHERE name = 'ИП Новиков' LIMIT 1), 'Перевозка техники', 'Компьютеры и периферия', 80.0, 120, 80, 60, (SELECT uuid FROM cities WHERE name = 'Тюмень' LIMIT 1), 'ул. Ленина, 400', (SELECT uuid FROM cities WHERE name = 'Екатеринбург' LIMIT 1), 'ул. Мира, 300', ARRAY['техника', 'компьютеры'], 3200),
    ((SELECT uuid FROM customers WHERE name = 'АО "Транс Доставка"' LIMIT 1), 'Доставка мебели', 'Спальный гарнитур', 90.0, 180, 120, 80, (SELECT uuid FROM cities WHERE name = 'Тольятти' LIMIT 1), 'ул. Юбилейная, 200', (SELECT uuid FROM cities WHERE name = 'Самара' LIMIT 1), 'ул. Московская, 800', ARRAY['мебель', 'спальная'], 2800),
    ((SELECT uuid FROM customers WHERE name = 'ООО "Груз Логистика"' LIMIT 1), 'Перевозка оборудования', 'Медицинское оборудование', 250.0, 200, 150, 120, (SELECT uuid FROM cities WHERE name = 'Ижевск' LIMIT 1), 'ул. Советская, 300', (SELECT uuid FROM cities WHERE name = 'Пермь' LIMIT 1), 'ул. Ленина, 600', ARRAY['оборудование', 'медицинское'], 8500),
    ((SELECT uuid FROM customers WHERE name = 'ИП Морозов' LIMIT 1), 'Доставка продуктов', 'Мясные полуфабрикаты', 60.0, 100, 80, 60, (SELECT uuid FROM cities WHERE name = 'Барнаул' LIMIT 1), 'ул. Ленина, 400', (SELECT uuid FROM cities WHERE name = 'Новосибирск' LIMIT 1), 'ул. Красная, 900', ARRAY['продукты', 'мясо'], 2000),
    ((SELECT uuid FROM customers WHERE name = 'АО "Логистика Доставка"' LIMIT 1), 'Перевозка одежды', 'Летняя коллекция', 70.0, 120, 80, 60, (SELECT uuid FROM cities WHERE name = 'Ульяновск' LIMIT 1), 'ул. Гончарова, 300', (SELECT uuid FROM cities WHERE name = 'Самара' LIMIT 1), 'ул. Московская, 1000', ARRAY['одежда', 'летняя'], 1800),
    ((SELECT uuid FROM customers WHERE name = 'ООО "Транспорт Сервис"' LIMIT 1), 'Доставка техники', 'Кухонная техника', 110.0, 150, 100, 80, (SELECT uuid FROM cities WHERE name = 'Иркутск' LIMIT 1), 'ул. Ленина, 600', (SELECT uuid FROM cities WHERE name = 'Хабаровск' LIMIT 1), 'ул. Ленина, 300', ARRAY['техника', 'кухонная'], 3800);

EOF

echo "База данных успешно очищена и заполнена начальными данными!"
echo "Добавлено:"
echo "- 30 городов"
echo "- 30 водителей" 
echo "- 30 заказчиков"
echo "- 30 заказов" 