package favorites

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Favorites struct {
	List []string
	Map  map[string]int
}

func (fav *Favorites) BuildMap() {
	m := map[string]int{}
	for i, s := range fav.List {
		m[s] = i
	}
	fav.Map = m
}

func (fav *Favorites) Has(item string) bool {
	_, ok := fav.Map[item]
	return ok
}

func (fav *Favorites) Add(item string) {
	fav.Map[item] = len(fav.List)
	fav.List = append(fav.List, item)
}

func (fav *Favorites) Remove(item string) int {
	index, ok := fav.Map[item]
	if !ok {
		return -1
	}
	fav.List = append(fav.List[:index], fav.List[index+1:]...)
	fav.BuildMap()
	return index
}

func (fav *Favorites) Load(fpath string) error {
	jsonBytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	err = json.Unmarshal(jsonBytes, &fav.List)
	if err != nil {
		return err
	}
	fav.BuildMap()
	return nil
}

func (fav *Favorites) Save(fpath string) error {
	jsonBytes, err := json.Marshal(fav.List)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fpath, jsonBytes, 0o644)
	if err != nil {
		return err
	}
	return nil
}
