package sql

import (
	"database/sql"
	"testing"

	"github.com/go-gorp/gorp/v3"
	"github.com/semaphoreui/semaphore/db"
	"github.com/semaphoreui/semaphore/util"
)

func TestValidatePort(t *testing.T) {
	d := SqlDb{}
	q := d.connection.prepareQueryWithDialect("select * from `test` where id = ?, email = ?", gorp.PostgresDialect{})
	if q != "select * from \"test\" where id = $1, email = $2" {
		t.Error("invalid postgres query")
	}
}

func TestOracleConnection(t *testing.T) {

	cfg := util.DbConfig{
		Dialect:       "oracle",
		Hostname:      "localhost",
		Username:      "KAI3",
		Password:      "KAI3",
		AdminPassword: "ChangeMe",
		DbName:        "KAI3",
		DbSid:         "XE",
		DbPort:        "1539",
	}

	connectionString, _ := cfg.GetConnectionString(true, true)
	dialect := cfg.Dialect
	conn, err := sql.Open(dialect, connectionString)

	// 	STATEMENT FOR CREATING THE DATABASE ON ORACLE
	// 	===============================================
	// 	create pluggable database KAI3
	// 	ADMIN USER kai IDENTIFIED BY KAI3
	// 	ROLES = (CONNECT, RESOURCE, SELECT_CATALOG_ROLE, PDB_DBA)
	// 	DEFAULT TABLESPACE KAI3
	// 	DATAFILE 'kai02.dbf' SIZE 250M REUSE AUTOEXTEND ON
	// 	FILE_NAME_CONVERT = ('/opt/oracle/oradata/XE/pdbseed','/opt/oracle/oradata/XE/KAI2')
	// 	PATH_PREFIX = '/opt/oracle/oradata/XE/KAI2'
	// 	TEMPFILE REUSE;

	// FOR QUERYING THE PDBSEED ORADATA BASE DIRECTORY
	// ================================================
	// 	SELECT
	//     file#,
	//     name,
	//     con_id
	// FROM
	//     v$datafile
	// WHERE
	//     con_id = (SELECT con_id FROM v$containers WHERE name = 'PDB$SEED') and lower(name) like '%system%'
	//     order by name
	//     fetch first 1 rows only;

	if err != nil {
		print("Opening failed")
	}

	defer conn.Close()

	result, err := conn.Exec("create pluggable database KAI3 ADMIN USER KAI3 IDENTIFIED BY KAI3 ROLES = (CONNECT, RESOURCE, SELECT_CATALOG_ROLE, PDB_DBA) DEFAULT TABLESPACE KAI3 DATAFILE 'kai03.dbf' SIZE 250M REUSE AUTOEXTEND ON FILE_NAME_CONVERT = ('/opt/oracle/oradata/XE/pdbseed','/opt/oracle/oradata/XE/KAI3') PATH_PREFIX = '/opt/oracle/oradata/XE/KAI3' TEMPFILE REUSE; ALTER PLUGGABLE DATABASE KAI3 OPEN;ALTER PLUGGABLE DATABASE pdb1 SAVE STATE;ALTER SESSION SET  CONTAINER = KAI3; alter user KAI3 quota unlimited on KAI3;")

	print(result)
	if err != nil {
		print("Creation of database failed")
	}
	//conn.Exec("create database " + cfg.GetDbName())
	conn.Close()

	print(connectionString)

	util.ConfigInit("/etc/semaphore.conf", false)

	store := SqlDb{}

	store.Connect("KAI3")
	db.Migrate(&store, nil)
}
