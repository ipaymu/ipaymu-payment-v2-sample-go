package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	var ipaymu_va = "1179000899"                      //your ipaymu va
	var ipaymu_key = "QbGcoO0Qds9sQFDmY0MWg1Tq.xtuh1" // your ipaymu api key

	url, _ := url.Parse("https://sandbox.ipaymu.com/api/v2/payment") //url sandbox mode
	//url, _ := url.Parse("https://my.ipaymu.com/api/v2/payment") //url production mode

	postBody, _ := json.Marshal(map[string]interface{}{
		"product":     []string{"Shoes", "Jacket"},
		"qty":         []int8{1, 2},
		"price":       []float64{350000, 200000},
		"returnUrl":   "http://your-website/thank-you-page", // your thank you page url
		"cancelUrl":   "http://your-website/cancel-page",    // your cancel page url
		"notifyUrl":   "http://your-website/callback-url",   // your callback url
		"referenceId": "TRX123",                             // reference id
		"buyerName":   "Putu",                               // optional
		"buyerEmail":  "putu@mail.com",                      // optional
		"buyerPhone":  "0812312312",                         // optional
	})

	bodyHash := sha256.Sum256([]byte(postBody))
	bodyHashToString := hex.EncodeToString(bodyHash[:])
	stringToSign := "POST:" + ipaymu_va + ":" + strings.ToLower(string(bodyHashToString)) + ":" + ipaymu_key

	h := hmac.New(sha256.New, []byte(ipaymu_key))
	h.Write([]byte(stringToSign))
	signature := hex.EncodeToString(h.Sum(nil))

	reqBody := ioutil.NopCloser(strings.NewReader(string(postBody)))

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
			"va":           {ipaymu_va},
			"signature":    {signature},
		},
		Body: reqBody,
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
}
