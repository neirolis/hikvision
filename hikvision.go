package hikvision

import (
	"encoding/binary"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"mime/multipart"

	"github.com/goware/urlx"
	"github.com/op/go-logging"
	dac "github.com/xinsnake/go-http-digest-auth-client"
)

var log = logging.MustGetLogger("HIKVISION-ISAPI")

type Client struct {
	addr, user, pass string
}

// New client to the hikvision device, check credentials and return new instance of Client
func New(addr, user, pass string) (*Client, error) {
	_, err := urlx.Parse(addr)
	if err != nil {
		return nil, err
	}
	if user == "" {
		return nil, errors.New("username is empty")
	}
	if pass == "" {
		return nil, errors.New("password is empty")
	}

	return &Client{
		addr: addr,
		user: user,
		pass: pass,
	}, nil
}

// URL concatinate address, path, and allow substitute variables to the path like `printf`
// no need to check error on parse address as address already checked with `New`
func (c *Client) URL(path string, a ...interface{}) string {
	u, _ := urlx.Parse(c.addr)

	u.Path = fmt.Sprintf(path, a...)

	return u.String()
}

//
//
//

type ThermalCapabilites struct {
	RealTimethermometry         bool `xml:"isSupportRealTimethermometry"`
	Power                       bool `xml:"isSupportPower"`
	RealtimeTempHumi            bool `xml:"isSupportRealtimeTempHumi"`
	ThermIntell                 bool `xml:"isSupportThermIntell"`
	ThermalPip                  bool `xml:"isSupportThermalPip"`
	ThermalIntelRuleDisplay     bool `xml:"isSupportThermalIntelRuleDisplay"`
	FaceThermometry             bool `xml:"isSupportFaceThermometry"`
	ThermalBlackBody            bool `xml:"isSupportThermalBlackBody"`
	ThermalStreamParam          bool `xml:"isSupportThermalStreamParam"`
	BodyTemperatureCompensation bool `xml:"isSupportBodyTemperatureCompensation"`
	TemperatureCorrection       bool `xml:"isSupportTemperatureCorrection"`
	ClickToThermometry          bool `xml:"isSupportClickToThermometry"`
	ThermometryHistorySearch    bool `xml:"isSupportThermometryHistorySearch"`
	BurningPrevention           bool `xml:"isSupportBurningPrevention"`
	JpegPicWithAppendData       bool `xml:"isSupportJpegPicWithAppendData"`
	RealTimethermometryForHTTP  bool `xml:"isSupportRealTimethermometryForHTTP"`
	FaceSnapThermometry         bool `xml:"isSupportFaceSnapThermometry"`
}

// ThermalCapabilites http://enpinfo.hikvision.com/unzip/20200513183429_69394_doc/GUID-376E37B7-834B-43D1-8E30-BBBCAECD07DB.html
func (c *Client) ThermalCapabilites() (data ThermalCapabilites, err error) {
	r := dac.NewRequest(c.user, c.pass, "GET", c.URL("/ISAPI/Thermal/capabilities"), "")
	resp, err := r.Execute()
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	err = xml.NewDecoder(resp.Body).Decode(&data)
	return
}

//
//
//

type ThermalData struct {
	Info         JpegPictureWithAppendData `json:"JpegPictureWithAppendData"`
	ThermalPic   []byte
	Temperatures []float32
	VisiblePic   []byte
}

type JpegPictureWithAppendData struct {
	Channel               int  `json:"channel"`
	JPEGPicLen            int  `json:"jpegPicLen"`
	JPEGPicWidth          int  `json:"jpegPicWidth"`
	JPEGPicHeight         int  `json:"jpegPicHeight"`
	P2PDataLen            int  `json:"p2pDataLen"`
	IsFreezedata          bool `json:"isFreezedata"`
	TemperatureDataLength int  `json:"temperatureDataLength"`
	Scale                 int  `json:"scale"`
	Offset                int  `json:"offset"`
	VisiblePicLen         int  `json:"visiblePicLen"`
}

// ThermalJPEGWithData http://enpinfo.hikvision.com/unzip/20200513183429_69394_doc/GUID-11C08C95-1F6C-424B-B73F-21A4BD1564D0.html#GUID-11C08C95-1F6C-424B-B73F-21A4BD1564D0
func (c *Client) ThermalJPEGWithData(channel string) (data ThermalData, err error) {
	r := dac.NewRequest(c.user, c.pass, "GET", c.URL("/ISAPI/Thermal/channels/%s/thermometry/jpegPicWithAppendData", channel)+"?format=json", "")
	resp, err := r.Execute()
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return data, errors.New(resp.Status)
	}

	reader := multipart.NewReader(resp.Body, "boundary")

	// decode json part
	jsonpart, err := reader.NextPart()
	if err != nil {
		return data, fmt.Errorf("failed to read part 'json', reason: %v", err)
	}
	if err = json.NewDecoder(jsonpart).Decode(&data); err != nil {
		return data, err
	}

	// initialize data
	data.Temperatures = make([]float32, 0, data.Info.JPEGPicWidth*data.Info.JPEGPicHeight)

	// read thermal jpeg image
	thermalPart, err := reader.NextPart()
	if err != nil {
		return data, fmt.Errorf("failed to read part 'thermalPic', reason: %v", err)
	}
	if data.ThermalPic, err = ioutil.ReadAll(thermalPart); err != nil {
		return data, fmt.Errorf("failed to read data 'thermalPic', reason: %v", err)
	}

	// read temperatures
	tempPart, err := reader.NextPart()
	if err != nil {
		return data, fmt.Errorf("failed to read part 'temperatures', reason: %v", err)
	}
	tempPoint := make([]byte, data.Info.TemperatureDataLength)
	for {
		if _, err := tempPart.Read(tempPoint); err != nil {
			break
		}

		var bits uint32
		switch data.Info.TemperatureDataLength {
		case 4:
			bits = binary.LittleEndian.Uint32(tempPoint)
		case 2:
			bits = uint32(binary.LittleEndian.Uint16(tempPoint))
		}

		data.Temperatures = append(data.Temperatures, math.Float32frombits(bits))
	}

	// read visible image
	visiblePart, err := reader.NextPart()
	if err != nil {
		return data, fmt.Errorf("failed to read part 'visiblePic', reason: %v", err)
	}
	if data.VisiblePic, err = ioutil.ReadAll(visiblePart); err != nil {
		return data, fmt.Errorf("failed to read data 'visiblePic', reason: %v", err)
	}

	return
}
