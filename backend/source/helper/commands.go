package helper

import "fmt"

// primitiveMap 所有支持的命令映射
var primitiveMap = map[string]string{
	"summary":                          "show protocols",
	"detail":                           "show protocols all %s",
	"route_from_protocol":              "show route protocol %s",
	"route_from_protocol_all":          "show route protocol %s all",
	"route_from_protocol_primary":      "show route protocol %s primary",
	"route_from_protocol_all_primary":  "show route protocol %s all primary",
	"route_filtered_from_protocol":     "show route filtered protocol %s",
	"route_filtered_from_protocol_all": "show route filtered protocol %s all",
	"route_from_origin":                "show route where bgp_path.last = %s",
	"route_from_origin_all":            "show route where bgp_path.last = %s all",
	"route_from_origin_primary":        "show route where bgp_path.last = %s primary",
	"route_from_origin_all_primary":    "show route where bgp_path.last = %s all primary",
	"route":                            "show route for %s",
	"route_all":                        "show route for %s all",
	"route_where":                      "show route where net ~ [ %s ]",
	"route_where_all":                  "show route where net ~ [ %s ] all",
	"route_generic":                    "show route %s",
	"generic":                          "show %s",
}

// GetBirdCommand 根据命令类型和参数获取对应的Bird命令
func GetBirdCommand(cmdType string, param string) string {
	if cmd, ok := primitiveMap[cmdType]; ok {
		if param != "" {
			return fmt.Sprintf(cmd, param)
		}
		return cmd
	}
	return ""
}
