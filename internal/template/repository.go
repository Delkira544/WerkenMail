package template

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, tpl *Template, vars []TemplateVariable) error
	GetByID(ctx context.Context, templateID uuid.UUID) (*Template, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]Template, error)
	GetVariableCounts(ctx context.Context, templateIDs []uuid.UUID) (map[uuid.UUID]int, error)
	GetVariablesByTemplateID(ctx context.Context, templateID uuid.UUID) ([]TemplateVariable, error)
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
            (id, template_id, key, type, default_value)
            VALUES (:id, :template_id, :key, :type, :default_value)`, vars[i])
		if err != nil {
			return err
		} // rollback automático por el defer
	}

	return tx.Commit()
}

func (r *repository) GetByID(ctx context.Context, templateID uuid.UUID) (*Template, error) {
	var tpl Template
	query := `SELECT id, project_id, name, subject, body_html, body_text, version, created_at, updated_at
			  FROM templates WHERE id = $1`
	err := r.db.GetContext(ctx, &tpl, query, templateID)
	if err != nil {
		return nil, err
	}
	return &tpl, nil
}

func (r *repository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]Template, error) {
	var templates []Template
	query := `SELECT id, project_id, name, subject, body_html, body_text, version, created_at, updated_at
			  FROM templates WHERE project_id = $1`
	err := r.db.SelectContext(ctx, &templates, query, projectID)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *repository) GetVariableCounts(ctx context.Context, templateIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	if len(templateIDs) == 0 {
		return map[uuid.UUID]int{}, nil
	}

	query, args, err := sqlx.In(
		`SELECT template_id, COUNT(*) AS count
         FROM template_variables
         WHERE template_id IN (?)
         GROUP BY template_id`,
		templateIDs,
	)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query) // convierte los "?" de sqlx.In a $1, $2... de postgres

	var rows []templateVariableCount
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	counts := make(map[uuid.UUID]int, len(rows))
	for _, row := range rows {
		counts[row.TemplateID] = row.Count
	}
	return counts, nil
}

func (r *repository) GetVariablesByTemplateID(ctx context.Context, templateID uuid.UUID) ([]TemplateVariable, error) {
	var vars []TemplateVariable
	query := `SELECT id, template_id, key, type, default_value
			  FROM template_variables WHERE template_id = $1`
	err := r.db.SelectContext(ctx, &vars, query, templateID)
	if err != nil {
		return nil, err
	}
	return vars, nil
}
