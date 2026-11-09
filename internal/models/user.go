package models

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Role struct {
	RoleName  string
	CanCreate bool
	CanRead   bool
	CanUpdate bool
	CanDelete bool
}

type User struct {
	ID        string
	Username  string
	Role      Role
	GroupName string
}
type UserModel struct{
	DB *sql.DB
}

func InitDB() *sql.DB {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "host=127.0.0.1 port=5432 user=postgres password=meleyke6 dbname=postgres sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func(m *UserModel) GetAllUsers() ([]User, error) {
	query := `
		SELECT 
			acc.id, 
			j.key AS meta_type,
			j.value->>'username' AS username,
			g.name AS group_name, 
			r.name AS role_name,
			BOOL_OR(a.name LIKE 'CREATE_%') AS can_create,
			BOOL_OR(a.name LIKE 'READ_%')   AS can_read,
			BOOL_OR(a.name LIKE 'UPDATE_%') AS can_update,
			BOOL_OR(a.name LIKE 'DELETE_%') AS can_delete    
		FROM accounts acc      
		JOIN account_groups ag     ON ag.account_id = acc.id AND ag.status = 1
		JOIN groups g              ON g.id = ag.group_id     AND g.status != 3
		JOIN rbac_group_roles gr   ON gr.group_id = g.id     AND gr.status = 1
		JOIN rbac_roles r          ON r.id = gr.role_id      AND r.status != 3      
		JOIN rbac_role_accesses ra ON ra.role_id = r.id      AND ra.status = 1    
		JOIN rbac_accesses a       ON a.id = ra.access_id    AND a.status != 3
		JOIN jsonb_each(acc.meta) AS j(key, value) ON true
		WHERE acc.status != 3
		  AND j.value->>'username' IS NOT NULL
		GROUP BY acc.id, j.key, j.value, g.name, r.name`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		var metaType string 
		err := rows.Scan(
			&u.ID, &metaType, &u.Username, &u.GroupName, 
			&u.Role.RoleName, &u.Role.CanCreate, &u.Role.CanRead,
			&u.Role.CanUpdate, &u.Role.CanDelete,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func(m *UserModel) GetUserByUsername(username string) ([]User, error) {
	query := `
		SELECT 
			acc.id, 
			j.key AS meta_type,
			j.value->>'username' AS username,
			g.name AS group_name, 
			r.name AS role_name,
			BOOL_OR(a.name LIKE 'CREATE_%') AS can_create,
			BOOL_OR(a.name LIKE 'READ_%')   AS can_read,
			BOOL_OR(a.name LIKE 'UPDATE_%') AS can_update,
			BOOL_OR(a.name LIKE 'DELETE_%') AS can_delete  
		FROM accounts acc      
		JOIN account_groups ag     ON ag.account_id = acc.id AND ag.status = 1
		JOIN groups g              ON g.id = ag.group_id     AND g.status != 3
		JOIN rbac_group_roles gr   ON gr.group_id = g.id     AND gr.status = 1
		JOIN rbac_roles r          ON r.id = gr.role_id      AND r.status != 3      
		JOIN rbac_role_accesses ra ON ra.role_id = r.id      AND ra.status = 1    
		JOIN rbac_accesses a       ON a.id = ra.access_id    AND a.status != 3 
		JOIN jsonb_each(acc.meta) AS j(key, value) ON true     
		WHERE acc.status != 3
		  AND j.value->>'username' IS NOT NULL
		  AND j.value->>'username' = $1
		GROUP BY acc.id, j.key, j.value, g.name, r.name`

	rows, err := m.DB.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		var metaType string 
		err := rows.Scan(
			&u.ID, &metaType, &u.Username, &u.GroupName, 
			&u.Role.RoleName, &u.Role.CanCreate, &u.Role.CanRead,
			&u.Role.CanUpdate, &u.Role.CanDelete,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}


func(m *UserModel) GetUserByUserID(id string)([]User, error){
	query := `SELECT 
						acc.id,
						j.key AS meta_type,
						j.value->>'username' AS username,
						g.name AS group_name,
						r.name AS role_name,
						BOOL_OR(a.name LIKE 'CREATE_%') AS can_create,
						BOOL_OR(a.name LIKE 'READ_%') AS can_read,
						BOOL_OR(a.name LIKE 'UPDATE_%') AS can_update,
						BOOL_OR(a.name LIKE 'DELETE_%') AS can_delete
						FROM accounts acc
						JOIN account_groups ag ON ag.account_id = acc.id AND ag.status = 1
						JOIN groups g ON g.id = ag.group_id AND g.status != 3
						JOIN rbac_group_roles gr ON gr.group_id = g.id AND gr.status = 1
						JOIN rbac_roles r ON r.id = gr.role_id  AND r.status != 3
						JOIN rbac_role_accesses ra ON ra.role_id = r.id AND ra.status = 1
						JOIN rbac_accesses a ON a.id = ra.access_id AND a.status != 3
						JOIN jsonb_each(acc.meta) AS j(key,value) ON true
						WHERE acc.status != 3
							AND j.value->>'username' IS NOT NULL
							AND acc.id IS NOT NULL
							AND acc.id = $1
						GROUP BY acc.id, j.key, j.value, g.name, r.name`


	rows, err := m.DB.Query(query, id)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next(){
		var u User
		var metaType string
		err := rows.Scan(&u.ID, &metaType, &u.Username,
										 &u.GroupName,&u.Role.RoleName, &u.Role.CanCreate, 
										 &u.Role.CanRead, &u.Role.CanUpdate, &u.Role.CanDelete,)
		if err != nil{
			return nil,err
		}

		users = append(users, u)
	}
	if err = rows.Err(); err != nil{
		return nil, err
	}
	return users, nil
}