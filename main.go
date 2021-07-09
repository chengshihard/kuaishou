package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type data struct {
	Data struct{
		PrivateFeeds struct{
			List []struct{
				ImgUrls []string `json:"imgUrls"`
			} `json:"list"`
		} `json:"privateFeeds"`
	} `json:"data"`
}

const path = ``

func main() {
	userId := ""
	count := "500"

	request(userId,count)
}

func request(userId string,count string) {
	param := make(map[string]interface{})
	param["operationName"] = "privateFeedsQuery"
	param["query"] = "query privateFeedsQuery($principalId: String, $pcursor: String, $count: Int) {\n  privateFeeds(principalId: $principalId, pcursor: $pcursor, count: $count) {pcursor\n  list {id\n            poster\n      workType\n      imgUrls\n        }\n    __typename\n  }\n}\n"
	variables := make(map[string]string)
	variables["count"] = count
	variables["pcursor"] = ""
	variables["principalId"] = userId
	param["variables"] = variables

	payload,err :=json.Marshal(param)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client {
	}
	req, err := http.NewRequest("POST", "https://live.kuaishou.com/live_graphql", bytes.NewReader(payload))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", "clientid=3; did=web_33a9b4647194c503c2fd52b54db1a5fd; client_key=65890b29; kpn=GAME_ZONE; Hm_lvt_86a27b7db2c5c0ae37fee4a8a35033ee=1621661634; didv=1625540611000; userId=219974107; userId=219974107; kuaishou.live.bfb1s=ac5f27b3b62895859c4c1622f49856a4; kuaishou.live.web_st=ChRrdWFpc2hvdS5saXZlLndlYi5zdBKgAbEItgpd3HCuOMj08RH_yV9p-zdBPCF3j478LycA8ypb8cJ3z6sdoTcNKUq1stYUJdTObFbWJBdDgp8HnmBuAwdxPIXa5zfFtBxoO6NUNY57GF7RijEDkrEV1lsqHXIIsZViHDW7W8tjSaB-Xnt0MUqpfCJYVYkvPCRdC3tOKkJ_wTwX561udqrIcJUujHW12Edq18nGzpfcvqHboGJTookaEgL1c1j1KEeWrOe8x-vTC5n9jyIgHRG4TZo08iqz7mBk9n-HajwVO-iKiP_khaS0JSLFqycoBTAB; kuaishou.live.web_ph=bc0ae7184180afe548019eafa4f8bf74c61f; clientid=3; did=web_66b83ce831ccc99a0192ccdcc6fae5f1")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	var d data
	err = json.Unmarshal(body,&d)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, imgUrls := range d.Data.PrivateFeeds.List {
		for _, url := range imgUrls.ImgUrls {
			imageDownload(url)
		}
	}
}

func imageDownload(imageUrl string){
	req, err := http.NewRequest("GET",imageUrl,nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("user-agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)

	split := strings.Split(imageUrl,"/")
	imageName := split[len(split) - 1]


	f, err := os.Create(path + imageName)
	if err != nil {
		panic(err)
	}
	io.Copy(f, resp.Body)

	fmt.Println("图片 " + imageName + " 下载成功")
}
