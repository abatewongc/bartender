package service

import (
	"errors"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/coltiebaby/bastion/client"
	cu "github.com/coltiebaby/bastion/client/clientutil"
	"io"
	"math/rand"
	"strings"
	"time"
)

type BartenderService struct {
	lcu           client.Client `json:client`
	isLocked      bool          `json:isLocked`
	hasRandomized bool          `json:hasRandomized`
}

func (this *BartenderService) Listen() {
	for range time.Tick(time.Millisecond * 500) {
		if this.isChampionLocked() && !this.isLocked {
			this.isLocked = true
			// TODO: GET CHAMPION INFO FROM LEAGUE API https://ddragon.leagueoflegends.com/cdn/12.7.1/data/en_US/champion.json
			// TODO: GET SKIN INFO FROM LEAGUE API - GET SKIN INFO API INFORMATION (XD)
			// pickable-skin-ids returns a list of all of the skins you own. I haven't tested to make sure that it doesn't return banned champion skin ids.
			// Chromas are just skins, internally, as far as I know. That means we might want to consider additional logic or data collection to allow for
			// Skin and then chroma randomization. Otherwise skins with more chromas will have a bigger section of the RNG.
			err := this.selectRandomChampionSkin()
			this.isLocked = false
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// GET /lol-champ-select/v1/current-champion
func (this *BartenderService) isChampionLocked() bool {
	fmt.Print("Checking if champion is picked... ")
	url, _ := this.lcu.URL(`/lol-champ-select/v1/current-champion`)
	raw, _ := this.lcu.Get(url)

	rawBody, _ := io.ReadAll(raw.Body)
	body := string(rawBody)

	canRandomize := this.canRandomize(body)
	if !canRandomize {
		return false
	}

	return true
}

func (this *BartenderService) canRandomize(body string) bool {
	if strings.Contains(body, "404") {
		//fmt.Println("We aren't in champ select, returning false!")
		this.hasRandomized = false
		return false
	}

	if strings.EqualFold(body, "0") {
		//fmt.Println("We haven't selected a champion, returning false!")
		return false
	}

	if this.hasRandomized {
		//fmt.Println("We have already randomized, returning false")
		return false
	}

	if this.isLocked {
		//fmt.Println("We are randomizing, returning false")
		return false
	}

	return true
}

func (this *BartenderService) selectRandomChampionSkin() error {
	// ask LCU for the skin carousel
	skins := this.executeLCUGetRequest(`/lol-champ-select/v1/skin-carousel-skins`)
	selected, err := this.selectRandomChampionSkinFromList(skins)
	if err != nil {
		return err
	}

	req, err := this.getPatchRequest(selected)
	if err != nil {
		return err
	}

	// select the skin
	err = this.executeLCUPatchRequest(`/lol-champ-select/v1/session/my-selection`, req)
	if err != nil {
		return err
	}
	this.hasRandomized = true
	return nil
}

func (this *BartenderService) selectRandomChampionSkinFromList(skins string) (int, error) {
	blob, err := gabs.ParseJSON([]byte(skins))

	if err != nil {
		return -1, err
	}

	var skinInfo []SkinInfo
	skinInfo = this.extractSkins(blob, skinInfo)

	if len(skinInfo) < 1 {
		return -1, errors.New("no skins for champion! this is most definitely a bug, please contact the maintainers")
	}

	return int(skinInfo[rand.Intn(len(skinInfo))].SkinId), nil
}

func (this *BartenderService) extractSkins(blob *gabs.Container, skinInfo []SkinInfo) []SkinInfo {
	for _, child := range blob.Children() {
		ok := this.isSelectable(child)
		if !ok {
			continue
		}

		championId, ok := child.Path("championId").Data().(float64)
		if !ok {
			continue
		}

		id, ok := child.Path("id").Data().(float64)
		if !ok {
			continue
		}

		name, ok := child.Path("name").Data().(string)
		if !ok {
			continue
		}

		si := SkinInfo{SkinId: id, ChampionId: championId, SkinName: name}
		skinInfo = append(skinInfo, si)

		if child.ExistsP("childSkins") {
			skinInfo = this.extractSkins(child.Path("childSkins"), skinInfo)
		}
	}
	return skinInfo
}

func (this *BartenderService) isSelectable(child *gabs.Container) bool {
	owned, ok := child.Path("ownership.owned").Data().(bool)
	if !owned || !ok {
		return false
	}

	disabled, _ := child.Path("disabled").Data().(bool)
	if disabled {
		return false
	}

	unlocked, _ := child.Path("unlocked").Data().(bool)
	if !unlocked {
		return false
	}

	return true
}

type SkinInfo struct {
	ChampionId float64 `json:championId`
	SkinName   string  `json:skinName`
	SkinId     float64 `json:skinId`
}

func (this *BartenderService) getPatchRequest(selected int) (string, error) {
	req := gabs.New()
	req.Set(selected, "selectedSkinId")

	return req.String(), nil
}

func (this *BartenderService) executeLCUPatchRequest(endpoint string, req string) error {
	fmt.Printf("Executing PATCH request: %s with payload %s\n", endpoint, req)
	url, _ := this.lcu.URL(endpoint)

	request, err := this.lcu.NewRequest("PATCH", url, []byte(req))
	if err != nil {
		return err
	}

	_, err = cu.HttpClient.Do(request)
	if err != nil {
		return err
	}
	//fmt.Println(resp)

	return nil
}

func (this *BartenderService) executeLCUGetRequest(endpoint string) string {
	//fmt.Println("Executing GET request:" + endpoint)
	url, _ := this.lcu.URL(endpoint)
	raw, _ := this.lcu.Get(url)

	rawBody, _ := io.ReadAll(raw.Body)
	body := string(rawBody)

	return body
}

func BuildBartenderService(client client.Client) (bartenderService BartenderService) {
	bartenderService.lcu = client
	return bartenderService
}
