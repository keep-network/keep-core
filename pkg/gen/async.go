package gen

//go:generate sh -c "rm -f ./async/*; go run github.com/keep-network/keep-common/tools/generators/promise/ -d ./async *event.RelayEntrySubmitted *event.GroupRegistration *event.RelayEntryRequested"
