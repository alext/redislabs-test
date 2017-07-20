package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Subscription struct {
	SID       int        `json:"sid"`
	PlanName  string     `json:"plan_name"`
	PlanUsage string     `json:"plan_usage"`
	Databases []Database `json:"sub_bdbs"`
}

type Database struct {
	ID     int    `json:"bdbID"`
	Name   string `json:"db_name"`
	URL    string `json:"db_master_url_public"`
	Status string `json:"db_status"`
}

func (c *Client) ListDBs() ([]Subscription, error) {

	resp, err := c.httpClient.Get(APIUrl + "/GetAvailableBDBS")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non-200 response %d", resp.StatusCode)
	}
	var data []Subscription
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) getSID() (int, error) {
	subs, err := c.ListDBs()
	if err != nil {
		return 0, err
	}
	if len(subs) < 1 {
		return 0, errors.New("No subscriptions found")
	}
	return subs[0].SID, nil
}

func (c *Client) GetDB(id int) (*Database, error) {
	subs, err := c.ListDBs()
	if err != nil {
		return nil, err
	}
	for _, d := range subs[0].Databases {
		if d.ID == id {
			return &d, nil
		}
	}
	return nil, errors.New("Database not found")
}
func (c *Client) GetDBByName(name string) (*Database, error) {
	subs, err := c.ListDBs()
	if err != nil {
		return nil, err
	}
	for _, d := range subs[0].Databases {
		if d.Name == name {
			return &d, nil
		}
	}
	return nil, errors.New("Database not found")
}

func (c *Client) ListDBsRaw() (interface{}, error) {
	resp, err := c.httpClient.Get(APIUrl + "/GetAvailableBDBS")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non-200 response %d", resp.StatusCode)
	}
	var data interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) ProvisionDB(name string) (*Database, error) {
	sid, err := c.getSID()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", APIUrl+"/CrudBdb", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("sid", strconv.Itoa(sid))
	q.Add("db_name", name)
	q.Add("repilcation", "true")
	q.Add("eviction", "volatile-lru")
	q.Add("memory", "100")
	q.Add("replication", "true")

	q.Add("security_group", "[]")
	q.Add("sip", "[]")
	q.Add("syncSources", "[]")
	q.Add("email_alerts", "[]")
	q.Add("max_cluster_mem_bdb", "500")
	q.Add("active", "true")
	req.URL.RawQuery = q.Encode()

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non-200 response %d", resp.StatusCode)
	}

	var data struct {
		Status string `json:"status"`
		DBID   int    `json:"bdbID"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	if data.Status != "success" {
		return nil, fmt.Errorf("Non-success response - %s", data.Status)
	}

	fmt.Print("Polling for database to be created")
	start := time.Now()
	var db *Database
	for {
		fmt.Print(".")
		db, err = c.GetDB(data.DBID)
		if err != nil {
			return nil, err
		}
		if db.Status == "active" {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("done in", time.Since(start))
	return db, nil
}

func (c *Client) DeleteDB(name string) error {
	sid, err := c.getSID()
	if err != nil {
		return err
	}

	db, err := c.GetDBByName(name)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", APIUrl+"/DeleteBDB", nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("sid", strconv.Itoa(sid))
	q.Add("bdbID", strconv.Itoa(db.ID))
	req.URL.RawQuery = q.Encode()

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Non-200 response %d", resp.StatusCode)
	}
	return nil
}
