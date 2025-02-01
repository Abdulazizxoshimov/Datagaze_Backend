package service

func GetTempColor(temp float64) string {
	switch {
	case temp <= -30:
		return "#003366"
	case temp <= -20:
		return "#4A90E2"
	case temp <= -10:
		return "#B3DFFD"
	case temp <= 0:
		return "#E6F7FF"
	case temp <= 10:
		return "#D1F2D3"
	case temp <= 20:
		return "#FFFACD"
	case temp <= 30:
		return "#FFCC80"
	case temp <= 40:
		return "#FF7043"
	default:
		return "#D32F2F"
	}
}

func GetWindColor(speed float64) string {
	switch {
	case speed <= 10:
		return "#E0F7FA"
	case speed <= 20:
		return "#B2EBF2"
	case speed <= 40:
		return "#4DD0E1"
	case speed <= 60:
		return "#0288D1"
	default:
		return "#01579B"
	}
}

func GetCloudColor(cloud int) string {
	switch {
	case cloud <= 10:
		return "#FFF9C4"
	case cloud <= 30:
		return "#FFF176"
	case cloud <= 60:
		return "#E0E0E0"
	case cloud <= 90:
		return "#9E9E9E"
	default:
		return "#616161"
	}
}
