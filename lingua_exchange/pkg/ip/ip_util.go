package ip

func IsIpFromChina(ip string) bool {
	return ip[:3] == "192" || ip[:3] == "10"
}
