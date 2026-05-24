package db

import (
	"github.com/FruKun/shop_tg_bot_urfu/models"
)

func (s *Database) GetProduct(id int64) (*models.Product, error) {
	product := &models.Product{}
	err := s.db.QueryRow(
		"SELECT id, name, description, price, quantity, image_url FROM products WHERE id=?", id,
	).Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Quantity, &product.ImageUrl)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *Database) CreateProduct(product *models.Product) error {
	_, err := s.db.Exec(
		"INSERT INTO products (name, description, price, quantity, image_url) VALUES (?,?,?,?,?)",
		product.Name, product.Description, product.Price, product.Quantity, product.ImageUrl,
	)
	return err
}

func (s *Database) UpdateProduct(product *models.Product) error {
	_, err := s.db.Exec(
		"UPDATE products SET name=?, description=?, price=?, quantity=?, image_url=? WHERE id=?",
		product.Name, product.Description, product.Price, product.Quantity, product.ImageUrl, product.Id,
	)
	return err
}

func (s *Database) GetAllProduct() ([]models.Product, error) {
	rows, err := s.db.Query(
		"SELECT id, name, description, price, quantity, image_url FROM products",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var f models.Product
		if err := rows.Scan(&f.Id, &f.Name, &f.Description, &f.Price, &f.Quantity, &f.ImageUrl); err != nil {
			return nil, err
		}
		products = append(products, f)
	}
	return products, nil
}

func (s *Database) DeleteProduct(id int64) error {
	_, err := s.db.Exec("DELETE FROM products WHERE id = ?", id)
	return err
}

func (s *Database) IsAdmin(userID int64) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM admins WHERE user_id = ?)", userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
