package clickhouse

import (
	"github.com/AMuzykus/risor/object"

	ch "github.com/ClickHouse/clickhouse-go/v2"	
)

func NewClickHouseOptions(options *object.Map) (*ch.Options, error) {

	chOptions := ch.Options{}

	if addrObj := options.Get("addr"); addrObj != nil {
		switch addrObj.Type() {
		case object.STRING: 
			if addr, errObj := object.AsString(addrObj); errObj != nil {
				return nil, errObj.Value()
			} else {
				chOptions.Addr = []string{addr}
			}
		case object.LIST:
			 if addr, errObj := object.AsList(addrObj); errObj != nil {
				return nil, errObj.Value()				
			 } else {
				for _, addrObj := range addr.Value() {
					if addr, errObj := object.AsString(addrObj); errObj != nil {
						return nil, errObj.Value()
					} else {
						chOptions.Addr = append(chOptions.Addr, addr)
					}
				} 
			 }
		default:
			return nil, object.TypeErrorf("type error: clickhouse.options.addr expected a string or list argument (got %s)", addrObj.Type())
		}
	}

	if authObj := options.Get("auth"); authObj != nil {
		if authMap, errObj := object.AsMap(authObj); errObj != nil {
			return nil, errObj.Value()
		} else {
			for key, value := range authMap.Value() {
				switch key {
				case "database":
					chOptions.Auth.Database, errObj = object.AsString(value)
					if errObj != nil {
						return nil, errObj.Value()
					}
				case "username":
					chOptions.Auth.Username, errObj = object.AsString(value)
					if errObj != nil {
						return nil, errObj.Value()
					}
				case "password":
					chOptions.Auth.Password, errObj = object.AsString(value)
					if errObj != nil {
						return nil, errObj.Value()
					}
				}
			}
		}
	}

	return &chOptions, nil
}
