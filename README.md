# HIKVISION ISAPI

WIP! Has only 2 methods:

    ThermalCapabilites - /ISAPI/Thermal/capabilities
    ThermalJPEGWithData - /ISAPI/Thermal/channels/<id>/thermometry/jpegPicWithAppendData?format=json

## Install

```sh
go get github.com/neirolis/hikvision
```

## Usage

```go
c, err := hikvision.New(addr, user, pass)
if err != nil {...}

resp, err := c.ThermalCapabilites()
if err != nil {...}

// where resp is:
// type ThermalCapabilites struct {
// 	RealTimethermometry         bool `xml:"isSupportRealTimethermometry"`
// 	Power                       bool `xml:"isSupportPower"`
// 	RealtimeTempHumi            bool `xml:"isSupportRealtimeTempHumi"`
// 	ThermIntell                 bool `xml:"isSupportThermIntell"`
// 	ThermalPip                  bool `xml:"isSupportThermalPip"`
// 	ThermalIntelRuleDisplay     bool `xml:"isSupportThermalIntelRuleDisplay"`
// 	FaceThermometry             bool `xml:"isSupportFaceThermometry"`
// 	ThermalBlackBody            bool `xml:"isSupportThermalBlackBody"`
// 	ThermalStreamParam          bool `xml:"isSupportThermalStreamParam"`
// 	BodyTemperatureCompensation bool `xml:"isSupportBodyTemperatureCompensation"`
// 	TemperatureCorrection       bool `xml:"isSupportTemperatureCorrection"`
// 	ClickToThermometry          bool `xml:"isSupportClickToThermometry"`
// 	ThermometryHistorySearch    bool `xml:"isSupportThermometryHistorySearch"`
// 	BurningPrevention           bool `xml:"isSupportBurningPrevention"`
// 	JpegPicWithAppendData       bool `xml:"isSupportJpegPicWithAppendData"`
// 	RealTimethermometryForHTTP  bool `xml:"isSupportRealTimethermometryForHTTP"`
// 	FaceSnapThermometry         bool `xml:"isSupportFaceSnapThermometry"`
// }


data, err := c.ThermalJPEGWithData()
if err != nil {...}

// where data is:
// type ThermalData struct {
// 	Info         JpegPictureWithAppendData `json:"JpegPictureWithAppendData"`
// 	ThermalPic   []byte
// 	Temperatures []float32
// 	VisiblePic   []byte
// }

// type JpegPictureWithAppendData struct {
// 	Channel               int  `json:"channel"`
// 	JPEGPicLen            int  `json:"jpegPicLen"`
// 	JPEGPicWidth          int  `json:"jpegPicWidth"`
// 	JPEGPicHeight         int  `json:"jpegPicHeight"`
// 	P2PDataLen            int  `json:"p2pDataLen"`
// 	IsFreezedata          bool `json:"isFreezedata"`
// 	TemperatureDataLength int  `json:"temperatureDataLength"`
// 	Scale                 int  `json:"scale"`
// 	Offset                int  `json:"offset"`
// 	VisiblePicLen         int  `json:"visiblePicLen"`
// }

```
