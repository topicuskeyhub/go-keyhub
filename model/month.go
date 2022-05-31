package model

const (
	MONTH_NAME_JAN = "JANUARY"
	MONTH_NAME_FEB = "FEBRUARY"
	MONTH_NAME_MAR = "MARCH"
	MONTH_NAME_APR = "APRIL"
	MONTH_NAME_MAY = "MAY"
	MONTH_NAME_JUN = "JUNE"
	MONTH_NAME_JUL = "JULY"
	MONTH_NAME_AUG = "AUGUST"
	MONTH_NAME_SEP = "SEPTEMBER"
	MONTH_NAME_OCT = "OCTOBER"
	MONTH_NAME_NOV = "NOVEMBER"
	MONTH_NAME_DEC = "DECEMBER"
)

type MonthSelection struct {
	Jan bool
	Feb bool
	Mar bool
	Apr bool
	May bool
	Jun bool
	Jul bool
	Aug bool
	Sep bool
	Oct bool
	Nov bool
	Dec bool
}

func (ms *MonthSelection) Enable(month ...string) {
	for _, m := range month {
		ms.Set(m, true)
	}
}

func (ms *MonthSelection) Disable(month ...string) {
	for _, m := range month {
		ms.Set(m, false)
	}
}
func (ms *MonthSelection) Set(month string, value bool) {
	switch month {
	case MONTH_NAME_JAN:
		ms.Jan = value
	case MONTH_NAME_FEB:
		ms.Feb = value
	case MONTH_NAME_MAR:
		ms.Mar = value
	case MONTH_NAME_APR:
		ms.Apr = value
	case MONTH_NAME_MAY:
		ms.May = value
	case MONTH_NAME_JUN:
		ms.Jun = value
	case MONTH_NAME_JUL:
		ms.Jul = value
	case MONTH_NAME_AUG:
		ms.Aug = value
	case MONTH_NAME_SEP:
		ms.Sep = value
	case MONTH_NAME_OCT:
		ms.Oct = value
	case MONTH_NAME_NOV:
		ms.Nov = value
	case MONTH_NAME_DEC:
		ms.Dec = value
	}
}

func (ms *MonthSelection) ToList() []string {

	months := []string{}

	if ms.Jan {
		months = append(months, MONTH_NAME_JAN)
	}
	if ms.Feb {
		months = append(months, MONTH_NAME_FEB)
	}
	if ms.Mar {
		months = append(months, MONTH_NAME_MAR)
	}
	if ms.Apr {
		months = append(months, MONTH_NAME_APR)
	}
	if ms.May {
		months = append(months, MONTH_NAME_MAY)
	}
	if ms.Jun {
		months = append(months, MONTH_NAME_JUN)
	}
	if ms.Jul {
		months = append(months, MONTH_NAME_JUL)
	}
	if ms.Aug {
		months = append(months, MONTH_NAME_AUG)
	}
	if ms.Sep {
		months = append(months, MONTH_NAME_SEP)
	}
	if ms.Oct {
		months = append(months, MONTH_NAME_OCT)
	}
	if ms.Nov {
		months = append(months, MONTH_NAME_NOV)
	}
	if ms.Dec {
		months = append(months, MONTH_NAME_DEC)
	}

	return months

}
