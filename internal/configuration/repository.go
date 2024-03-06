package configuration

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/unm4sked/finch/internal/entities"
	"github.com/unm4sked/finch/pkg/postgres"
)

type Repository interface {
	CreateConfiguration() error
	GetConfigurationById(id string) (entities.Configuration, error)
	GetConfigurations() ([]entities.Configuration, error)
	DeleteConfiguration(id string) error
	UpdateConfiguration(id string, description string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db postgres.Postgres) Repository {
	return &repository{
		db: db.Pool,
	}
}

func (r *repository) CreateConfiguration() error {
	description := "text text text"
	tag, err := r.db.Exec(context.Background(), `INSERT INTO configurations (description) VALUES ($1)`, description)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	fmt.Println(tag)
	return errors.New("hello")
}

func (r *repository) GetConfigurationById(id string) (entities.Configuration, error) {
	var deafultConfiguration entities.Configuration
	rows, err := r.db.Query(context.Background(), "SELECT * FROM configurations WHERE id LIKE $1 LIMIT 1", id)
	if err != nil {
		return deafultConfiguration, err
	}
	configuration, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[entities.Configuration])
	if err != nil {
		return deafultConfiguration, err
	}
	return configuration, nil
}

func (r *repository) GetConfigurations() ([]entities.Configuration, error) {
	emptyConfiguratios := make([]entities.Configuration, 0)
	rows, err := r.db.Query(context.Background(), "SELECT * FROM configurations")
	if err != nil {
		fmt.Println("Error: ", err)
		return emptyConfiguratios, err
	}

	defer rows.Close()

	configuratios, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.Configuration])

	if err != nil {
		fmt.Println("Error while collecting rows", err)
		return emptyConfiguratios, err
	}

	return configuratios, nil
}

func (r *repository) DeleteConfiguration(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM configurations WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateConfiguration(id string, description string) error {
	config, err := r.GetConfigurationById(id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(context.Background(), "UPDATE configurations SET description = $1 WHERE id = $2", description, config.Id)

	if err != nil {
		return err
	}

	return nil
}
