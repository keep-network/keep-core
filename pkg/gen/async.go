package gen

//go:generate sh -c "rm -f ./async/*; go run github.com/keep-network/keep-common/tools/generators/promise/ -d ./async *event.EntrySubmitted *event.GroupRegistration *event.Request"
