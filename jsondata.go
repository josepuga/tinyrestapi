package main

import (
	"encoding/json"
	"errors"
	"sync"
)

// JSONData is a struct to manage the data used by the handlers.
//

type JSONData struct {
	items             []Item            // Unmarshall JSON content
	idCount           int               // Counter to assign new ID
	aliasesToInternal map[string]string // Aliases for JSON fields
	internalToAliases map[string]string
	mutex             sync.Mutex // Safety idCount and items modifications

}

type Item struct {
	ID           int     `json:"id"` // Must be correlative to set idCount
	Name         string  `json:"name"`
	Surname      string  `json:"surname"`
	Age          int     `json:"age"`
	Phone        string  `json:"phone"`
	CountryCode2 string  `json:"country_code_2"`
	CountryCode3 string  `json:"country_code_3"`
	CountryName  string  `json:"country_name"`
	Address      string  `json:"address"`
	ZipCode4     int     `json:"zipcode4"`
	ZipCode5     int     `json:"zipcode5"`
	City         string  `json:"city"`
	Province     string  `json:"province"`
	Email        string  `json:"email"`
	URL          string  `json:"url"`
	Check1       bool    `json:"check1"`
	Check2       bool    `json:"check2"`
	EAN          string  `json:"ean"`
	ISBN         string  `json:"isnb"`
	Price99      float64 `json:"price99"`
	Price999     float64 `json:"price999"`
	Text60       string  `json:"text60"`
	Text256      string  `json:"text256"`
	Comment      string  `json:"comment"`
}

func NewJSONData() *JSONData {
	return &JSONData{
		aliasesToInternal: make(map[string]string),
		internalToAliases: make(map[string]string),
	}
}

func (j *JSONData) AddAlias(alias, internal string) {
	j.aliasesToInternal[alias] = internal
	j.internalToAliases[internal] = alias
}

// normalizeJSONFields Transform aliases (personalized names) into internal names
func (j *JSONData) normalizeJSONFields(input []byte) ([]byte, error) {
	return j.changeJSONFields(input, j.aliasesToInternal)
}

// Not necessary...
/*
func (j *JSONData) aliasJSONFields(input []byte) ([]byte, error) {
    return j.changeJSONFields(input, j.internalToAliases)
}
*/

func (j *JSONData) changeJSONFields(input []byte, aliases map[string]string) ([]byte, error) {
	var original map[string]any
	if err := json.Unmarshal(input, &original); err != nil {
		return nil, err
	}

	// Create a  new "normalized" JSON
	normalized := make(map[string]any)
	for key, value := range original {
		if newKey, exists := aliases[key]; exists {
			normalized[newKey] = value
		} else {
			normalized[key] = value
		}
	}
	return json.Marshal(normalized)
}

// SetData sets the JSON data as []byte
func (j *JSONData) SetData(data []byte) error {
	j.items = nil
	err := json.Unmarshal(data, &j.items)
	if err != nil {
		return err
	}
	// Next Free ID.
	j.idCount = len(j.items)
	return nil
}

// GetItems Returns all JSON items
func (j *JSONData) GetItems() []map[string]any {
    j.mutex.Lock()
    defer j.mutex.Unlock()
    var itemsWithAliases []map[string]any
    for _, item := range j.items {
        itemMap, err := convertItemToMap(&item)
        if err != nil {
            continue // Ignore errors for individual items
        }
        itemsWithAliases = append(itemsWithAliases, j.aliasMap(itemMap))
    }
    return itemsWithAliases
	//return j.items
}        

// GetItemsByID Return only 1 items by ID and error
// GET
func (j *JSONData) GetItemByID(id int) (map[string]any, error) {

	for _, item := range j.items {
        if item.ID == id {
            // Convert Item to map[string]any with internal names
            itemMap, err := convertItemToMap(&item)
            if err != nil {
                return nil, err
            }

            // Change internal names to aliases
            return j.aliasMap(itemMap), nil
        }
	}
	return nil, errors.New("ID not found")
}

