package bartender

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"

	"github.com/abatewongc/bartender-bastion/client"
	cu "github.com/abatewongc/bartender-bastion/client/clientutil"
)

var errNoSkinsFound = errors.New("no skins for champion! this is most definitely a bug, please contact the maintainers")

var endpoints = Endpoints{
	CurrentChamp: "/lol-champ-select/v1/current-champion",
	SkinCarousel: "/lol-champ-select/v1/skin-carousel-skins",
	MySelection:  "/lol-champ-select/v1/session/my-selection",
}

type SkinInfo struct {
	ChampionId float64    `json:"championId"`
	SkinName   string     `json:"skinName"`
	SkinId     float64    `json:"skinId"`
	Chromas    []SkinInfo `json:"chromas"`
}

type Endpoints struct {
	CurrentChamp string
	SkinCarousel string
	MySelection  string
}

type service struct {
	tickrate       time.Duration
	inGameTickrate time.Duration
	skinBlacklist  map[float64]struct{} // Slice of skin IDs
	endpoints      Endpoints
	lcu            client.Client
	isLocked       bool
	hasRandomized  bool
}

func New(client client.Client, options ...func(*service)) *service {
	svc := &service{
		tickrate:       time.Millisecond * 500,
		inGameTickrate: time.Minute * 8,
		skinBlacklist:  map[float64]struct{}{},
		endpoints:      endpoints,
		lcu:            client,
	}

	for _, option := range options {
		option(svc)
	}

	return svc
}

func WithTickrate(t time.Duration) func(*service) {
	return func(svc *service) {
		svc.tickrate = t
	}
}

func WithGameTickrate(t time.Duration) func(*service) {
	return func(svc *service) {
		svc.inGameTickrate = t
	}
}

func WithBlacklist(bl map[float64]struct{}) func(*service) {
	return func(svc *service) {
		svc.skinBlacklist = bl
	}
}

func (svc *service) Listen() {
	fmt.Print("Checking if champion is picked...")
	for range time.Tick(svc.tickrate) {
		if isLocked, err := svc.isChampionLocked(); isLocked && !svc.isLocked {
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
		} else if err != nil {
			fmt.Printf("\nerror encountered: %v\n", err)
		}
		fmt.Print(".")
	}
}

func (svc *service) isChampionLocked() (bool, error) {
	url, err := svc.lcu.URL(svc.endpoints.CurrentChamp)
	if err != nil {
		return false, err
	}

	resp, err := svc.lcu.Get(url)
	if err != nil {
		return false, err
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return false, errors.New("bad response from server: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return svc.canRandomize(string(body)), nil
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
	var err error

	skins, err := svc.getSkinCarousel()
	if err != nil {
		return err
	}

	selectedSkinId, err := svc.randomSkinIdFromList(skins)
	if err != nil {
		return err
	}

	err = svc.selectSkin(selectedSkinId)
	if err != nil {
		return err
	}
	svc.hasRandomized = true
	return nil
}

func (svc *service) randomSkinIdFromList(skins string) (int, error) {
	blob, err := gabs.ParseJSON([]byte(skins))

	if err != nil {
		return 0, err
	}

	var skinInfo []SkinInfo
	skinInfo = svc.extractSkins(blob, skinInfo)

	if len(skinInfo) < 1 {
		return 0, errNoSkinsFound
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

		if _, exists := svc.skinBlacklist[skin.SkinId]; !exists {
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

func (svc *service) selectSkin(skinId int) error {
	var err error

	req := gabs.New()
	req.Set(skinId, "selectedSkinId")
	url, err := svc.lcu.URL(svc.endpoints.MySelection)
	if err != nil {
		return err
	}

	fmt.Printf("\nExecuting PATCH request: %s with payload %s\n", svc.endpoints.MySelection, req.String())
	request, err := svc.lcu.NewRequest("PATCH", url, []byte(req.String()))
	if err != nil {
		return err
	}

	resp, err := cu.HttpClient.Do(request)
	if err != nil {
		return err
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New("bad response from server: " + resp.Status)
	}

	return nil
}

func (svc *service) getSkinCarousel() (string, error) {
	url, err := svc.lcu.URL(svc.endpoints.SkinCarousel)
	if err != nil {
		return "", err
	}

	resp, err := svc.lcu.Get(url)
	if err != nil {
		return "", err
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", errors.New("bad response from server: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
