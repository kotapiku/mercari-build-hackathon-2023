package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kotapiku/mercari-build-hackathon-2023/backend/domain"
)

type UserRepository interface {
	AddUser(ctx context.Context, user domain.User) (int64, error)
	GetUser(ctx context.Context, id int64) (domain.User, error)
	GetUserByName(ctx context.Context, userName string) (domain.User, error)
	UpdateBalance(ctx context.Context, id int64, balance int64) error
}

type UserDBRepository struct {
	*sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserDBRepository{DB: db}
}

var (
	ErrConflict = errors.New("id conflict occurs")
)

func (r *UserDBRepository) AddUser(ctx context.Context, user domain.User) (int64, error) {
	rst, err := r.ExecContext(ctx, "INSERT OR ABORT INTO users (name, password) VALUES (?, ?)", user.Name, user.Password)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.name" {
			return 0, ErrConflict
		}
		return 0, err
	}

	lastID, err := rst.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (r *UserDBRepository) GetUser(ctx context.Context, id int64) (domain.User, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", id)

	var user domain.User
	return user, row.Scan(&user.ID, &user.Name, &user.Password, &user.Balance)
}

func (r *UserDBRepository) GetUserByName(ctx context.Context, userName string) (domain.User, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM users WHERE name = ?", userName)

	var user domain.User
	return user, row.Scan(&user.ID, &user.Name, &user.Password, &user.Balance)
}

func (r *UserDBRepository) UpdateBalance(ctx context.Context, id int64, balance int64) error {
	if _, err := r.ExecContext(ctx, "UPDATE users SET balance = ? WHERE id = ?", balance, id); err != nil {
		return err
	}
	return nil
}

type ItemRepository interface {
	AddItem(ctx context.Context, item domain.Item) (int64, error)
	GetItem(ctx context.Context, id int32) (domain.Item, error)
	GetItemImage(ctx context.Context, id int32) ([]byte, error)
	GetItems(ctx context.Context, status domain.ItemStatus) ([]domain.ItemWithCategory, error)
	GetItemsByUserID(ctx context.Context, userID int64) ([]domain.Item, error)
	GetCategory(ctx context.Context, id int64) (domain.Category, error)
	GetCategories(ctx context.Context) ([]domain.Category, error)
	UpdateItemStatus(ctx context.Context, id int32, status domain.ItemStatus) error
	SearchItem(ctx context.Context, itemName string) ([]domain.ItemWithCategory, error)
}

type ItemDBRepository struct {
	*sql.DB
}

func NewItemRepository(db *sql.DB) ItemRepository {
	return &ItemDBRepository{DB: db}
}

func (r *ItemDBRepository) AddItem(ctx context.Context, item domain.Item) (int64, error) {
	rst, err := r.ExecContext(ctx, "INSERT INTO items (name, price, description, category_id, seller_id, image, status) VALUES (?, ?, ?, ?, ?, ?, ?)", item.Name, item.Price, item.Description, item.CategoryID, item.UserID, item.Image, item.Status)
	if err != nil {
		return 0, err
	}
	lastID, err2 := rst.LastInsertId()
	if err2 != nil {
		return 0, ErrConflict // idのconflictがおきたとき
	}

	return lastID, nil
}

func (r *ItemDBRepository) GetItem(ctx context.Context, id int32) (domain.Item, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM items WHERE id = ?", id)

	var item domain.Item
	return item, row.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &item.CategoryID, &item.UserID, &item.Image, &item.Status, &item.CreatedAt, &item.UpdatedAt)
}

const selectItemsWithCat = `
		SELECT *
		FROM items
		LEFT OUTER JOIN category
		ON items.category_id = category.id
		`

func (r *ItemDBRepository) SearchItem(ctx context.Context, itemName string) ([]domain.ItemWithCategory, error) {
	rows, err := r.QueryContext(ctx, selectItemsWithCat+"WHERE items.name LIKE ?", "%"+itemName+"%")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := make([]domain.ItemWithCategory, 0)
	for rows.Next() {
		var item domain.Item
		var category domain.Category
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &item.CategoryID, &item.UserID, &item.Image, &item.Status, &item.CreatedAt, &item.UpdatedAt, &category.ID, &category.Name); err != nil {
			return nil, err
		}
		items = append(items, domain.ItemWithCategory{Item: item, Category: category})
	}
	if rows.Err() != nil {
		return nil, err
	}
	return items, nil
}

func (r *ItemDBRepository) GetItemImage(ctx context.Context, id int32) ([]byte, error) {
	row := r.QueryRowContext(ctx, "SELECT image FROM items WHERE id = ?", id)
	var image []byte
	return image, row.Scan(&image)
}

func (r *ItemDBRepository) GetItems(ctx context.Context, status domain.ItemStatus) ([]domain.ItemWithCategory, error) {
	rows, err := r.QueryContext(ctx, selectItemsWithCat+"WHERE status = ? ORDER BY updated_at desc", status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.ItemWithCategory, 0)
	for rows.Next() {
		var item domain.Item
		var category domain.Category
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &item.CategoryID, &item.UserID, &item.Image, &item.Status, &item.CreatedAt, &item.UpdatedAt, &category.ID, &category.Name); err != nil {
			return nil, err
		}
		items = append(items, domain.ItemWithCategory{Item: item, Category: category})
	}
	if rows.Err() != nil {
		return nil, err
	}
	return items, nil
}

func (r *ItemDBRepository) GetItemsByUserID(ctx context.Context, userID int64) ([]domain.Item, error) {
	rows, err := r.QueryContext(ctx, "SELECT * FROM items WHERE seller_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Description, &item.CategoryID, &item.UserID, &item.Image, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ItemDBRepository) UpdateItemStatus(ctx context.Context, id int32, status domain.ItemStatus) error {
	if _, err := r.ExecContext(ctx, "UPDATE items SET status = ? WHERE id = ?", status, id); err != nil {
		return err
	}
	return nil
}

func (r *ItemDBRepository) GetCategory(ctx context.Context, id int64) (domain.Category, error) {
	row := r.QueryRowContext(ctx, "SELECT * FROM category WHERE id = ?", id)

	var cat domain.Category
	return cat, row.Scan(&cat.ID, &cat.Name)
}

func (r *ItemDBRepository) GetCategories(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.QueryContext(ctx, "SELECT * FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []domain.Category
	for rows.Next() {
		var cat domain.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cats, nil
}
