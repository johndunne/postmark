package postmark

import (
    "io"
    "net/http"
    "bytes"
    "fmt"
    "encoding/json"
    "io/ioutil"
)

const emailUrl = "https://api.postmarkapp.com/email"
const batchUrl = "https://api.postmarkapp.com/email/batch"

type Postmark struct {
    key string
}

func NewPostmark(apikey string)(*Postmark){
    return &Postmark{ key: apikey }
}

func (p *Postmark) Send(m *Message)(* Response, error){

    data, err := m.Marshal()
    if err != nil {
        return nil, err
    }
    postData := bytes.NewBuffer(data)
    req, err := http.NewRequest("POST", emailUrl, postData)
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Postmark-Server-Token", p.key)

    rsp, err :=  http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    switch {
        case rsp.StatusCode == 401:
            return nil, fmt.Errorf("Missing of incorrect API key header")
        case rsp.StatusCode == 422:
            bytes,_:=ioutil.ReadAll(rsp.Body)
            return nil, fmt.Errorf(string(bytes))
        case rsp.StatusCode == 500:
            return nil, fmt.Errorf("Postmark seems to be down!")
    }

    var body bytes.Buffer
    _, err = io.Copy(&body, rsp.Body)
    rsp.Body.Close()
    if err != nil {
        return nil, err
    }

    prsp, err := UnmarshalResponse([]byte(body.String()))
    if err != nil {
        return nil, err
    }
    return prsp, nil
}

func (p *Postmark) SendBatch(msgs []*Message)(* Response, error){

    data, err := json.Marshal(msgs)
    if err != nil {
        return nil, err
    }
    postData := bytes.NewBuffer(data)
    req, err := http.NewRequest("POST", batchUrl, postData)
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Postmark-Server-Token", p.key)

    rsp, err :=  http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    switch {
        case rsp.StatusCode == 401:
            return nil, fmt.Errorf("Missing of incorrect API key header")
        case rsp.StatusCode == 422:
            bytes,_:=ioutil.ReadAll(rsp.Body)
            return nil, fmt.Errorf(string(bytes))
        case rsp.StatusCode == 500:
            return nil, fmt.Errorf("Postmark seems to be down!")
    }

    var body bytes.Buffer
    _, err = io.Copy(&body, rsp.Body)
    rsp.Body.Close()
    if err != nil {
        return nil, err
    }

    prsp, err := UnmarshalResponse([]byte(body.String()))
    if err != nil {
        return nil, err
    }
    return prsp, nil
}
