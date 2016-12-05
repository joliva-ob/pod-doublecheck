package actuator


type HealthJson struct {
	Description string `json:"description"`
	Status string `json:"status"`
	Discovery struct {
			    Description string `json:"description"`
			    Status string `json:"status"`
			    DiscoveryClient struct {
						Description string `json:"description"`
						Status string `json:"status"`
						Services []string `json:"services"`
					} `json:"discoveryClient"`
		    } `json:"discovery"`
	DiskSpace struct {
			    Status string `json:"status"`
			    Free int64 `json:"free"`
			    Threshold int `json:"threshold"`
		    } `json:"diskSpace"`
	ConfigServer struct {
			    Status string `json:"status"`
			    PropertySources []string `json:"propertySources"`
		    } `json:"configServer"`
	Hystrix struct {
			    Status string `json:"status"`
		    } `json:"hystrix"`
}


type InfoJson struct {
	Git struct {
		    Branch string `json:"branch"`
		    Commit struct {
				   ID string `json:"id"`
				   Time string `json:"time"`
			   } `json:"commit"`
	    } `json:"git"`
}
