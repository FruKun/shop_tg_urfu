package db

import (
	"fmt"

	"github.com/FruKun/shop_tg_bot_urfu/models"
)

func (s *Database) GetProduct(id int64) (*models.Product, error) {
	product := &models.Product{}
	err := s.db.QueryRow(
		"SELECT id, name, description, price, quantity, image_url FROM products WHERE id=?", id,
	).Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Quantity, &product.ImageUrl)
	if err != nil {
		s.logger.Error("product id = id, error: %s", err)
		return nil, err
	}
	return product, nil
}

func (s *Database) CreateProduct(product *models.Product) error {
	_, err := s.db.Exec(
		"INSERT INTO products (name, description, price, quantity, image_url) VALUES (?,?,?,?,?)",
		product.Name, product.Description, product.Price, product.Quantity, product.ImageUrl,
	)
	if err != nil {
		s.logger.Error("%s", err)
		return err
	}
	s.logger.Debug("%s, %s, %f, %d, %s", product.Name, product.Description, product.Price, product.Quantity, product.ImageUrl)
	return nil
}

func (s *Database) UpdateProduct(product *models.Product) error {
	_, err := s.db.Exec(
		"UPDATE products SET name=?, description=?, price=?, quantity=?, image_url=? WHERE id=?",
		product.Name, product.Description, product.Price, product.Quantity, product.ImageUrl, product.Id,
	)
	if err != nil {
		s.logger.Error("%s", err)
		return err
	}
	s.logger.Debug("created product: %s, %s, %f, %d, %s", product.Name, product.Description, product.Price, product.Quantity, product.ImageUrl)
	return nil
}

func (s *Database) GetAllProducts() ([]models.Product, error) {
	rows, err := s.db.Query(
		"SELECT id, name, description, price, quantity, image_url FROM products ORDER BY id",
	)
	if err != nil {
		s.logger.Error("%s", err)
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var f models.Product
		if err := rows.Scan(&f.Id, &f.Name, &f.Description, &f.Price, &f.Quantity, &f.ImageUrl); err != nil {
			s.logger.Error("%s", err)
			return nil, err
		}
		products = append(products, f)
	}
	s.logger.Debug("created products: %d", len(products))
	return products, nil
}

func (s *Database) GetPaginatedProducts(itemsPerPage int, page int) ([]models.Product, int, error) {
	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM products WHERE quantity > 0").Scan(&total)
	if err != nil {
		s.logger.Error("%s", err)
		return nil, 0, fmt.Errorf("ошибка базы данных")
	}

	offset := page * itemsPerPage
	rows, err := s.db.Query(
		`SELECT id, name, description, price, quantity, image_url
         FROM products 
         WHERE quantity > 0 
         ORDER BY id 
         LIMIT ? OFFSET ?`,
		itemsPerPage, offset,
	)
	if err != nil {
		s.logger.Error("%s", err)
		return nil, total, fmt.Errorf("ошибка базы данных")
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.Price,
			&p.Quantity, &p.ImageUrl)
		if err != nil {
			s.logger.Error("%s", err)
			return nil, total, fmt.Errorf("ошибка базы данных")
		}
		products = append(products, p)
	}
	return products, total, nil
}

func (s *Database) DeleteProduct(id int64) error {
	_, err := s.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		s.logger.Error("%s", err)
		return err
	}
	s.logger.Debug("deleted product: %d", id)
	return nil
}
