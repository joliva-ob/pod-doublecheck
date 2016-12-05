package service


import (
	 "time"

	"github.com/oneboxtm/integrations-platform/int-avet-conciliation/config"
	"github.com/oneboxtm/integrations-platform/int-avet-conciliation/ticketing_avet_service"
)



type bookingComparisson struct {
	ObBookingId string
	ExtBookingId string
	ObBarCode string
	ExtBarCode string
	ObStatus string
	ExtStatus string
	IsDifferent bool
}



// Goroutine to schedule and keep looking for bookings to check if Onebox databases
// by quering elasticsearch and will double-check each one with the 3rd party status
func bookingChecker( clubWsUrl string, clubCheckIntervalTimeSec int,  clubName string, statusChan chan string ) {

	ticker := time.Tick(time.Duration(clubCheckIntervalTimeSec * 1000) * time.Millisecond)
	clubName = clubName
	config.Log.Infof("BookingChecker started for club: %v", clubName)

	for {
		select {
		case <- ticker:
			processBookings(clubWsUrl, clubCheckIntervalTimeSec, clubName)
		case <- statusChan:
			config.Log.Debugf("BookingChecker for club: %v is up and running.", clubWsUrl)
		}
	}

}



func processBookings( clubWsUrl string, clubCheckIntervalTimeSec int, clubName string ) {

	config.Log.Infof("Processing bookings from club: %v...", clubName)
	bookingStatusMap := make(map[string]*bookingComparisson) // k: ONEBOX booking id v: bookingComparisson

	// Retrieve yesterday bookings from ONEBOX Elk


	bookings := make([]string, 1)
	bookings[0] = "A394AAE/17102016/100537"

	// For each ONEBOX booking from yesterday
	for i, extBookingId := range bookings {

		// Connecting to the 3rd party in order to retrieve yesterday bookings
		extBarCode, extStatus := avet_ticketing.AvetBookingStatus( clubWsUrl, extBookingId )
		bookingComp := bookingComparisson{ObBookingId: "", ExtBookingId: extBookingId, ObBarCode: "", ExtBarCode: extBarCode, ObStatus: "", ExtStatus: extStatus, IsDifferent: false}
		bookingStatusMap[extBookingId] = &bookingComp

		config.Log.Debugf( "%v -> Status for booking id: %v status is %v and barCode is %v", i, extBookingId, extStatus, extBarCode )
	}

	// Compare bookings and report results
	go compareToReport( bookingStatusMap )
}




func compareToReport( bookingStatusMap map[string]*bookingComparisson ){

	config.Log.Infof("Compare and report...")

	for i, comparisson := range bookingStatusMap {

		if comparisson.ExtStatus != comparisson.ObStatus {
			comparisson.IsDifferent = true
			config.Log.Infof("%v -> Found status difference between OB booking id: %v (%v) and Ext booking id: %v (%v)", i, comparisson.ObBookingId, comparisson.ObStatus, comparisson.ExtBookingId, comparisson.ExtStatus)
		}
	}
}