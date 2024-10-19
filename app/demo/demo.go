package demo

type DemoRequest struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     int    `json:"age"`
	Detail  string `json:"detail"`
}
