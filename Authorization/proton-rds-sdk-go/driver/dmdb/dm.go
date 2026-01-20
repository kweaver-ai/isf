package dmdb

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"gitee.com/chunanyong/dm"
)

var replaceParam = []string{"timeout", "autocommit"}
var dmParam = []string{"connectTimeout", "autoCommit"}
var (
	errInvalidDSNMissingSymbol = errors.New("invalid DSN: missing '@' or '(' or ')' separating the necessary parts")
	errInvalidDSNNoSlash       = errors.New("invalid DSN: missing the slash separating the database name")
	errInvalidDSNUnescaped     = errors.New("invalid DSN: did you forget to escape a param value?")
	errInvalidDSNAddr          = errors.New("invalid DSN: network address not terminated (missing closing brace)")
	errNoDMSVCConf             = errors.New("invalid DMSVCConf: no dm_svc_conf,may permission problem?please check env")
)

type RDSConn struct {
	driver.Conn
}

type RDSStmt struct {
	//dm.DmStatement
	driver.Stmt
}

func (rdsStmt RDSStmt) Exec(args []driver.Value) (driver.Result, error) {
	for i, v := range args {
		if _, ok := v.([]byte); ok {
			args[i] = string(v.([]byte))
		}
	}
	if os.Getenv("RDS_SDK_DM_DEBUG") == "1" {
		fmt.Println("prepare: ", args)
	}
	return rdsStmt.Stmt.Exec(args)
}

func Open(dsn string) (driver.Conn, error) {
	dmdsn, err := NewDmdsn(dsn)
	if err != nil {
		return nil, err
	}
	dmConn, err := (&dm.DmDriver{}).Open(dmdsn)
	if err != nil {
		return nil, err
	}
	return &RDSConn{dmConn}, err
}

func (rdsConn *RDSConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	query = newDmQuery(query, args)
	for i, v := range args {
		if _, ok := v.Value.([]byte); ok {
			args[i].Value = string(v.Value.([]byte))
		}
	}
	if os.Getenv("RDS_SDK_DM_DEBUG") == "1" {
		fmt.Println(query, args)
	}
	return rdsConn.Conn.(driver.ExecerContext).ExecContext(ctx, query, args)
}

func (rdsConn *RDSConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	query = newDmQuery(query, args)
	for i, v := range args {
		if _, ok := v.Value.([]byte); ok {
			args[i].Value = string(v.Value.([]byte))
		}
	}
	if os.Getenv("RDS_SDK_DM_DEBUG") == "1" {
		fmt.Println(query, args)
	}
	return rdsConn.Conn.(driver.QueryerContext).QueryContext(ctx, query, args)
}

func (rdsConn *RDSConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	query = newDmQuery(query, nil)
	if os.Getenv("RDS_SDK_DM_DEBUG") == "1" {
		fmt.Println(query)
	}
	stmt, err := rdsConn.Conn.(driver.ConnPrepareContext).PrepareContext(ctx, query)
	return RDSStmt{stmt}, err
}

type RDSConnector struct {
	driver.Connector
}

func OpenConnector(dsn string) (driver.Connector, error) {
	dmdsn, err := NewDmdsn(dsn)
	if err != nil {
		return nil, err
	}
	dmConnector, err := (&dm.DmDriver{}).OpenConnector(dmdsn)
	return &RDSConnector{dmConnector}, err
}

func (rdsConnector *RDSConnector) Connect(ctx context.Context) (driver.Conn, error) {
	dmConn, err := rdsConnector.Connector.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return &RDSConn{dmConn}, err
}

func NewDmdsn(dsn string) (dmdsn string, err error) {
	// [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
	// Find the last '/' (since the password or the net addr might contain a '/')
	foundSlash := false
	for i := len(dsn) - 1; i >= 0; i-- {
		if dsn[i] == '/' {
			foundSlash = true
			var j, k int
			// left part is empty if i <= 0
			if i > 0 {
				// [username[:password]@][protocol[(address)]]
				// Find the last '@' in dsn[:i]
				for j = i; j >= 0; j-- {
					if dsn[j] == '@' || j == 0 {
						break
					}
				}
				// [protocol[(address)]]
				// Find the first '(' in dsn[j+1:i]
				for k = j + 1; k < i; k++ {
					if dsn[k] == '(' {
						// dsn[i-1] must be == ')' if an address is specified
						if dsn[i-1] != ')' {
							if strings.ContainsRune(dsn[k+1:i], ')') {
								return "", errInvalidDSNUnescaped
							}
							return "", errInvalidDSNAddr
						}
						break
					}
				}
				// compatible for dm dsn
				dsnpre := dsn[:j]
				for n := 0; n < len(dsnpre); n++ {
					if dsnpre[n] == ':' {
						pwd := dsnpre[n+1:]
						dsnpre = dsn[:n+1] + url.PathEscape(pwd)
					}
				}
				dsn = dsn[j:]
				for i, param := range replaceParam {
					if strings.Contains(dsn, param) {
						dsn = strings.ReplaceAll(dsn, replaceParam[i], dmParam[i])
					}
				}
				dsn, err = replaceDSN(dsn)
				if err != nil {
					return dmdsn, err
				}
				dmdsn = dsnpre + dsn
				dmdsn = strings.TrimSuffix(dmdsn, "&")
			}
			break
		}
	}

	if !foundSlash && len(dsn) > 0 {
		return "", errInvalidDSNNoSlash
	}
	dmdsn, _, err = changeDSNDbnameANDParam(dmdsn)
	if err != nil {
		return dmdsn, err
	}
	return dmdsn, nil
}

