package gen

//go:generate sh -c "rm -f ./async/*; go run ./../../vendor/github.com/keep-network/keep-common/tools/generators/promise/ -d ./async *event.Entry *event.GroupTicketSubmission *event.GroupRegistration *event.Request *event.DKGResultSubmission"
