package repository

import (
	"database/sql"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/domain"
	_ "github.com/lib/pq"
)

type DroneRepository interface {
	CreateDrone(drone *domain.Drone) error
	GetDroneByID(id string) (*domain.Drone, error)
	GetDroneByName(name string) (*domain.Drone, error)
	GetIdleDrones() ([]*domain.Drone, error)
	GetActiveDrones() ([]*domain.Drone, error)
	UpdateDrone(drone *domain.Drone) error
}

type OrderRepository interface {
	CreateOrder(order *domain.Order) error
	GetOrderByID(id string) (*domain.Order, error)
	GetActiveOrderByDroneID(droneID string) (*domain.Order, error)
	GetNextPendingOrder() (*domain.Order, error)
	ClaimNextPendingOrder(droneID string) (*domain.Order, error)
	UpdateOrder(order *domain.Order) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(connStr string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

// Close closes the database connection
func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

// --- Drone Implementation ---

func (r *PostgresRepository) CreateDrone(drone *domain.Drone) error {
	query := `INSERT INTO drones (id, name, status, latitude, longitude, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, drone.ID, drone.Name, drone.Status, drone.Latitude, drone.Longitude, drone.CreatedAt)
	return err
}

func (r *PostgresRepository) GetDroneByID(id string) (*domain.Drone, error) {
	query := `SELECT id, name, status, latitude, longitude, created_at FROM drones WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var drone domain.Drone
	err := row.Scan(&drone.ID, &drone.Name, &drone.Status, &drone.Latitude, &drone.Longitude, &drone.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &drone, err
}

func (r *PostgresRepository) GetDroneByName(name string) (*domain.Drone, error) {
	query := `SELECT id, name, status, latitude, longitude, created_at FROM drones WHERE name = $1`
	row := r.db.QueryRow(query, name)

	var drone domain.Drone
	err := row.Scan(&drone.ID, &drone.Name, &drone.Status, &drone.Latitude, &drone.Longitude, &drone.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &drone, err
}

func (r *PostgresRepository) GetIdleDrones() ([]*domain.Drone, error) {
	query := `SELECT id, name, status, latitude, longitude, created_at, updated_at FROM drones WHERE status = 'IDLE'`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drones []*domain.Drone
	for rows.Next() {
		var drone domain.Drone
		err := rows.Scan(&drone.ID, &drone.Name, &drone.Status, &drone.Latitude, &drone.Longitude, &drone.CreatedAt, &drone.UpdatedAt)
		if err != nil {
			return nil, err
		}
		drones = append(drones, &drone)
	}
	return drones, nil
}

func (r *PostgresRepository) GetActiveDrones() ([]*domain.Drone, error) {
	query := `SELECT id, name, status, latitude, longitude, created_at, updated_at FROM drones WHERE status IN ('IDLE', 'DELIVERING')`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drones []*domain.Drone
	for rows.Next() {
		var drone domain.Drone
		err := rows.Scan(&drone.ID, &drone.Name, &drone.Status, &drone.Latitude, &drone.Longitude, &drone.CreatedAt, &drone.UpdatedAt)
		if err != nil {
			return nil, err
		}
		drones = append(drones, &drone)
	}
	return drones, nil
}

func (r *PostgresRepository) UpdateDrone(drone *domain.Drone) error {
	query := `UPDATE drones SET status = $1, latitude = $2, longitude = $3 WHERE id = $4`
	_, err := r.db.Exec(query, drone.Status, drone.Latitude, drone.Longitude, drone.ID)
	return err
}

// --- Order Implementation ---

func (r *PostgresRepository) CreateOrder(order *domain.Order) error {
	query := `INSERT INTO orders (id, status, origin_lat, origin_lon, dest_lat, dest_lon, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, order.ID, order.Status, order.OriginLat, order.OriginLon, order.DestLat, order.DestLon, order.CreatedAt, order.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetOrderByID(id string) (*domain.Order, error) {
	query := `SELECT id, status, origin_lat, origin_lon, dest_lat, dest_lon, drone_id, created_at, updated_at FROM orders WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var order domain.Order
	err := row.Scan(&order.ID, &order.Status, &order.OriginLat, &order.OriginLon, &order.DestLat, &order.DestLon, &order.DroneID, &order.CreatedAt, &order.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &order, err
}

func (r *PostgresRepository) GetActiveOrderByDroneID(droneID string) (*domain.Order, error) {
	query := `SELECT id, status, origin_lat, origin_lon, dest_lat, dest_lon, drone_id, created_at, updated_at 
	          FROM orders WHERE drone_id = $1 AND status IN ('RESERVED', 'PICKED_UP') LIMIT 1`
	row := r.db.QueryRow(query, droneID)

	var order domain.Order
	err := row.Scan(&order.ID, &order.Status, &order.OriginLat, &order.OriginLon, &order.DestLat, &order.DestLon, &order.DroneID, &order.CreatedAt, &order.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &order, err
}

func (r *PostgresRepository) GetNextPendingOrder() (*domain.Order, error) {
	query := `SELECT id, status, origin_lat, origin_lon, dest_lat, dest_lon, drone_id, created_at, updated_at 
	          FROM orders WHERE status = 'PENDING' ORDER BY created_at ASC LIMIT 1`
	row := r.db.QueryRow(query)

	var order domain.Order
	err := row.Scan(&order.ID, &order.Status, &order.OriginLat, &order.OriginLon, &order.DestLat, &order.DestLon, &order.DroneID, &order.CreatedAt, &order.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &order, err
}

func (r *PostgresRepository) ClaimNextPendingOrder(droneID string) (*domain.Order, error) {
	// Atomic reservation using FOR UPDATE SKIP LOCKED
	// This finds the next pending order, locks it (skipping already locked ones), and updates it.
	query := `
		UPDATE orders
		SET status = 'RESERVED', drone_id = $1, updated_at = NOW()
		WHERE id = (
			SELECT id
			FROM orders
			WHERE status = 'PENDING'
			ORDER BY created_at ASC
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		RETURNING id, status, origin_lat, origin_lon, dest_lat, dest_lon, drone_id, created_at, updated_at
	`
	row := r.db.QueryRow(query, droneID)

	var order domain.Order
	err := row.Scan(&order.ID, &order.Status, &order.OriginLat, &order.OriginLon, &order.DestLat, &order.DestLon, &order.DroneID, &order.CreatedAt, &order.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &order, err
}

func (r *PostgresRepository) UpdateOrder(order *domain.Order) error {
	query := `UPDATE orders SET status = $1, drone_id = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(query, order.Status, order.DroneID, order.UpdatedAt, order.ID)
	return err
}
