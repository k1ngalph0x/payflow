module github.com/k1ngalph0x/payflow/client/wallet

go 1.25.2

require (
	github.com/k1ngalph0x/payflow/wallet-service v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.78.0
)

require (
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/k1ngalph0x/payflow/wallet-service => ../services/wallet-service
