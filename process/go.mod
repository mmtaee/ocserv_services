module process

go 1.23.5

require (
	github.com/mmtaee/go-oc-utils v0.0.11
	github.com/segmentio/kafka-go v0.4.47
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/oklog/ulid/v2 v2.1.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gorm.io/driver/postgres v1.5.11 // indirect
	gorm.io/gorm v1.25.12 // indirect
)

//replace github.com/klauspost/compress v1.17.11 => ./dev-pkg/github/compress-1.17.11
//replace github.com/pierrec/lz4/v4 v4.1.22 => ./dev-pkg/github/lz4-4.1.22
