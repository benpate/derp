package derp

/*
type weirdTestTemplate struct {
	TemplateID  string `json:"templateId"  bson:"templateId"`  // Internal name/token other objects (like streams) will use to reference this Template.
	Label       string `json:"label"       bson:"label"`       // Human-readable label used in management UI.
	Description string `json:"description" bson:"description"` // Human-readable long-description text used in management UI.
	Category    string `json:"category"    bson:"category"`    // Human-readable category (grouping) used in management UI.
	IconURL     string `json:"iconUrl"     bson:"iconUrl"`     // Icon image used in management UI.
	URL         string `json:"url"         bson:"url"`         // URL where this template is published
	//	Schema      schema.Schema    `json:"schema"      bson:"schema"`      // JSON Schema that describes the data required to populate this Template.
	//	States      map[string]State `json:"states"      bson:"states"`      // Map of States (by state.ID) that Streams of this Template can be in.
	//	Views       map[string]View  `json:"views"       bson:"views"`       // Map of Views (by view.ID) that are available to Streams of this Template.
}

func TestWeirdError(t *testing.T) {

	jsonData := []byte(`{
		"label": "Article Template",
		"schema": {
			"url": "example.com/test-template",
			"title": "Test Template Schema",
			"type": "object",
			"properties": {
				"class": {
					"type": "string"
				},
				"title": {
					"type": "string",
					"description": "The human-readable title for this article"
				},
				"body": {
					"type": "string",
					"description": "The HTML content for this article"
				},
				"persons": {
					"description": "Array of people to render on the page",
					"type": "array",
					"items": {
						"type": "object",
						"properties": {
							"name": {
								"type": "string"
							},
							"email": {
								"type":"string"
							}
						}
					}
				},
				"friends": {
				"type" : "array",
				"items" : { "title" : "REFERENCE", "$ref" : "#" }
				}
			},
			"required": ["class", "title", "body", "persons"]
		},
		"states": {
			"DEFAULT": {
				"label": "Default",
				"views": ["DEFAULT", "DETAIL"]
			}
		},
		"views": {
			"DEFAULT": {
				"label": "Default",
				"file":  "default.html"
			},
			"DETAIL": {
				"label": "Detail",
				"file":  "detail.html"
			}
		}
	}

	`)

	result := Template{}

	if errr := json.Unmarshal(data, &result); errr != nil {
		spew.Dump(errr, string(data))
		errr2 := derp.Wrap(errr, "ghost.service.templateSource.File.Load", "Invalid JSON in template.json", string(data))
		derp.Report(errr2)
		return model.Template{}, errr2
	}

}
*/
