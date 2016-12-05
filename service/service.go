package service


import (
	"github.com/oneboxtm/integrations-platform/int-avet-conciliation/elasticsearch"
)


var (
	statusChan = make(chan string)
)



func StartService() {

	elasticsearch.InitializeElasticsearch()

	// launch a goroutine for each avet club from configuration
	// for ... {
	go bookingChecker("http://10.1.9.30:10002/Ticketing/TicketingService.svc", 30, "Valencia CF", statusChan)
	// }

}




