package db

import (
	"fmt"

	"github.com/FruKun/shop_tg_bot_urfu/models"
)

func (s *Database) GetCart(userID int64) ([]models.CartItem, error) {
	rows, err := s.db.Query("SELECT product_id, quantity FROM cart_items WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		item.UserID = userID
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
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

func (s *Database) AddToCart(userID int64, productID int64, quantity int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("internal error %w", err)
	}
	defer tx.Rollback()

	var stock int

	err = tx.QueryRow("SELECT quantity FROM products WHERE id = ?", productID).Scan(&stock)
	if err != nil {
		return fmt.Errorf("товар не найден")
	}

	inCart := 0

	err = tx.QueryRow("SELECT quantity FROM cart_items WHERE user_id = ? AND product_id = ?", userID, productID).Scan(&inCart)

	if quantity+inCart > stock {
		return fmt.Errorf("недостаточно товара, в корзине: %d, в каталоге: %d", inCart, stock)
	}

	_, err = tx.Exec(
		`INSERT OR REPLACE INTO cart_items 
		(user_id, product_id, quantity) VALUES 
		(?, ?, COALESCE(
		(SELECT quantity FROM cart_items WHERE user_id=? AND product_id=?),
		0) + ?)`,
		userID, productID, userID, productID, quantity,
	)
	if err != nil {
		return fmt.Errorf("ошибка добавления в корзину %w", err)
	}

	return tx.Commit()
}

func (s *Database) Order(UserId int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("internal error: %w", err)
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT product_id, quantity FROM cart_items WHERE user_id=?", UserId)
	if err != nil {
		return fmt.Errorf("ошибка получения корзины: %w", err)
	}
	defer rows.Close()

	var itemsCount int
	for rows.Next() {
		var cartProduct models.CartItem
		if err := rows.Scan(&cartProduct.ProductID, &cartProduct.Quantity); err != nil {
			return fmt.Errorf("ошибка получения корзины: %w", err)
		}

		var productQuantity int

		err := tx.QueryRow(
			"SELECT quantity FROM products WHERE id=?", cartProduct.ProductID,
		).Scan(&productQuantity)
		if err != nil {
			return fmt.Errorf("товар не найден id=%d", cartProduct.ProductID)
		}

		if cartProduct.Quantity > productQuantity {
			return fmt.Errorf("недостаточно товара: id:%d в корзине %d, в каталоге: %d", cartProduct.ProductID, cartProduct.Quantity, productQuantity)
		}

		_, err = tx.Exec(
			"UPDATE products SET quantity = quantity - ? WHERE id = ?",
			cartProduct.Quantity, cartProduct.ProductID,
		)
		if err != nil {
			return fmt.Errorf("ошибка списания товара id=%d\n%w", cartProduct.ProductID, err)
		}

		itemsCount++
	}
	if itemsCount == 0 {
		return fmt.Errorf("корзина пуста")
	}
	// сдесь должно быть формирование заказа

	_, err = tx.Exec("DELETE FROM cart_items WHERE user_id=?", UserId)
	if err != nil {
		return fmt.Errorf("ошибка очистки корзины: %w", err)
	}

	return tx.Commit()
}
