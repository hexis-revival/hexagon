module github.com/hexis-revival/hexagon/hscore

go 1.22.7

replace github.com/hexis-revival/hexagon/common => ../common

require github.com/hexis-revival/hexagon/common v0.0.0-20241005191625-f1f1499b1fd1

require github.com/lib/pq v1.10.9 // indirect

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gorilla/mux v1.8.1
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/redis/go-redis/v9 v9.6.1 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/exp v0.0.0-20241004190924-225e2abe05e6 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gorm.io/driver/postgres v1.5.9 // indirect
	gorm.io/gorm v1.25.12 // indirect
)
