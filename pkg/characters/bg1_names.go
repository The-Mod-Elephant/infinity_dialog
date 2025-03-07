package characters

const (
	BAJANT BG1_Characters = iota
	BALORA
	BBRANW
	BCORAN
	BDORN
	BDYNAH
	BEDWIN
	BELDOT
	BFALDO
	BGARRI
	BIMOEN
	BJAHEI
	BKAGAI
	BKHALI
	BKIVAN
	BMINSC
	BMONTA
	BNEERA
	BQUAYL
	BRASAAD
	BSAFAN
	BSHART
	BSKIE
	BTIAX
	BVICON
	BXANNN
	BXZAR
	BYESLI
)

type BG1_Characters int64

func (s BG1_Characters) String() string {
	switch s {
	case BKIVAN:
		return "KIVAN"
	case BALORA:
		return "ALORA"
	case BMINSC:
		return "MINSC"
	case BDYNAH:
		return "DYNAHEIR"
	case BYESLI:
		return "YESLICK"
	case BCORAN:
		return "CORAN"
	case BAJANT:
		return "AJANTIS"
	case BKHALI:
		return "KHALID"
	case BJAHEI:
		return "JAHEIRA"
	case BGARRI:
		return "GARRICK"
	case BSAFAN:
		return "SAFANA"
	case BFALDO:
		return "FALDORN"
	case BBRANW:
		return "BRANWEN"
	case BQUAYL:
		return "QUAYLE"
	case BXANNN:
		return "XAN"
	case BSKIE:
		return "SKIE"
	case BELDOT:
		return "ELDOTH"
	case BXZAR:
		return "XZAR"
	case BMONTA:
		return "MONTARON"
	case BTIAX:
		return "TIAX"
	case BKAGAI:
		return "KAGAIN"
	case BSHART:
		return "SHARTEEL"
	case BEDWIN:
		return "EDWIN"
	case BVICON:
		return "VICONIA"
	case BIMOEN:
		return "IMOEN"
	case BNEERA:
		return "NEERA"
	case BDORN:
		return "DORN"
	case BRASAAD:
		return "RASAAD"
	}
	return ""
}
