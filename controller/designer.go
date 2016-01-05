package controller

import (
	"encoding/json"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/webtemplate/helper"
	"io/ioutil"
	"strings"
)

type DesignerController struct {
	AppViewsPath string
}

func (t *DesignerController) getConfig(_id string) (map[string]interface{}, error) {
	bytes, err := ioutil.ReadFile(t.AppViewsPath + "data/page/page-" + _id + ".json")
	if !helper.HandleError(err) {
		return nil, err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(bytes, &data)
	if !helper.HandleError(err) {
		return nil, err
	}

	return data, nil
}
func (t *DesignerController) setConfig(_id string, config map[string]interface{}) error {
	filename := t.AppViewsPath + "data/page/page-" + _id + ".json"
	bytes, err := json.Marshal(config)
	if !helper.HandleError(err) {
		return err
	}

	err = ioutil.WriteFile(filename, bytes, 0644)
	if !helper.HandleError(err) {
		return err
	}

	return nil
}

func (t *DesignerController) GetConfig(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := map[string]string{}
	err := r.GetForms(&payload)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}

	data, err := t.getConfig(payload["_id"])
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}

	return helper.Result(true, data, "")
}

func (t *DesignerController) SetDataSource(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := map[string]string{}
	err := r.GetForms(&payload)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}
	_id := payload["_id"]

	config, err := t.getConfig(_id)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}
	config["datasources"] = strings.Split(payload["datasources"], ",")

	err = t.setConfig(_id, config)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}

	return helper.Result(true, nil, "")
}

func (t *DesignerController) GetWidgets(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := map[string]string{}
	err := r.GetForms(&payload)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}

	if payload["type"] == "chart" {
		bytes, err := ioutil.ReadFile(t.AppViewsPath + "data/chart.json")
		if !helper.HandleError(err) {
			return helper.Result(false, nil, err.Error())
		}

		data := []map[string]interface{}{}
		err = json.Unmarshal(bytes, &data)
		if !helper.HandleError(err) {
			return helper.Result(false, nil, err.Error())
		}

		return helper.Result(true, data, "")
	} else if payload["type"] == "grid" {
		bytes, err := ioutil.ReadFile(t.AppViewsPath + "data/mapgrid.json")
		if !helper.HandleError(err) {
			return helper.Result(false, nil, err.Error())
		}

		data := []map[string]interface{}{}
		err = json.Unmarshal(bytes, &data)
		if !helper.HandleError(err) {
			return helper.Result(false, nil, err.Error())
		}

		return helper.Result(true, data[0]["data"], "")
	}

	return helper.Result(true, []map[string]interface{}{}, "")
}

func (t *DesignerController) GetWidget(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := map[string]string{}
	err := r.GetForms(&payload)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}

	if payload["type"] == "chart" {
		bytes, err := ioutil.ReadFile(t.AppViewsPath + "data/chart/chart-" + payload["widgetID"] + ".json")
		if !helper.HandleError(err) {
			return helper.Result(false, nil, err.Error())
		}

		data := map[string]interface{}{}
		err = json.Unmarshal(bytes, &data)
		if !helper.HandleError(err) {
			return helper.Result(false, nil, err.Error())
		}

		return helper.Result(true, data, "")
	} else if payload["type"] == "grid" {
		connection, err := helper.LoadConfig(t.AppViewsPath + "data/grid/" + payload["widgetID"] + ".json")
		if !helper.HandleError(err) {
			return helper.Result(false, nil, err.Error())
		}
		defer connection.Close()

		cursor, err := connection.NewQuery().Select("*").Cursor(nil)
		if !helper.HandleError(err) {
			return helper.Result(false, nil, err.Error())
		}
		defer cursor.Close()

		dataSource, err := cursor.Fetch(nil, 0, false)
		if !helper.HandleError(err) {
			return helper.Result(false, nil, err.Error())
		}

		return helper.Result(true, dataSource.Data, "")
	}

	return helper.Result(true, map[string]interface{}{}, "")
}

func (t *DesignerController) AddWidget(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := map[string]string{}
	err := r.GetForms(&payload)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}
	_id := payload["_id"]

	config, err := t.getConfig(_id)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}
	content := config["content"].([]interface{})
	contentNew := map[string]interface{}{
		"dataSource": payload["dataSource"],
		"title":      payload["title"],
		"type":       payload["type"],
		"widgetID":   payload["widgetID"],
	}

	for i, eachRaw := range content {
		each := eachRaw.(map[string]interface{})
		if each["panelID"] == payload["panelID"] {
			each["content"] = append([]interface{}{contentNew}, each["content"].([]interface{})...)
		}

		config["content"].([]interface{})[i] = each
	}

	err = t.setConfig(_id, config)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}

	return helper.Result(true, nil, "")
}

func (t *DesignerController) AddPanel(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := map[string]interface{}{}
	err := r.GetForms(&payload)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}

	_id := payload["_id"].(string)
	title := payload["title"].(string)
	var width int = int(payload["width"].(float64))
	panelID := helper.RandomIDWithPrefix("p")

	config, err := t.getConfig(_id)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}
	contentOld := config["content"].([]interface{})
	contentNew := map[string]interface{}{
		"panelID": panelID,
		"title":   title,
		"width":   width,
		"content": []interface{}{},
	}
	config["content"] = append([]interface{}{contentNew}, contentOld...)

	err = t.setConfig(_id, config)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}

	return helper.Result(true, panelID, "")
}

func (t *DesignerController) RemovePanel(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := map[string]string{}
	err := r.GetForms(&payload)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}

	_id := payload["_id"]
	panelID := payload["panelID"]

	config, err := t.getConfig(_id)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}
	contentOld := config["content"].([]interface{})
	contentNew := []interface{}{}

	for _, each := range contentOld {
		if each.(map[string]interface{})["panelID"] == panelID {
			continue
		}

		contentNew = append(contentNew, each)
	}

	config["content"] = contentNew

	err = t.setConfig(_id, config)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}

	return helper.Result(true, nil, "")
}

func (t *DesignerController) SetHideShow(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := map[string]string{}
	err := r.GetForms(&payload)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}
	_id := payload["_id"]

	config, err := t.getConfig(_id)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}
	panelsid := strings.Split(payload["panelid"], ",")

	contentOld := config["content"].([]interface{})
	contentNew := []interface{}{}

	for _, each := range contentOld {
		for _, eachRaw := range panelsid {
			if eachRaw == each.(map[string]interface{})["panelID"] {
				each.(map[string]interface{})["hide"] = true
			} else {
				each.(map[string]interface{})["hide"] = false
			}
		}
		contentNew = append(contentNew, each)
	}

	config["content"] = contentNew

	err = t.setConfig(_id, config)
	if !helper.HandleError(err) {
		return helper.Result(false, nil, err.Error())
	}

	return helper.Result(true, nil, "")
}
