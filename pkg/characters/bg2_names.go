package characters

const (
	BAERIE BG2_Characters = iota
	BAERIE25
	BANOME25
	BANOMEN
	BCERND
	BCERND25
	BDORN25
	BEDWIN25
	BHAERD25
	BHAERDA
	BHEXXA25
	BHEXXAT
	BIMOEN2
	BIMOEN25
	BJAHEI25
	BJAHEIR
	BJAN
	BJAN25
	BKELDO25
	BKELDOR
	BKORGA25
	BKORGAN
	BMAZZY
	BMAZZY25
	BMINSC25
	BNALIA
	BNALIA25
	BNEERA25
	BRASAA25
	BSAREV25
	BVALYG25
	BVALYGA
	BVICON25
	BVICONI
	BYOSHIM
	CERNDJ
	DORNJ
	WILSON
	ZDBAE25B
	ZDBAEB
)

type BG2_Characters int64

func (s BG2_Characters) String() string {
	switch s {
	case BAERIE, BAERIE25:
		return "Aerie"
	case BANOMEN, BANOME25:
		return "Anomen"
	case BCERND, BCERND25:
		return "Cernd"
	case BEDWIN25:
		return "Edwin"
	case BJAHEIR, BJAHEI25:
		return "Jaheria"
	case BJAN, BJAN25:
		return "Jan"
	case BKELDOR, BKELDO25:
		return "Keldorn"
	case BKORGAN, BKORGA25:
		return "Korgan"
	case BHAERDA, BHAERD25:
		return "Haedalis"
	case BMAZZY, BMAZZY25:
		return "Mazzy"
	case BMINSC25:
		return "Minsc"
	case BNALIA, BNALIA25:
		return "Nalia"
	case BHEXXAT, BHEXXA25:
		return "Hexxat"
	case BIMOEN2, BIMOEN25:
		return "Imeon"
	case BNEERA25:
		return "Neera"
	case BDORN25:
		return "Dorn"
	case BRASAA25:
		return "Rasaad"
	case BSAREV25:
		return "Sarevok"
	case BVALYGA, BVALYG25:
		return "Valygar"
	case BVICONI, BVICON25:
		return "Viconia"
	case BYOSHIM:
		return "Yoshimo"
	case WILSON:
		return "Wilson"
	case ZDBAEB, ZDBAE25B:
		return "Baeloth"
	}
	return ""
}
