package common

import (
	"fmt"
	"math/rand"
	"net"
	"time"
	"os"
)

// 使用真实有效的美国ip
// https://lite.ip2location.com/united-states-of-america-ip-address-ranges
// https://cdn-lite.ip2location.com/datasets/US.json?_=1683336720620
//
//	async function getIpRange() {
//	  const results = await fetch(`https://cdn-lite.ip2location.com/datasets/US.json?_=${Date.now()}`)
//	    .then((res) => res.json())
//	    .then((res) => {
//	      const limitCount = 10000;
//	      return res.data.filter((x) => parseInt(x[2].replace(/,/g,'')) >= limitCount).map((x) => `{"${x[0]}", "${x[1]}"}, //${x[2]}`);
//	    });
//	    console.log(`results : `,results);
//	    return results.join('\n');
//	}
//
// copy(await getIpRange());
var IP_RANGE = [][]string{
	{"4.150.64.0", "4.150.127.255"},      // Azure Cloud EastUS2 16382
	{"4.152.0.0", "4.153.255.255"},       // Azure Cloud EastUS2 131070
	{"13.68.0.0", "13.68.127.255"},       // Azure Cloud EastUS2 32766
	{"13.104.216.0", "13.104.216.255"},   // Azure EastUS2 256
	{"20.1.128.0", "20.1.255.255"},       // Azure Cloud EastUS2 32766
	{"20.7.0.0", "20.7.255.255"},         // Azure Cloud EastUS2 65534
	{"20.22.0.0", "20.22.255.255"},       // Azure Cloud EastUS2 65534
	{"40.84.0.0", "40.84.127.255"},       // Azure Cloud EastUS2 32766
	{"40.123.0.0", "40.123.127.255"},     // Azure Cloud EastUS2 32766
	{"4.214.0.0", "4.215.255.255"},       // Azure Cloud JapanEast 131070
	{"4.241.0.0", "4.241.255.255"},       // Azure Cloud JapanEast 65534
	{"40.115.128.0", "40.115.255.255"},   // Azure Cloud JapanEast 32766
	{"52.140.192.0", "52.140.255.255"},   // Azure Cloud JapanEast 16382
	{"104.41.160.0", "104.41.191.255"},   // Azure Cloud JapanEast 8190
	{"138.91.0.0", "138.91.15.255"},      // Azure Cloud JapanEast 4094
	{"151.206.65.0", "151.206.79.255"},   // Azure Cloud JapanEast 256
	{"191.237.240.0", "191.237.241.255"}, // Azure Cloud JapanEast 512
	{"4.208.0.0", "4.209.255.255"},       // Azure Cloud NorthEurope 131070
	{"52.169.0.0", "52.169.255.255"},     // Azure Cloud NorthEurope 65534
	{"68.219.0.0", "68.219.127.255"},     // Azure Cloud NorthEurope 32766
	{"65.52.64.0", "65.52.79.255"},       // Azure Cloud NorthEurope 4094
	{"98.71.0.0", "98.71.127.255"},       // Azure Cloud NorthEurope 32766
	{"74.234.0.0", "74.234.127.255"},     // Azure Cloud NorthEurope 32766
	{"4.151.0.0", "4.151.255.255"},       // Azure Cloud SouthCentralUS 65534
	{"13.84.0.0", "13.85.255.255"},       // Azure Cloud SouthCentralUS 131070
	{"4.255.128.0", "4.255.255.255"},     // Azure Cloud WestCentralUS 32766
	{"13.78.128.0", "13.78.255.255"},     // Azure Cloud WestCentralUS 32766
	{"4.175.0.0", "4.175.255.255"},       // Azure Cloud WestEurope 65534
	{"13.80.0.0", "13.81.255.255"},       // Azure Cloud WestEurope 131070
	{"20.73.0.0", "20.73.255.255"},       // Azure Cloud WestEurope 65534
}

var MUID_ADDRESSES = []string{
	"074AD7F106536BC6392FC4C907CA6AEA",
	"019546D2D9086B1C238C555FD84B6A2A",
	"226B78B3878768AE2C6A6B3E864069FA",
	"1BD4F74902356D8C047AE4C403F26C19",
	"2BA6A324A2FC66D834B4B0A9A34F6722",
	"22897C804CF66A462F436F0D4DB56B13",
	"3B06D20557436C4D1D98C18856F06DC9",
	"0D0604C3DD7469723A9B174EDC3768CD",
	"3FE6DC08443C6B280FFBCF8545FB6AFA",
	//
	"06FCC9A5B8C26DD20FBEDA3CB9C56CFB",
	"38436E625858684411327DFB595F69B2",
	"0DF624B60F646DE8128F372F0ED36C4A",
	// 添加更多的IP地址
	"302E487C4D4E664139B85BB44C3B6714",
	"35FB747F5FD86DA22ADF67B45E976C79",
	"03DFB40F16DC6190345BA7C4176B604A",
	"07E349B2778C630D10165A79763B62F6",
	"065BFA8B7F1A62C3119DE9407EAD63DB",
	"1B949776444761C02FDE84BD450860D5",
	"0FAD13FA5D066034007600315C1461A8",
	"3B9B864207FB65702CC9958906E964A0",
	"33CE6C22A9D567A13FC87FE9A8C766B8",
	"154EEEC90A1769CF2D2CFD020B056830",
}


func generateRandomString(length int) string {
	charset := "ABCDEF1234567890"
	rand.Seed(time.Now().UnixNano())

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}



// 获取真实有效的随机IP
func GetRandomIP() string {
	xfip := os.Getenv("X_For_IP")


	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))

	//if USER_MUID == "" {
		IPSTR := MUID_ADDRESSES[rng.Intn(len(MUID_ADDRESSES))]
		trimmedIPStr := IPSTR[:len(IPSTR)-2]
		randomString := generateRandomString(2)
		USER_MUID = trimmedIPStr + randomString
	//	}


	if xfip == "" {
	
	// 生成随机索引
	randomIndex := rng.Intn(len(IP_RANGE))

	// 获取随机 IP 地址范围
	startIP := IP_RANGE[randomIndex][0]
	endIP := IP_RANGE[randomIndex][1]

	// 将起始 IP 地址转换为整数形式
	startIPInt := ipToUint32(net.ParseIP(startIP))
	// 将结束 IP 地址转换为整数形式
	endIPInt := ipToUint32(net.ParseIP(endIP))

	// 生成随机 IP 地址
	randomIPInt := rng.Uint32()%(endIPInt-startIPInt+1) + startIPInt
	randomIP := uint32ToIP(randomIPInt)

	return randomIP
	}else {
		return xfip
	}
}

// 将 IP 地址转换为 uint32
func ipToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	var result uint32
	result += uint32(ip[0]) << 24
	result += uint32(ip[1]) << 16
	result += uint32(ip[2]) << 8
	result += uint32(ip[3])
	return result
}

// 将 uint32 转换为 IP 地址
func uint32ToIP(intIP uint32) string {
	ip := fmt.Sprintf("%d.%d.%d.%d", byte(intIP>>24), byte(intIP>>16), byte(intIP>>8), byte(intIP))
	return ip
}
