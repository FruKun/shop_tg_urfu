package webadmin

const dashboardTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Админка магазина</title>
    <meta charset="utf-8">
    <style>
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #4CAF50; color: white; }
        .btn { padding: 5px 10px; text-decoration: none; color: white; border-radius: 3px; }
        .btn-danger { background-color: #f44336; }
        .btn-primary { background-color: #008CBA; }
        .product-img { max-width: 80px; max-height: 80px; }
    </style>
</head>
<body>
    <h1>Товары в магазине</h1>
    <a href="/admin/add" class="btn btn-primary">Добавить товар</a>
    <table>
        <tr>
            <th>ID</th><th>Изображение</th><th>Название</th><th>Цена</th><th>Количество</th><th>Действия</th>
        </tr>
        {{range .}}
        <tr>
            <td>{{.Id}}</td>
            <td>
                {{if .ImageUrl}}
                <a href="{{.ImageUrl}}"><img src="{{.ImageUrl}}" class="product-img" alt="{{.Name}}"> </a>
                {{else}}
                —
                {{end}}
            </td>
            <td>{{.Name}}</td>
            <td>{{.Price}} руб.</td>
            <td>{{.Quantity}}</td>
            <td>
                <a href="/admin/delete?id={{.Id}}" class="btn btn-danger" onclick="return confirm('Удалить товар?')">Удалить</a>
            </td>
        </tr>
        {{end}}
    </table>
</body>
</html>`

const addTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Добавить товар</title>
    <meta charset="utf-8">
</head>
<body>
    <h1>Добавление товара</h1>
    <form action="/admin/add/product" method="POST" enctype="multipart/form-data">
        <label>Название: <input type="text" name="name" required></label><br>
        <label>Описание: <textarea name="description"></textarea></label><br>
        <label>Цена: <input type="number" step="0.01" name="price" required></label><br>
        <label>Количество: <input type="number" name="quantity" required></label><br>
        <label>Изображение: <input type="file" name="image" accept="image/*"></label><br>
        <button type="submit">Добавить</button>
    </form>
    <a href="/admin/dashboard">Назад</a>
</body>
</html>`
