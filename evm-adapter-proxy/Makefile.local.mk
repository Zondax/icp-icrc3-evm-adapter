install-goic:
	go install github.com/aviate-labs/agent-go/cmd/goic@latest

generate-client:
	goic generate did ../canisters/logger_canister/src/logger_canister_backend/logger_canister_backend.did client --output=internal/icp/logger_backend.go --packageName=icp