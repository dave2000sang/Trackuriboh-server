package models

// Response JSON object from https://docs.tcgplayer.com/reference/catalog_getcategorygroups-1
type CategoryResponse struct {
	TotalItems int            `json:"totalItems"`
	Success    bool           `json:"success"`
	Errors     []string       `json:"errors"`
	Results    []CategoryData `json:"results"`
}

type CategoryData struct {
	GroupId        int    `json:"groupId"`
	Name           string `json:"name"`
	Abbreviation   string `json:"abbreviation"`
	IsSupplemental bool   `json:"isSupplemental"`
	PublishedOn    string `json:"publishedOn"`
	ModifiedOn     string `json:"modifiedOn"`
	CategoryId     int    `json:"categoryId"`
}

type ProductData struct {
	ProductId  int    `json:"productId"`
	Name       string `json:"name"`
	CleanName  string `json:"cleanName"`
	ImageUrl   string `json:"imageUrl"`
	CategoryId int    `json:"categoryId"`
	GroupId    int    `json:"groupId"`
	Url        string `json:"url"`
	ModifiedOn string `json:"modifiedOn"`
}

type SKUData struct {
	SkuId       int `json:"skuId"`
	ProductId   int `json:"productId"`
	LanguageId  int `json:"languageId"`
	PrintingId  int `json:"printingId"`
	ConditionId int `json:"conditionId"`
}

type ConditionData struct {
	ConditionId  int `json:"conditionId"`
	Name         int `json:"name"`
	Abbreviation int `json:"abbreviation"`
	DisplayOrder int `json:"displayOrder"`
}
