package main

import "net/http"
import "bytes"
import "fmt"
import "io/ioutil"
import "encoding/json"
import "time"
import "strconv"

func main() {
  var url string = "http://***/api_jsonrpc.php"
  var apiuser string = "***"
  var apipassword string = "***"

  type Params struct {
    Password string `json:"password"`
    User string     `json:"user"`
  }
  parameters := Params{Password: apipassword, User: apiuser}

  type Request struct {
    Param Params    `json:"params"`
    Jsonrpc string  `json:"jsonrpc"`
    Method string   `json:"method"`
    Id int          `json:"id"`
  }

  type Request_auth struct {
    Param interface{}    `json:"params"`
    Jsonrpc string  `json:"jsonrpc"`
    Method string   `json:"method"`
    Id int          `json:"id"`
    Auth string     `json:"auth"`
  }

  type Response struct {
    Jsonrpc, Result string
    Id int
  }

  // Login
  req_login := Request{parameters, "2.0", "user.login", 0}

  req_login_json, _ := json.Marshal(req_login)
  fmt.Println(string(req_login_json))

  resp_login, err_login := http.Post(url, "application/json-rpc", bytes.NewBuffer(req_login_json))
  if err_login != nil {
  	fmt.Println("Something wired happened!!!")
  }
  defer resp_login.Body.Close()
  body_login, _ := ioutil.ReadAll(resp_login.Body)
  fmt.Println(string(body_login))
  // decode json from response []byte

  var m Response
  err2 := json.Unmarshal(body_login, &m)
  if err2 != nil {
    fmt.Println("Something wired happened!!!")
  }
  var authtoken string = m.Result
  fmt.Println(authtoken)


  type SearchParams struct {
    Output string `json:"output"`
    Filter interface{} `json:"filter"`
  }
  search_params := SearchParams{"extend", map[string]string{"host":"mgmjump01.oev-cloud.de"}}
  req_host_search := Request_auth{search_params, "2.0", "host.get", 2, authtoken}
  req_host_search_json, _ := json.Marshal(req_host_search)
  resp_host_search, err_host_search := http.Post(url, "application/json-rpc", bytes.NewBuffer(req_host_search_json))
  if err_host_search != nil {
    fmt.Println("Something wired happened!!!")
  }
  defer resp_host_search.Body.Close()
  resp_host_search_body, _ := ioutil.ReadAll(resp_host_search.Body)
  fmt.Println(string(resp_host_search_body))
  var f interface{}
  json.Unmarshal(resp_host_search_body, &f)
  /*
  fmt.Println(f)
  fm := f.(map[string]interface{})
  fmt.Println(fm["result"])
  fmm := fm["result"].([]interface{})
  fmt.Println(fmm)
  fmmm := fmm[0].(map[string]interface{})
  fmt.Println(fmmm["hostid"])
  */
  hostid := (((f.(map[string]interface{}))["result"].([]interface{}))[0].(map[string]interface{}))["hostid"]
  fmt.Println(hostid)


  fmt.Println("Create Maintenance_____________________________________")

  var duration int64 = 1800
  now := time.Now()
  fmt.Println(now)
  start := now.Unix()
  fmt.Println(start)
  end := start + duration
  fmt.Println(end)
  var buffer bytes.Buffer
  buffer.WriteString(`{"jsonrpc":"2.0","method":"maintenance.create","params":[{"groupids":[], "hostids":[`)
  buffer.WriteString(hostid.(string))
  buffer.WriteString(`], "name": "Maintenance for `)
  buffer.WriteString(hostid.(string))
  buffer.WriteString(`", "maintenance_type": "0", "description": "Scripted maintenance", "active_since": `)
  buffer.WriteString(strconv.FormatInt(start,10))
  buffer.WriteString(`, "active_till": `)
  buffer.WriteString(strconv.FormatInt(end,10))
  buffer.WriteString(`, "timeperiods": [{"timeperiod_type": 0, "start_date": `)
  buffer.WriteString(strconv.FormatInt(start,10))
  buffer.WriteString(`, "period": `)
  buffer.WriteString(strconv.FormatInt(duration,10))
  buffer.WriteString(`}]}],"auth":"`)
  buffer.WriteString(string(authtoken))
  buffer.WriteString(`","id":3}`)

  fmt.Println(string(buffer.Bytes()))

  resp_create_maintenance, err_create_maintenance := http.Post(url, "application/json-rpc", bytes.NewBuffer(buffer.Bytes()))
  if err_create_maintenance != nil {
    fmt.Println("Something wired happened!!!")
  }
  defer resp_create_maintenance.Body.Close()
  body_create_maintenance, _ := ioutil.ReadAll(resp_create_maintenance.Body)
  fmt.Println(string(body_create_maintenance))


  // Logout
  req_logout := Request_auth{parameters, "2.0", "user.logout", 1, authtoken}
  req_logout_json, _ := json.Marshal(req_logout)
  resp_logout, err_logout := http.Post(url, "application/json-rpc", bytes.NewBuffer(req_logout_json))
  if err_logout != nil {
    fmt.Println("Something wired happened!!!")
  }
  defer resp_logout.Body.Close()
  body_logout, _ := ioutil.ReadAll(resp_logout.Body)
  fmt.Println(string(body_logout))
}
