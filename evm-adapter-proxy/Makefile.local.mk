install-goic:
	go install github.com/aviate-labs/agent-go/cmd/goic@latest

generate-client:
	goic generate did ../canisters/logger_canister/src/logger_canister_backend/logger_canister_backend.did client --output=internal/icp/clients/logger/client.go --packageName=icpLogger
	goic generate did ../canisters/dex_canister/src/dex_canister_backend/dex_canister_backend.did client --output=internal/icp/clients/dex/client.go --packageName=icpDex