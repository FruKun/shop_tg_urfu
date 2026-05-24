package db

import (
	"github.com/FruKun/shop_tg_bot_urfu/models"
)

func (s *Database) GetCart(userID int64) ([]models.CartItem, error) {
	rows, err := s.db.Query("SELECT user_id, product_id, quantity FROM cart_items WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		if err := rows.Scan(&item.UserID, &item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *Database) ClearCart(userID int64) error {
	_, err := s.db.Exec("DELETE FROM cart_items WHERE user_id=?", userID)
	return err
}

func (s *Database) AddToCart(userID int64, productID int, quantity int) error {
	_, err := s.db.Exec(
		`INSERT OR REPLACE INTO cart_items 
		(user_id, product_id, quantity) VALUES 
		(?, ?, COALESCE(
		(SELECT quantity FROM cart_items WHERE user_id=? AND product_id=?),
		0) + ?)`,
		userID, productID, userID, productID, quantity,
	)
	return err
}
