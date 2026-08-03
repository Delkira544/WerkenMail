package template

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, tpl *Template, vars []TemplateVariable) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, tpl *Template, vars []TemplateVariable) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `INSERT INTO templates (id, project_id, name, subject, body_html, body_text, version)
              VALUES (:id, :project_id, :name, :subject, :body_html, :body_text, :version)
              RETURNING created_at, updated_at`
	rows, err := tx.NamedQuery(query, tpl)
	if err != nil {
		return err
	}
	if rows.Next() {
		rows.StructScan(tpl)
	}
	rows.Close()

	for i := range vars {
		vars[i].TemplateID = tpl.ID
		_, err := tx.NamedExec(`INSERT INTO template_variables
            (id, template_id, key, type, required, default_value)
            VALUES (:id, :template_id, :key, :type, :required, :default_value)`, vars[i])
		if err != nil {
			return err
		} // rollback automático por el defer
	}

	return tx.Commit()
}