func replaceDSN(dsn string) (dmdsn string, err error) {
	var prefix, mid, suffix []string
	var invalid = false
	var host string
	if strings.Contains(dsn, "@") {
		prefix = strings.Split(dsn, "@")
	} else {
		invalid = true
	}
	if strings.Contains(dsn, "(") {
		mid = strings.Split(prefix[1], "(")
	} else {
		invalid = true
	}
	if strings.Contains(dsn, ")") {
		suffix = strings.Split(mid[1], ")")
	} else {
		invalid = true
	}
	if !invalid {
		if strings.Contains(suffix[0], ",") {
			err = customDMSVCConf(suffix[0])
			if err != nil {
				return "", errNoDMSVCConf
			}
			host = "DM"
		} else {
			host = suffix[0]
		}
		dmdsn = prefix[0] + "@" + host + suffix[1]
	} else {
		return "", errInvalidDSNMissingSymbol
	}
	return dmdsn, nil
}

func changeDSNDbnameANDParam(dsn string) (dmdsn, params string, err error) {
	length := len(dsn)
	for i := length - 1; i >= 0; i-- {
		if dsn[i] == '/' {
			var j int
			// dbname[?param1=value1&...&paramN=valueN]
			// Find the first '?' in dsn[i+1:]
			for j = i + 1; j < length; j++ {
				if dsn[j] == '?' {
					dsnpre := dsn[:j+1]
					params = dsn[j+1:]
					// change parameter values in dm dsn
					resslice := strings.Split(params, "&")
					for k, v := range resslice {
						param := strings.SplitN(v, "=", 2)
						if len(param) != 2 {
							continue
						}
						switch value := param[1]; param[0] {
						case dmParam[0]:
							t, err := time.ParseDuration(value)
							if err != nil {
								return dmdsn, params, err
							}
							resslice[k] = param[0] + "=" + strconv.FormatInt(t.Milliseconds(), 10)
						case dmParam[1]:
						default:
							resslice[k] = ""
						}
					}
					index := 0
					for _, v := range resslice {
						if v != "" {
							resslice[index] = v
							index++
						}
					}
					dsn = dsnpre + strings.Join(resslice[:index], "&")
					break
				}
			}
			dbname := dsn[i+1 : j]
			// if dbName, err = url.PathUnescape(dbname); err != nil {
			// 	return "", fmt.Errorf("invalid dbname %q: %w", dbname, err)
			// }
			dmdsn = dsn[:i] + dsn[j:]
			if j == length {
				dmdsn = "dm://" + dmdsn + "?schema=" + dbname + "&compatibleMode=mysql&escapeProcess=true&svcConfPath=/tmp/dm_svc.conf"
			} else {
				dmdsn = "dm://" + dmdsn + "&schema=" + dbname + "&compatibleMode=mysql&escapeProcess=true&svcConfPath=/tmp/dm_svc.conf"
				dmdsn = strings.Replace(dmdsn, "?&", "?", 1)
			}
			break
		}
	}
	return
}

func newDmQuery(query string, args []driver.NamedValue) (dmquery string) {
	dmquery = strings.ReplaceAll(query, "`", "\"")
	return dmquery
}

func customDMSVCConf(hostAndPort string) error {
	var dmsvc string
	if strings.Contains(hostAndPort, "]") {
		// ipv6: [ip1,ip2]:port
		ipsAndPort := strings.Split(hostAndPort, "]")
		ips := strings.TrimPrefix(ipsAndPort[0], "[")
		ip1, ip2 := strings.Split(ips, ",")[0], strings.Split(ips, ",")[1]
		port := strings.TrimPrefix(ipsAndPort[1], ":")
		dmsvc = fmt.Sprintf("DM=([%s]:%s,[%s]:%s)", ip1, port, ip2, port)
	} else {
		// ipv4: ip1,ip2:port
		ipsAndPort := strings.Split(hostAndPort, ":")
		ip1, ip2 := strings.Split(ipsAndPort[0], ",")[0], strings.Split(ipsAndPort[0], ",")[1]
		dmsvc = fmt.Sprintf("DM=(%s:%s,%s:%s)", ip1, ipsAndPort[1], ip2, ipsAndPort[1])
	}

	file, err := os.Create("/tmp/dm_svc.conf")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(dmsvc)
	if err != nil {
		return err
	}
	return nil
}
