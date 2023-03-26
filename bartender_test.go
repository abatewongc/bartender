package bartender

import (
	"bartender/internal/mock"
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// isChampionLocked properly returns true
func TestIsChampionLocked(t *testing.T) {
	client := mock.Newclient()
	client.GetResponse = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBuffer([]byte{})),
	}
	svc := New(client)

	actual, err := svc.isChampionLocked()

	assert.NoError(t, err)
	assert.True(t, actual)
}

// isChampionLocked returns error on error fetching URL
func TestIsChampionLockedErrorURLError(t *testing.T) {
	client := mock.Newclient()
	client.URLError = errEmpty
	svc := New(client)

	_, err := svc.isChampionLocked()

	assert.Error(t, err)
}

// isChampionLocked returns error on error sending get request
func TestIsChampionLockedGetError(t *testing.T) {
	client := mock.Newclient()
	client.GetResponse = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBuffer([]byte{})),
	}
	client.GetError = errEmpty
	svc := New(client)

	_, err := svc.isChampionLocked()

	assert.Error(t, err)
}

// isChampionLocked returns error on a bad status code from the get request
func TestIsChampionLockedGetBadStatusError(t *testing.T) {
	client := mock.Newclient()
	client.GetResponse = &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(bytes.NewBuffer([]byte{})),
	}
	svc := New(client)

	_, err := svc.isChampionLocked()

	assert.Error(t, err)
}

// TestSelectRandomSkin tests selectRandomChampionSkin works
func TestSelectRandomSkin(t *testing.T) {
	var err error
	client := mock.Newclient()
	client.GetResponse = &getSkinCarouselOk
	client.DoResponse = &selectSkinOk
	svc := New(client, WithHTTPClient(client))

	err = svc.selectRandomChampionSkin()

	assert.Nil(t, err)
	assert.Equal(t, 1, client.CalledMethods["Get"])
	assert.Equal(t, 1, client.CalledMethods["Do"])
	assert.Equal(t, 1, client.CalledMethods["PATCH"])
}

// TestURLError tests selectRandomChampionSkin properly returns an error
// if the lcu client returns an error when creating a url.URL object
func TestURLError(t *testing.T) {
	//TODO
}

// TestGetSkinCarouselError tests selectRandomChampionSkin properly returns an
// error if there is an error getting the skin carousel
func TestGetSkinCarouselError(t *testing.T) {
	//TODO
}

// TestGetSkinCarouselBadResponse tests selectRandomChampionSkin properly
// returns an error on a non-200 status code when getting the skin
// carousel
func TestGetSkinCarouselBadResponse(t *testing.T) {
	//TODO
}

// TestEmptySkinList tests selectRandomChampionSkin properly returns an error
// if an empty skin list is returned from getting the skin carousel
func TestEmptySkinList(t *testing.T) {
	//TODO
}

// TestSelectSkinResponseError tests selectRandomChampionSkin properly returns
// an error if there is an error selecting the skin.
func TestSelectSkinResponseError(t *testing.T) {
	//TODO
}

// TestBadSelectSkinResponse tests selectRandomChampionSkin properly returns
// an error if a non-200 status code is returned while selecting a skin.
func TestBadSelectSkinResponse(t *testing.T) {
	//TODO
}

var errEmpty = errors.New("")

var championLockedOk = http.Response{
	StatusCode: 200,
	Body:       io.NopCloser(bytes.NewBuffer([]byte{})),
}

var championLockedNotFound = http.Response{
	StatusCode: 404,
	Body:       io.NopCloser(bytes.NewBuffer([]byte{})),
}

var championLockedBody = []byte("164")

var getSkinCarouselOk = http.Response{
	StatusCode: 200,
	Body:       io.NopCloser(bytes.NewBuffer(getSkinCarouselBody)),
	// Body:       io.NopCloser(bytes.NewBuffer([]byte{})),
}

var selectSkinOk = http.Response{
	StatusCode: 204,
}

