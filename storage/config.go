package storage

import (
	"github.com/codingXiang/go-orm"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

const (
	Primary            = "Primary"            // Default mode. All operations read from the current replica set primary.
	PrimaryPreferred   = "PrimaryPreferred"   // Read from the primary if available. Read from the secondary otherwise.
	Secondary          = "Secondary"          // Read from one of the nearest secondary members of the replica set.
	SecondaryPreferred = "SecondaryPreferred" // Read from one of the nearest secondaries if available. Read from primary otherwise.
	Nearest            = "Nearest"            // Read from one of the nearest members, irrespective of it being primary or secondary.
	Eventual           = "Eventual"           // Same as Nearest, but may change servers between reads.
	Monotonic          = "Monotonic"          // Same as SecondaryPreferred before first write. Same as Primary after first write.
	Strong             = "Strong"             // Same as Primary.
)

type ConfigSource struct {
	config *viper.Viper
	mongo  *mongo
	orm    *gorm.DB
}

type mongo struct {
	session *mgo.Session
	db      *mgo.Database
}

func NewConfigSource(config *viper.Viper, o orm.OrmInterface) *ConfigSource {
	s := new(ConfigSource).
		setConfig(config).
		initMongo().
		initOrm(o)
	return s
}

func (s *ConfigSource) setConfig(config *viper.Viper) *ConfigSource {
	s.config = config
	return s
}

func (s *ConfigSource) initOrm(o orm.OrmInterface) *ConfigSource {
	s.orm = o.GetInstance()
	return s
}

func (s *ConfigSource) initMongo() *ConfigSource {
	var err error
	s.mongo = new(mongo)
	s.mongo.session, err = mgo.Dial(s.config.GetString("mongo.url") + ":" + s.config.GetString("mongo.port"))
	cred := s.getCredential(s.config.GetString("mongo.username"), s.config.GetString("mongo.password"))
	err = s.sessionLogin(cred)
	if err != nil {
		panic(err)
	}
	sessionMode := s.config.GetString("mongo.mode")
	s.mongo.session.SetMode(s.getSessionMode(sessionMode), true)
	if err != nil {
		panic(err)
	}
	s.mongo.db = s.mongo.session.DB(s.config.GetString("mongo.database"))
	return s
}

//sessionLogin mongo 登入
func (c *ConfigSource) sessionLogin(cred *mgo.Credential) error {
	log.Println("login")
	if cred != nil {
		return c.mongo.session.Login(cred)
	} else {
		log.Println("cred is null, not login")
		return nil
	}
}

//getCredential 取得 mongo 連線憑證
func (c *ConfigSource) getCredential(username, password string) *mgo.Credential {
	if username == "" || password == "" {
		return nil
	}
	return &mgo.Credential{
		Username: username,
		Password: password,
	}
}

type LoginInfo struct {
	UserID     int64     `json:"userId"`
	ClientIP   string    `json:"clientIP"`
	LoginState string    `json:"loginState"`
	LoginTime  time.Time `json:"loginTime"`
}

//getSessionMode 取得 mongo 的 session 模式
func (c *ConfigSource) getSessionMode(mode string) mgo.Mode {
	switch mode {
	case Primary:
		return mgo.Primary
	case PrimaryPreferred:
		return mgo.PrimaryPreferred
	case Secondary:
		return mgo.Secondary
	case SecondaryPreferred:
		return mgo.SecondaryPreferred
	case Nearest:
		return mgo.Nearest
	case Eventual:
		return mgo.Eventual
	case Monotonic:
		return mgo.Monotonic
	case Strong:
		return mgo.Strong
	default:
		return mgo.Primary
	}
}
