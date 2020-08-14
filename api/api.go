package api

// ExecV1Path defines path for ExecV1 API.
const (
	ExecV1Path   = "/api/exec/v1/"
	AttachV1Path = "/api/attach/v1/"
	EchoV1Path   = "/api/echo/v1/"
)

// PrefixBase64 prefixes base64-encoded data.
const PrefixBase64 = "base64:"

// ExecV1RequestBody defines request body for ExecV1.
type ExecV1RequestBody struct {
	Stdin string   `yaml:"Stdin"`
	Args  []string `yaml:"Args"`
}

// ExecV1ResponseBody defines response body for ExecV1.
type ExecV1ResponseBody struct {
	HTTPStatus int    `yaml:"HTTPStatus"`
	ExitStatus int    `yaml:"ExitStatus"`
	Output     string `yaml:"Output"`
	Error      string `yaml:"Error"`
}

// AttachV1Message defines message for AttachV1.
type AttachV1Message struct {
	Args []string `yaml:"Args"`
}
