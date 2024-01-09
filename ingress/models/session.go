package models

type Session struct {
	UserAgent  string  `bson:"user_agent" json:"ua"`
	IP         string  `bson:"ip"`
	TimeOrigin float64 `bson:"time_origin" json:"to"`
	DeviceType int     `bson:"device_type" json:"dt"`
	Referrer   string  `bson:"referrer" json:"rf"`
	Domain     string  `bson:"domain" json:"dm"`
}
