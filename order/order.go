package order



type Order struct{
	Id int  `json:"id,omitempty"`
	Type string	`json:"type,omitempty"`
	Status int 
}

