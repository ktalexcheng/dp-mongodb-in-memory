module github.com/ktalexcheng/dp-mongodb-in-memory

go 1.19

retract v1.4.0 // Now considered unsafe as the semantics of the Start() function implicitly changed

require (
	github.com/ONSdigital/log.go/v2 v2.0.9
	github.com/smartystreets/goconvey v1.7.2
	github.com/spf13/afero v1.8.2
	go.mongodb.org/mongo-driver v1.8.0
	golang.org/x/crypto v0.0.0-20211117183948-ae814b36b871
)

require (
	github.com/fatih/color v1.13.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20210202160940-bed99a852dfe // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.3.8 // indirect
)

//Licensing error DO NOT USE
retract v1.0.0
