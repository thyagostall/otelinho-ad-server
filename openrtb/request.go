package openrtb

type BidRequest struct {
	ID     string `json:"id"`
	At     int    `json:"at"`
	Device Device `json:"device"`
	User   User   `json:"user"`
	Regs   Regs   `json:"regs"`
	Site   site   `json:"site"`
	Imp    []Imp  `json:"imp"`
}

type Device struct {
	UA         string     `json:"ua"`
	Geo        Geo        `json:"geo"`
	IP         string     `json:"ip"`
	DeviceType DeviceType `json:"devicetype"`
	Make       string     `json:"make"`
	OS         string     `json:"os"`
	Language   string     `json:"language"`
}

type DeviceType int

const (
	MobileTablet DeviceType = iota + 1
	PersonalComputer
	ConnectedTV
	Phone
	Tablet
	ConnectedDevice
	SetTopBox
)

type Geo struct {
	Lat       float32        `json:"lat"`
	Lon       float32        `json:"lon"`
	Type      SourceLocation `json:"type"`
	IPService IPService      `json:"ipservice"`
	Country   string         `json:"country"`
	Region    string         `json:"region"`
	City      string         `json:"city"`
	ZIP       string         `json:"zip"`
}

type SourceLocation int

const (
	GPSLocationServices SourceLocation = iota + 1
	IPAddress
	UserProvided
)

type IPService int

const (
	IP2Location IPService = iota + 1
	Neustar
	MaxMind
	NetAcuity
)

type User struct {
	ID   string `json:"id"`
	YOB  int    `json:"yob"`
	Data []Data `json:"data"`
}

type Data struct {
	ID      string    `json:"id"`
	Segment []Segment `json:"segment"`
}

type Segment struct {
	ID string `json:"id"`
}

type Regs struct {
	COPPA int `json:"coppa"`
}

type Imp struct {
	ID     string `json:"id"`
	TagID  string `json:"tagid"`
	Secure int    `json:"secure"`
}

type site struct {
}
