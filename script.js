console.log('=== SCRIPT.JS ЗАГРУЖЕН ===');

// DOM элементы
const loadingEl = document.getElementById('loading');
const errorEl = document.getElementById('error');
const ordersEl = document.getElementById('orders');

// Загрузка заказов
async function loadOrders() {
    try {
        showLoading();
        
        const response = await fetch('/v1/orders');
        
        if (!response.ok) {
            if (response.status === 400 || response.status === 500) {
                throw new Error(`HTTP ${response.status}`);
            }
            throw new Error('Network error');
        }
        
        const orders = await response.json();
        console.log('Получены заказы:', orders); // Логируем данные
        
        if (!orders || orders.length === 0) {
            console.log('Заказов нет');
            ordersEl.innerHTML = '<div class="loading">Заказов не найдено</div>';
            showOrders();
            return;
        }
        
        displayOrders(orders);
        
    } catch (error) {
        console.error('Error loading orders:', error);
        showError();
    }
}

// Отображение заказов
function displayOrders(orders) {
    console.log('displayOrders вызван с:', orders); // Логируем вызов
    
    if (!orders || orders.length === 0) {
        console.log('Заказов нет в displayOrders');
        ordersEl.innerHTML = '<div class="loading">Заказов не найдено</div>';
        showOrders();
        return;
    }
    
    console.log('Создаем HTML для', orders.length, 'заказов');
    const ordersHTML = orders.map(order => createOrderHTML(order)).join('');
    console.log('HTML создан:', ordersHTML.substring(0, 200) + '...'); // Показываем начало HTML
    
    ordersEl.innerHTML = ordersHTML;
    showOrders();
}

// Создание HTML для заказа
function createOrderHTML(order) {
    return `
        <div class="order-card">
            <div class="order-header">
                <div class="order-id">#${order.uuid || order.id || 'N/A'}</div>
                <div class="order-title">${order.title || 'Без названия'}</div>
            </div>
            <div class="order-details">
                <div class="order-detail">
                    <strong>Описание:</strong> ${order.description || 'Не указано'}
                </div>
                <div class="order-detail">
                    <strong>Клиент:</strong> ${order.customer_name || order.customer || 'Не указано'} (${order.customer_phone || order.phone || 'нет телефона'})
                </div>
                <div class="order-detail">
                    <strong>Маршрут:</strong> ${order.from_location || order.from || 'Не указано'} → ${order.to_location || order.to || 'Не указано'}
                </div>
                <div class="order-detail">
                    <strong>Вес:</strong> ${order.weight ? order.weight + ' кг' : 'Не указан'}
                </div>
                <div class="order-detail">
                    <strong>Размеры:</strong> ${order.length_cm && order.width_cm && order.height_cm ? 
                        `${order.length_cm}×${order.width_cm}×${order.width_cm} см` : 
                        (order.dimensions || 'Не указаны')}
                </div>
                <div class="order-detail">
                    <strong>Цена:</strong> ${order.price ? order.price + ' ₽' : 'Не указана'}
                </div>
                <div class="order-detail">
                    <strong>Дата:</strong> ${order.created_at ? new Date(order.created_at).toLocaleDateString('ru-RU') : (order.date || 'Не указана')}
                </div>
            </div>
            ${order.tags && order.tags.length > 0 ? `<div class="order-tags">${order.tags.map(tag => `<span class="tag">${tag}</span>`).join('')}</div>` : ''}
        </div>
    `;
}

// Показать загрузку
function showLoading() {
    loadingEl.classList.remove('hidden');
    errorEl.classList.add('hidden');
    ordersEl.classList.add('hidden');
}

// Показать ошибку
function showError() {
    loadingEl.classList.add('hidden');
    errorEl.classList.remove('hidden');
    ordersEl.classList.add('hidden');
}

// Показать заказы
function showOrders() {
    loadingEl.classList.add('hidden');
    errorEl.classList.add('hidden');
    ordersEl.classList.remove('hidden');
}

// Инициализация: показываем загрузку при старте
document.addEventListener('DOMContentLoaded', function() {
    showLoading();
    loadOrders();
});

 