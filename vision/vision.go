package vision

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type TheImage struct {
	Content string `json:"content"`
}

type TheFeature struct {
	Type       string `json:"type"`
	MaxResults int    `json:"maxResults" bson:",omitempty"`
}

type Request struct {
	Image    TheImage     `json:"image"`
	Features []TheFeature `json:"features"`
}

type TheRequests []Request

type Body struct {
	Requests TheRequests `json:"requests"`
}

type Response struct {
	Responses []struct {
		LogoAnnotations []struct {
			Mid          string  `json:"mid"`
			Description  string  `json:"description"`
			Score        float64 `json:"score"`
			BoundingPoly struct {
				Vertices []struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"vertices"`
			} `json:"boundingPoly"`
		} `json:"logoAnnotations"`
		LabelAnnotations []struct {
			Mid         string  `json:"mid"`
			Description string  `json:"description"`
			Score       float64 `json:"score"`
			Topicality  float64 `json:"topicality"`
		} `json:"labelAnnotations"`
	} `json:"responses"`
}

func CloudVision() []string {
	var results []string
	var features []TheFeature
	var body Body

	features = append(features, TheFeature{
		Type:       "LABEL_DETECTION",
		MaxResults: 10,
	}, TheFeature{
		Type: "LOGO_DETECTION",
	})

	imgFile, err := os.Open("./image.png")
	if err != nil {
		log.Fatal(err)
	}

	defer imgFile.Close()

	fInfo, _ := imgFile.Stat()
	var size int64 = fInfo.Size()
	buffer := make([]byte, size)

	fReader := bufio.NewReader(imgFile)
	fReader.Read(buffer)

	imgBase64Str := base64.StdEncoding.EncodeToString(buffer)

	var requests TheRequests

	requests = append(requests, Request{
		Image: TheImage{
			Content: imgBase64Str,
		},
		Features: features,
	})

	body = Body{
		Requests: requests,
	}

	jsonBody, err := json.Marshal(body)

	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Post("https://vision.googleapis.com/v1/images:annotate?key="+os.Getenv("CLOUD_VISION_KEY"), "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	var obj Response

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bodyBytes, &obj)
	if err != nil {
		log.Fatal(err)
	}

	for i, resp := range obj.Responses {
		if len(resp.LogoAnnotations) > 0 {
			fmt.Println(i)
			for _, logos := range resp.LogoAnnotations {
				results = append(results, logos.Description)
			}
		}
	}

	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	return results
}
