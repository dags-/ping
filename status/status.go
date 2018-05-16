package status

type Status struct {
	Type  string `json:"type"`
	Error error  `json:"error,omitempty"`
	Data  *Data  `json:"data,omitempty"`
}

type Data struct {
	Version     Version     `json:"version"`
	Description Description `json:"description"`
	Players     Players     `json:"players"`
	ModInfo     *ModInfo    `json:"modinfo"`
	Favicon     string      `json:"favicon"`
}

type Description struct {
	Text string `json:"text"`
}

type Players struct {
	Max    int      `json:"max"`
	Online int      `json:"online"`
	Sample []Player `json:"sample"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type ModInfo struct {
	Type    string `json:"type"`
	ModList []Mod  `json:"modList"`
}

type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Mod struct {
	ModID   string `json:"modid"`
	Version string `json:"version"`
}