var getSkinCarouselBody = []byte(`[{"championId":164,"childSkins":[],"chromaPreviewPath":null,"disabled":false,"emblems":[],"groupSplash":"","id":164000,"isBase":true,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Camille","ownership":{"loyaltyReward":false,"owned":true,"rental":{"rented":false},"xboxGPReward":false},"rarityGemPath":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164000.jpg","splashVideoPath":null,"stillObtainable":false,"tilePath":"/lol-game-data/assets/v1/champion-tiles/164/164000.jpg","unlocked":true},{"championId":164,"childSkins":[],"chromaPreviewPath":null,"disabled":false,"emblems":[],"groupSplash":"","id":164001,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Program Camille","ownership":{"loyaltyReward":false,"owned":true,"rental":{"rented":false},"xboxGPReward":false},"rarityGemPath":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164001.jpg","splashVideoPath":null,"stillObtainable":false,"tilePath":"/lol-game-data/assets/v1/champion-tiles/164/164001.jpg","unlocked":true},{"championId":164,"childSkins":[{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164003.png","colors":["#54209B","#54209B"],"disabled":false,"id":164003,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Coven Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164002,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164002.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164004.png","colors":["#6ABBEE","#6ABBEE"],"disabled":false,"id":164004,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Coven Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164002,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164002.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164005.png","colors":["#27211C","#27211C"],"disabled":false,"id":164005,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Coven Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164002,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164002.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164006.png","colors":["#E58BA5","#E58BA5"],"disabled":false,"id":164006,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Coven Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164002,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164002.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164007.png","colors":["#D33528","#D33528"],"disabled":false,"id":164007,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Coven Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164002,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164002.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164008.png","colors":["#ECF9F8","#ECF9F8"],"disabled":false,"id":164008,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Coven Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164002,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164002.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164009.png","colors":["#73BFBE","#73BFBE"],"disabled":false,"id":164009,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Coven Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164002,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164002.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164020.png","colors":["#162B30","#8E0A38"],"disabled":false,"id":164020,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Coven Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164002,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164002.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false}],"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164002.png","disabled":false,"emblems":[],"groupSplash":"","id":164002,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Coven Camille","ownership":{"loyaltyReward":false,"owned":true,"rental":{"rented":false},"xboxGPReward":false},"rarityGemPath":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164002.jpg","splashVideoPath":null,"stillObtainable":false,"tilePath":"/lol-game-data/assets/v1/champion-tiles/164/164002.jpg","unlocked":true},{"championId":164,"childSkins":[],"chromaPreviewPath":null,"disabled":false,"emblems":[],"groupSplash":"","id":164010,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"iG Camille","ownership":{"loyaltyReward":false,"owned":true,"rental":{"rented":false},"xboxGPReward":false},"rarityGemPath":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164010.jpg","splashVideoPath":null,"stillObtainable":false,"tilePath":"/lol-game-data/assets/v1/champion-tiles/164/164010.jpg","unlocked":true},{"championId":164,"childSkins":[{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164012.png","colors":["#D33528","#D33528"],"disabled":false,"id":164012,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Arcana Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164011,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164011.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164013.png","colors":["#DF9117","#DF9117"],"disabled":false,"id":164013,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Arcana Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164011,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164011.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164014.png","colors":["#54209B","#54209B"],"disabled":false,"id":164014,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Arcana Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164011,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164011.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164015.png","colors":["#2DA130","#2DA130"],"disabled":false,"id":164015,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Arcana Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164011,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164011.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164016.png","colors":["#2756CE","#2756CE"],"disabled":false,"id":164016,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Arcana Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164011,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164011.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164017.png","colors":["#E58BA5","#E58BA5"],"disabled":false,"id":164017,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Arcana Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164011,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164011.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164018.png","colors":["#ECF9F8","#ECF9F8"],"disabled":false,"id":164018,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Arcana Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164011,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164011.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164019.png","colors":["#27211C","#27211C"],"disabled":false,"id":164019,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Arcana Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164011,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164011.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false}],"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164011.png","disabled":false,"emblems":[],"groupSplash":"","id":164011,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Arcana Camille","ownership":{"loyaltyReward":false,"owned":true,"rental":{"rented":false},"xboxGPReward":false},"rarityGemPath":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164011.jpg","splashVideoPath":null,"stillObtainable":false,"tilePath":"/lol-game-data/assets/v1/champion-tiles/164/164011.jpg","unlocked":true},{"championId":164,"childSkins":[{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164022.png","colors":["#D33528","#D33528"],"disabled":false,"id":164022,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Strike Commander Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164021,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164021.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164023.png","colors":["#DF9117","#DF9117"],"disabled":false,"id":164023,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Strike Commander Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164021,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164021.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164024.png","colors":["#E58BA5","#E58BA5"],"disabled":false,"id":164024,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Strike Commander Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164021,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164021.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164025.png","colors":["#6ABBEE","#6ABBEE"],"disabled":false,"id":164025,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Strike Commander Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164021,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164021.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164026.png","colors":["#54209B","#54209B"],"disabled":false,"id":164026,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Strike Commander Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164021,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164021.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164027.png","colors":["#73BFBE","#73BFBE"],"disabled":false,"id":164027,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Strike Commander Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164021,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164021.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164028.png","colors":["#27211C","#27211C"],"disabled":false,"id":164028,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Strike Commander Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164021,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164021.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164029.png","colors":["#ECF9F8","#ECF9F8"],"disabled":false,"id":164029,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Strike Commander Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164021,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164021.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false},{"championId":164,"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164030.png","colors":["#FF70BA","#464646"],"disabled":false,"id":164030,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Strike Commander Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"parentSkinId":164021,"shortName":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164021.jpg","splashVideoPath":null,"stage":0,"stillObtainable":false,"tilePath":"","unlocked":false}],"chromaPreviewPath":"/lol-game-data/assets/v1/champion-chroma-images/164/164021.png","disabled":false,"emblems":[],"groupSplash":"","id":164021,"isBase":false,"isChampionUnlocked":true,"isUnlockedFromEntitledFeature":false,"name":"Strike Commander Camille","ownership":{"loyaltyReward":false,"owned":false,"rental":{"rented":false},"xboxGPReward":false},"rarityGemPath":"","splashPath":"/lol-game-data/assets/v1/champion-splashes/164/164021.jpg","splashVideoPath":null,"stillObtainable":false,"tilePath":"/lol-game-data/assets/v1/champion-tiles/164/164021.jpg","unlocked":false}]`)

var selectSkinResponse = http.Response{
	StatusCode: 200,
}
