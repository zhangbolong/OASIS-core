module OASIS-core

go 1.26.1

require google.golang.org/grpc v1.79.3

require oasis-data v0.0.0-00010101000000-000000000000

require zhangbolong/OASIS-hr v0.0.0-00010101000000-000000000000

replace oasis-data => ../OASIS-data

replace zhangbolong/OASIS-hr => ../OASIS-hr

require (
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