// Convert internal fields to alias
func (j *JSONData) aliasMap(input map[string]any) map[string]any {
	aliased := make(map[string]any)
	for key, value := range input {
		if newKey, exists := j.internalToAliases[key]; exists {
			aliased[newKey] = value
		} else {
			aliased[key] = value
		}
	}
	return aliased
}

// Insert an Item in the JSON Data. Returns map[string]any, for best handling for the server
// POST
func (j *JSONData) AddItem(rawData []byte) (map[string]any, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	// Process aliases alias->normal
	normalizedJSON, err := j.normalizeJSONFields(rawData)
	if err != nil {
		return nil, errors.New("cannot normalize JSON")
	}

	// Convert JSON in Item struct
	newItem, err := convertJSONToItem(normalizedJSON)
	if err != nil {
		return nil, err
	}

	// Set the ID
	j.idCount++
	newItem.ID = j.idCount

	// Add the new item
	j.items = append(j.items, newItem)

	// Convert Item to map[string]any with internal names
	itemMap, err := convertItemToMap(&newItem)
	if err != nil {
		return nil, err
	}

	// Change internal names to aliases
	return j.aliasMap(itemMap), nil
}

// PUT
func (j *JSONData) UpdateItem(id int, rawData []byte) (map[string]any, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	// Process aliases alias->normal
	normalizedJSON, err := j.normalizeJSONFields(rawData)
	if err != nil {
		return nil, errors.New("cannot normalize JSON")
	}

	// Convert to Item
	updatedItem, err := convertJSONToItem(normalizedJSON)
	if err != nil {
		return nil, err
	}

	for i, item := range j.items {
		if item.ID == id {
			updatedItem.ID = id // Set the same ID
			j.items[i] = updatedItem

			// Convert Item to map[string]any with internal names
			itemMap, err := convertItemToMap(&updatedItem)
			if err != nil {
				return nil, err
			}

			// Change internal names to aliases
			return j.aliasMap(itemMap), nil
		}

	}
	return nil, errors.New("ID not found")
}

// PATCH
func (j *JSONData) PatchItem(id int, rawData []byte) (map[string]any, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()


    // Process aliases alias->normal
    normalizedJSON, err := j.normalizeJSONFields(rawData)
    if err != nil {
        return nil, errors.New("cannot normalize JSON")
    }

    // Convert JSON data to dynamic map. Ie: ["id"] = 2, ["name"] = "Peter", ....
	var updates map[string]any
	if err := json.Unmarshal(normalizedJSON, &updates); err != nil {
		return nil, errors.New("Invalid JSON")
	}

	for i, item := range j.items {
		if item.ID == id {
			itemMap, err := convertItemToMap(&item)
			if err != nil {
				return nil, err
			}

			// Update only sending fields. `updates` contains that fields/values
			for k, v := range updates {
				itemMap[k] = v
			}

			// Convert to Item
			patchedItem, err := convertJSONToItem(rawData)
			if err != nil {
				return nil, err
			}

			// Replace in the list of Items
			j.items[i] = patchedItem

            // Convert to map with alias names
            return j.aliasMap(itemMap), nil
		}
	}
	return nil, errors.New("ID not found")
}

// DELETE
func (j *JSONData) DeleteItem(id int) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	for i, item := range j.items {
		if item.ID == id {
			// No good practice modify j.items inside a `for loop` of... j.items,
			// but in this case, the loop will end up just after that
			j.items = append(j.items[:i], j.items[i+1:]...) // 2nd Arg expanded!!!
			return nil
		}
	}
	return errors.New("ID not found")
}

// convertItemToMap converts an Item into a map[string]any. This value is use
// for the response of gin-gonic.
func convertItemToMap(item *Item) (map[string]any, error) {
	itemMap := make(map[string]any)

	// Convert Item struct to [] byte
	itemBytes, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	// Convert []byte to map[string]any
	if err := json.Unmarshal(itemBytes, &itemMap); err != nil {
		return nil, err
	}
	return itemMap, nil
}

func convertJSONToItem(rawData []byte) (Item, error) {
	result := Item{}
	err := json.Unmarshal(rawData, &result)
	return result, err

}
