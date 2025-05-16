package utils

import "time"

var OzonProducerNames = []string{`ООО "Озон"`, "Озон ООО/Озон Фарм ООО", "Озон ООО", "Озон Фарм ООО", "Озон Фарм, Россия", "Озон, Россия", "Озон/Озон Фарм, Россия"}
var AgroFarmProducerNames = []string{`Арго-Фарм`, `АРГО-ФАРМ ООО`, `Арго-Фарм, Россия`}
var KirovFFProducerNames = []string{`Кировская фармацевтическая фабрика`, `Кировская фармацевтическая фабрика АО`, `Кировская фармфабрика, Россия`}

var RequestDelay = 2 * time.Second

var ZSRegions = map[string]string{
	"moscowregion": "Москва и область",
	"spb":          "Санкт-Петербург и область",
	"novosibirsk":  "Новосибирск",
	"ekb":          "Екатеринбург",
	"kazan":        "Казань",
	"nnovgorod":    "Нижний Новгород",
	"ufa":          "Уфа",
	"rostov":       "Ростов-на-Дону",
}

var ARRegions = map[string]string{
	"5e57803249af4c0001d64407": "Москва и область",
	"5e5778f0f85fde0001489df1": "Санкт-Петербург",
	"5e57858e2b690a0001b0977f": "Новосибирск",
	"5e58b8fee8d43f0001a55baf": "Екатеринбург",
	"6681727f3210fbd198563743": "Казань",
	"5e58a675cd37fa0001a7cb79": "Нижний Новгород",
	"5e57acd4752ac70001593b7f": "Уфа",
	"5e577e97a2efde00019e63e4": "Ростов-на-Дону",
}

var EARegions = map[string]string{
	"msk":            "Москва и область",
	"spb":            "Санкт-Петербург и область",
	"novosibirsk":    "Новосибирск",
	"ekb":            "Екатеринбург",
	"kazan":          "Казань",
	"nn":             "Нижний Новгород",
	"ufa":            "Уфа",
	"rostov-na-donu": "Ростов-на-Дону",
}

var PillsList = map[string]string{
	"antistenMV":                 "Антистен",
	"artrafik":                   "артрафик",
	"cerpehol":                   "Церпехол",
	"citipigam":                  "Цитипигам",
	"detravenol":                 "Детравенол",
	"doksiciklin":                "Доксициклин",
	"jekurohol":                  "Экурохол",
	"jekziter":                   "Экзитер",
	"jeljufor":                   "Элюфор",
	"esslial":                    "эсслиал форте",
	"estetatet":                  "эстет а тет",
	"jezlor":                     "Эзлор",
	"fazostabil":                 "Фазостабил",
	"flebofa":                    "Флебофа",
	"frejmitus":                  "Фреймитус",
	"holikron":                   "Холикрон",
	"izislip":                    "Изислип",
	"kontrazud":                  "Контразуд",
	"ksilometazolin":             "Ксилометазолин",
	"kruoksaban":                 "круоксабан",
	"magnepasit":                 "Магнепасит",
	"mebespalin":                 "Мебеспалин",
	"meladapt":                   "Меладапт",
	"memoritab":                  "Меморитаб",
	"motilorus":                  "Мотилорус",
	"mukocil":                    "Мукоцил",
	"naftoderil":                 "Нафтодерил",
	"nekspra":                    "Некспра",
	"nikurilly":                  "Никуриллы",
	"noocil":                     "Нооцил",
	"orungamin":                  "Орунгамин",
	"pronokognil":                "Пронокогнил",
	"safiston":                   "Сафистон",
	"seltavir":                   "Сельтавир",
	"skinoklir":                  "Скиноклир",
	"sumatrolid":                 "Суматролид",
	"tiloram":                    "Тилорам",
	"tolizor":                    "Толизор",
	"toraksol-soljushn-tablets":  "Тораксол солюшн таблетс",
	"ulblok":                     "Ульблок",
	"vildegra":                   "Вилдегра",
	"vismuta-trikalija-dicitrat": "Висмута трикалия дицитрат",
}
