package projects

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, input *Project) error
	Update(ctx context.Context, input *Project) error
	Delete(ctx context.Context, id string) error
	GetProjectByID(ctx context.Context, id string) (*Project, error)
	List(ctx context.Context, filter ProjectFilter) ([]Project, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

// Crea un nuevo proyecto en la base de datos
func (r *repository) Create(ctx context.Context, input *Project) error {
	query := `INSERT INTO projects (id, user_id, name, description)
			   VALUES (:id, :user_id, :name, :description)
			   RETURNING created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, input)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.StructScan(input); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) Update(ctx context.Context, input *Project) error {
	query := `UPDATE projects
			  SET name = :name, description = :description, updated_at = NOW()
			  WHERE id = :id
			  RETURNING updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, input)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.StructScan(input); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *repository) GetProjectByID(ctx context.Context, id string) (*Project, error) {
	var project Project
	query := `SELECT id, user_id, name, description, created_at, updated_at
			  FROM projects
			  WHERE id = $1`
	err := r.db.GetContext(ctx, &project, query, id)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *repository) List(ctx context.Context, filter ProjectFilter) ([]Project, error) {
	query := `SELECT id, user_id, name, description, created_at, updated_at
              FROM projects WHERE 1=1`
	args := []interface{}{}
	idx := 1
	if filter.UserID != nil {
		query += fmt.Sprintf(` AND user_id = $%d`, idx)
		args = append(args, *filter.UserID)
		idx++
	}

	if filter.Search != nil {
		query += fmt.Sprintf(` AND name ILIKE $%d`, idx) // ILIKE = case-insensitive
		args = append(args, "%"+*filter.Search+"%")
		idx++
	}
	query += ` ORDER BY created_at DESC`
	if filter.Limit != nil {
		query += fmt.Sprintf(` LIMIT $%d`, idx)
		args = append(args, *filter.Limit)
		idx++
	}
	if filter.Offset != nil {
		query += fmt.Sprintf(` OFFSET $%d`, idx)
		args = append(args, *filter.Offset)
		idx++
	}
	var projects []Project
	err := r.db.SelectContext(ctx, &projects, query, args...)
	return projects, err
}
