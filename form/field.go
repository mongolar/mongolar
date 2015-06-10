package form

// Form fields structure
type Field struct {
	Type            string           `json:"type" bson:"type"`
	Hide            bool             `json:"hide,omitempty" bson:"hide,omitempty"`
	Key             string           `json:"key" bson:"key"`
	TemplateOptions *TemplateOptions `json:"templateOptions" bson:"templateOptions"`
	HideExpression  string           `json:"hideExpression,omitempty" bson:"hideExpression"`
	Validator       string           `json:"-" bson:"validator,omitempty"`
}

type TemplateOptions struct {
	Options     []map[string]string `json:"options,omitempty" bson:"options,omitempty"`
	Label       string              `json:"label,omitempty" bson:"label"`
	Required    bool                `json:"required,omitempty" bson:"required"`
	Placeholder string              `json:"placeholder,omitempty" bson:"placeholder,omitempty"`
	Rows        int                 `json:"rows,omitempty" bson:"rows,omitempty"`
	Cols        int                 `json:"cols,omitempty" bson:"cols,omitempty"`
	Fields      []*Field            `json:"fields,omitempty" bson:"fields,omitempty"`
	ButtonText  string              `json:"btnText,omitempty" bson:"btnText,omitempty"`
}

// Add label to field
func (f *Field) AddLabel(l string) *Field {
	f.TemplateOptions.Label = l
	return f
}

// Add
func (f *Field) AddPlaceHolder(p string) *Field {
	f.TemplateOptions.Placeholder = p
	return f
}

func (f *Field) AddRowsCols(r int, c int) *Field {
	f.TemplateOptions.Rows = r
	f.TemplateOptions.Cols = c
	return f
}

func (f *Field) AddHideExpression(he string) *Field {
	//TODO: Fix this
	f.HideExpression = he
	return f
}

func (f *Field) Required() *Field {
	f.TemplateOptions.Required = true
	return f
}

func (f *Field) Hidden() *Field {
	f.Hide = true
	return f
}
