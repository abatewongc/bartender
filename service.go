package bartender

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/coltiebaby/bastion/client"
	cu "github.com/coltiebaby/bastion/client/clientutil"
)

type SkinInfo struct {
	ChampionId float64    `json:"championId"`
	SkinName   string     `json:"skinName"`
	SkinId     float64    `json:"skinId"`
	Chromas    []SkinInfo `json:"chromas"`
}

type Config struct {
	SkinBlacklist  map[float64]struct{} // Slice of skin IDs
	Tickrate       time.Duration
	InGameTickrate time.Duration
}

type service struct {
	cnf           Config
	lcu           client.Client
	isLocked      bool
	hasRandomized bool
}

func New(cnf Config, client client.Client) *service {
	return &service{
		lcu: client,
	}
}

func (svc *service) Listen() {
	fmt.Print("Checking if champion is picked...")
	for range time.Tick(svc.cnf.Tickrate) {
		if svc.isChampionLocked() && !svc.isLocked {
			svc.isLocked = true
			// TODO: GET CHAMPION INFO FROM LEAGUE API https://ddragon.leagueoflegends.com/cdn/12.7.1/data/en_US/champion.json
			// TODO: GET SKIN INFO FROM LEAGUE API - GET SKIN INFO API INFORMATION (XD)
			// pickable-skin-ids returns a list of all of the skins you own. I haven't tested to make sure that it doesn't return banned champion skin ids.
			// Chromas are just skins, internally, as far as I know. That means we might want to consider additional logic or data collection to allow for
			// Skin and then chroma randomization. Otherwise skins with more chromas will have a bigger section of the RNG.
			err := svc.selectRandomChampionSkin()
			svc.isLocked = false
			if err != nil {
				fmt.Println(err)
			}
		}
		fmt.Print(".")
	}
}

// GET /lol-champ-select/v1/current-champion
func (svc *service) isChampionLocked() bool {
	url, _ := svc.lcu.URL(`/lol-champ-select/v1/current-champion`)
	raw, _ := svc.lcu.Get(url)

	rawBody, _ := io.ReadAll(raw.Body)
	body := string(rawBody)

	return svc.canRandomize(body)
}

func (svc *service) canRandomize(body string) bool {
	switch true {
	case strings.Contains(body, "404"):
		svc.hasRandomized = false
		return false
	case strings.EqualFold(body, "0"):
		return false
	case svc.hasRandomized:
		return false
	case svc.isLocked:
		return false
	default:
		return true
	}
}

func (svc *service) selectRandomChampionSkin() error {
	// ask LCU for the skin carousel
	skins := svc.executeLCUGetRequest(`/lol-champ-select/v1/skin-carousel-skins`)

	selected, err := svc.selectRandomChampionSkinFromList(skins)
	if err != nil {
		return err
	}

	req, err := svc.getPatchRequest(selected)
	if err != nil {
		return err
	}

	// select the skin
	err = svc.executeLCUPatchRequest(`/lol-champ-select/v1/session/my-selection`, req)
	if err != nil {
		return err
	}
	svc.hasRandomized = true
	return nil
}

func (svc *service) selectRandomChampionSkinFromList(skins string) (int, error) {
	blob, err := gabs.ParseJSON([]byte(skins))

	if err != nil {
		return -1, err
	}

	var skinInfo []SkinInfo
	skinInfo = svc.extractSkins(blob, skinInfo)

	if len(skinInfo) < 1 {
		return -1, errors.New("no skins for champion! this is most definitely a bug, please contact the maintainers")
	}

	// Reroll until a skin not in the blacklist is rolled
	var skin SkinInfo

	for {
		skin = skinInfo[rand.Intn(len(skinInfo))]

		// Reroll random chroma
		if len(skin.Chromas) > 0 {
			var skinAndChromas []SkinInfo = append(skin.Chromas, skin)
			skin = skinAndChromas[rand.Intn(len(skinAndChromas))]
		}

		if _, exists := svc.cnf.SkinBlacklist[skin.SkinId]; !exists {
			break
		}
	}

	return int(skin.SkinId), nil
}

func (svc *service) extractSkins(blob *gabs.Container, skinInfos []SkinInfo) []SkinInfo {
	for _, child := range blob.Children() {
		ok := svc.isSelectable(child)
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

		var chromas []SkinInfo
		if child.ExistsP("childSkins") {
			chromas = svc.extractSkins(child.Path("childSkins"), chromas)
		}

		si := SkinInfo{SkinId: id, ChampionId: championId, SkinName: name, Chromas: chromas}
		skinInfos = append(skinInfos, si)
	}
	return skinInfos
}

func (svc *service) isSelectable(child *gabs.Container) bool {
	owned, ok := child.Path("ownership.owned").Data().(bool)
	if !owned || !ok {
		return false
	}

	disabled, _ := child.Path("disabled").Data().(bool)
	if disabled {
		return false
	}

	unlocked, _ := child.Path("unlocked").Data().(bool)
	return unlocked
}

func (svc *service) getPatchRequest(selected int) (string, error) {
	req := gabs.New()
	req.Set(selected, "selectedSkinId")

	return req.String(), nil
}

func (svc *service) executeLCUPatchRequest(endpoint string, req string) error {
	fmt.Printf("Executing PATCH request: %s with payload %s\n", endpoint, req)
	url, _ := svc.lcu.URL(endpoint)

	request, err := svc.lcu.NewRequest("PATCH", url, []byte(req))
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

func (svc *service) executeLCUGetRequest(endpoint string) string {
	//fmt.Println("Executing GET request:" + endpoint)
	url, _ := svc.lcu.URL(endpoint)
	raw, _ := svc.lcu.Get(url)

	rawBody, _ := io.ReadAll(raw.Body)
	body := string(rawBody)

	return body
}
