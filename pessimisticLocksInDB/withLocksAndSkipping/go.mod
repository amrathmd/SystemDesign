module withLocksAndSkipping

go 1.23.3

replace pessimisticLocksInDB/common => ../common

require pessimisticLocksInDB/common v0.0.0-00010101000000-000000000000

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
)
