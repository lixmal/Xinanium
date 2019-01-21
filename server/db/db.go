package db

import (
    "database/sql"
    "errors"
    "log"
    "net"
    "time"

    commonconfig "../../common/config"
    "github.com/coopernurse/gorp"
    _ "github.com/lib/pq"
    "golang.org/x/crypto/bcrypt"
)

const (
    DBName     = "dev"
    DBUser     = "dev"
    DBPassword = "dev"
    UserTable  = "user"
    UserPK     = "handle"
)

// TODO: Add login game last login?
type User struct {
    commonconfig.Credentials
    Name          string    `db:"name"`
    LoginTime     time.Time `db:"login_time"`
    LoginDuration int64     `db:"login_duration"`
    Created       time.Time `db:"created"`
    Modified      time.Time `db:"modified"`
    IP            string    `db:"ip"`
    Enabled       bool      `db:"enabled"`
}

var dbmap *gorp.DbMap

func init() {
    dbh, err := sql.Open("postgres", "user="+DBUser+" dbname="+DBName+" password="+DBPassword)
    if err != nil {
        log.Fatal(err)
    }

    db := &gorp.DbMap{Db: dbh, Dialect: gorp.PostgresDialect{}}
    db.AddTableWithName(User{}, UserTable).SetKeys(false, UserPK)
    if err = db.CreateTablesIfNotExists(); err != nil {
        log.Fatal(err)
    }

    dbmap = db
    log.Println("Connection to PostgreSQL DB established")
    if err = AddUser("vik", "secret"); err != nil {
        log.Println(err)
    }
}

func GetUser(handle string, t ...*gorp.Transaction) (*User, error) {
    //err := dbmap.SelectOne(&user, "SELECT * FROM "+UserTable+" WHERE 'handle' = $1", handle)
    var usr interface{}
    var err error
    if t == nil {
        usr, err = dbmap.Get(User{}, handle)
    } else {
        usr, err = t[0].Get(User{}, handle)
    }
    if usr == nil {
        return nil, errors.New("Handle not found")
    }
    if err != nil {
        return nil, err
    }

    return usr.(*User), nil
}

func Login(creds *commonconfig.Credentials, addr net.Addr) (*User, error) {
    var handle = creds.Handle
    if len(handle) < commonconfig.HandleMinLength || len(handle) > commonconfig.HandleMaxLength {
        return nil, errors.New("Invalid handle length")
    }

    t, err := dbmap.Begin()
    if err != nil {
        return nil, err
    }
    user, err := GetUser(handle, t)
    if err != nil {
        return nil, err
    }
    if err := bcrypt.CompareHashAndPassword(user.Password, creds.Password); err != nil {
        return nil, err
    }
    // Enabled check after password check, so attacker cannot get any information
    if !user.Enabled {
        return nil, errors.New("Account is disabled")
    }
    user.LoginTime = time.Now()
    if ip, _, err := net.SplitHostPort(addr.String()); err == nil {
        user.IP = ip
    }
    if _, err := t.Update(user); err != nil {
        return nil, err
    }
    if err := t.Commit(); err != nil {
        return nil, err
    }
    return user, nil
}

func AddUser(handle, password string) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), commonconfig.HashCost)
    if err != nil {
        return err
    }

    // start transaction to make sure nobody creates a user
    // between get and insert
    t, err := dbmap.Begin()
    if err != nil {
        return err
    }
    if _, err = GetUser(handle, t); err == nil {
        return errors.New("User exists or connection error")
    }
    if err = t.Insert(
        &User{
            Credentials: commonconfig.Credentials{
                Password: hash,
                Handle:   handle,
            },
            Name:    handle,
            // Created: time.Now(),
            // TODO: check if init values are inserted (like a date 1970)
        },
    ); err != nil {
        return err
    }
    if err := t.Commit(); err != nil {
        return err
    }

    return nil
}
