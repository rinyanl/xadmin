package conf

import "go.mongodb.org/mongo-driver/bson/primitive"

// xary

type Api struct {
	Services []string `json:"services"`
	Tag      string   `json:"tag"`
}

type Log struct {
	Loglevel string `json:"loglevel"`
	Access   string `json:"access"`
	Error    string `json:"error"`
}

type Dns struct {
	Servers []string `json:"servers"`
}

type Rules struct {
	Type        string   `json:"type,omitempty"`
	Ip          []string `json:"ip,omitempty"`
	OutboundTag string   `json:"outboundTag,omitempty"`
	Domain      []string `json:"domain,omitempty"`
	InboundTag  []string `json:"inboundTag,omitempty"`
}

type Routing struct {
	DomainStrategy string  `json:"domainStrategy"`
	Strategy       string  `json:"strategy"`
	Rules          []Rules `json:"rules"`
}

type Zero struct {
	StatsUserUplink   bool `json:"statsUserUplink"`
	StatsUserDownlink bool `json:"statsUserDownlink"`
}

type Levels struct {
	Zero Zero `json:"0"`
}

type Policy struct {
	Levels Levels `json:"levels"`
}

type Outbounds struct {
	Tag      string `json:"tag"`
	Protocol string `json:"protocol"`
}

type Clients struct {
	Password string `json:"password"`
	Flow     string `json:"flow"`
	Level    int    `json:"level"`
	Email    string `json:"email"`
}

type Fallbacks struct {
	Dest int `json:"dest"`
}

type Settings struct {
	Clients    []Clients   `json:"clients,omitempty"`
	Decryption string      `json:"decryption,omitempty"`
	Address    string      `json:"address,omitempty"`
	Fallbacks  []Fallbacks `json:"fallbacks,omitempty"`
}

type Certificates struct {
	CertificateFile string `json:"certificateFile"`
	KeyFile         string `json:"keyFile"`
}

type XtlsSettings struct {
	AllowInsecure bool           `json:"allowInsecure"`
	MinVersion    string         `json:"minVersion,omitempty"`
	Alpn          []string       `json:"alpn,omitempty"`
	Certificates  []Certificates `json:"certificates,omitempty"`
}

type StreamSettings struct {
	Network      string        `json:"network,omitempty"`
	Security     string        `json:"security,omitempty"`
	XtlsSettings *XtlsSettings `json:"xtlsSettings,omitempty"`
}

type Inbounds struct {
	Tag            string          `json:"tag,omitempty"`
	Port           int             `json:"port,omitempty"`
	Listen         string          `json:"listen,omitempty"`
	Protocol       string          `json:"protocol,omitempty"`
	Settings       Settings        `json:"settings,omitempty"`
	StreamSettings *StreamSettings `json:"streamSettings,omitempty"`
}

type Stats struct{}

type XaryConfig struct {
	Stats     Stats       `json:"stats"`
	Api       Api         `json:"api"`
	Log       Log         `json:"log"`
	Dns       Dns         `json:"dns"`
	Routing   Routing     `json:"routing"`
	Inbounds  []Inbounds  `json:"inbounds"`
	Policy    Policy      `json:"policy"`
	Outbounds []Outbounds `json:"outbounds"`
}

// custom

type UserTraffic struct {
	UserEmail string
	Uplink    int64
	Downlink  int64
}

type Port struct {
	Port int `json:"port,omitempty"`
}

type UseRecord struct {
	TotalFlow  int64  `bson:"totalFlow"`
	CreateTime string `bson:"createTime"`
	TimeStamp  int64  `bson:"timeStamp"`
}

type RunStatus struct {
	Xary  bool `json:"xary"`
	Nginx bool `json:"nginx"`
	Mongo bool `json:"mongo"`
}

type UserCollection struct {
	Id            primitive.ObjectID `bson:"_id" json:"_id"`
	UserEmail     string             `bson:"userEmail" json:"userEmail"`
	UserPassword  string             `bson:"userPassword" json:"userPassword"`
	Uplink        int64              `bson:"uplink" json:"uplink"`
	Downlink      int64              `bson:"downlink" json:"downlink"`
	TotalUplink   int64              `bson:"totalUplink" json:"totalUplink"`
	TotalDownlink int64              `bson:"totalDownlink" json:"totalDownlink"`
	SubDir        string             `bson:"subDir"  json:"subDir"`
	CreateTime    string             `bson:"createTime" json:"createTime"`
	TimeStamp     int64              `bson:"timeStamp" json:"timeStamp"`
}

type InboundCollection struct {
	Id        primitive.ObjectID `json:"_id"`
	Tag       string             `json:"tag"`
	Port      int                `json:"port"`
	Protocol  string             `json:"protocol"`
	UserTotal int                `json:"usertotal"`
}

// database

type Userlist struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int64       `json:"page"`
}

// sub

type Proxies struct {
	Id     primitive.ObjectID `bson:"_id" json:"_id"`
	Name   string             `bson:"name" json:"name"`
	Server string             `bson:"server" json:"server"`
	Port   int                `bson:"port" json:"port"`
}

type ProxiesJson struct {
	Id     string `bson:"_id" json:"_id"`
	Name   string `bson:"name" json:"name"`
	Server string `bson:"server" json:"server"`
	Port   int    `bson:"port" json:"port"`
}
