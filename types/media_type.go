package types

type MediaType int

const (
	Video MediaType = iota
	Image
	Archive
	Unknown
)

func (mt MediaType) String() string {
	switch mt {
	case Video:
		return "Video"
	case Image:
		return "Image"
	case Archive:
		return "Archive"
	case Unknown:
		return "Unknown"
	}

	return "MISSING"
}
