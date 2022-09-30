package utils

const (
	BRI                         = "BRI"                         // 002
	MANDIRI                     = "MANDIRI"                     // 008
	BNI                         = "BNI"                         // 009
	DANAMON                     = "DANAMON"                     // 011
	PERMATA                     = "PERMATA"                     // 013
	PERMATA_UUS                 = "PERMATA_UUS"                 // 013
	BCA                         = "BCA"                         // 014
	MAYBANK                     = "MAYBANK"                     // 016
	MAYBANK_SYR                 = "MAYBANK_SYR"                 // 016
	PANIN                       = "PANIN"                       // 019
	PANIN_SYR                   = "PANIN_SYR"                   // 019
	CIMB                        = "CIMB"                        // 022
	UOB                         = "UOB"                         // 023
	OCBC                        = "OCBC"                        // 028
	CITIBANK                    = "CITIBANK"                    // 031
	ARTHA                       = "ARTHA"                       // 037
	TOKYO                       = "TOKYO"                       // 042
	DBS                         = "DBS"                         // 046
	STANDARD_CHARTERED          = "STANDARD_CHARTERED"          // 050
	CAPITAL                     = "CAPITAL"                     // 054
	ANZ                         = "ANZ"                         // 061
	BOC                         = "BOC"                         // 069
	BUMI_ARTA                   = "BUMI_ARTA"                   // 076
	HSBC                        = "HSBC"                        // 087
	RABOBANK                    = "RABOBANK"                    // 089
	JTRUST                      = "JTRUST"                      // 095
	MAYAPADA                    = "MAYAPADA"                    // 097
	JAWA_BARAT                  = "JAWA_BARAT"                  // 110 not supported by xendit
	DKI                         = "DKI"                         // 111
	DAERAH_ISTIMEWA             = "DAERAH_ISTIMEWA"             // 112
	JAWA_TENGAH                 = "JAWA_TENGAH"                 // 113
	JAWA_TIMUR                  = "JAWA_TIMUR"                  // 114
	JAMBI                       = "JAMBI"                       // 115
	JAMBI_UUS                   = "JAMBI_UUS"                   // 115
	ACEH                        = "ACEH"                        // 116
	ACEH_UUS                    = "ACEH_UUS"                    // 116
	SUMATERA_UTARA              = "SUMATERA_UTARA"              // 117 not supported by xendit
	NAGARI                      = "NAGARI"                      // 118 not supported by xendit
	RIAU_DAN_KEPRI              = "RIAU_DAN_KEPRI"              // 119
	RIAU_DAN_KEPRI_UUS          = "RIAU_DAN_KEPRI_UUS"          // 119
	SUMSEL_DAN_BABEL            = "SUMSEL_DAN_BABEL"            // 120
	SUMSEL_DAN_BABEL_UUS        = "SUMSEL_DAN_BABEL_UUS"        // 120
	LAMPUNG                     = "LAMPUNG"                     // 121
	KALIMANTAN_SELATAN          = "KALIMANTAN_SELATAN"          // 122
	KALIMANTAN_BARAT            = "KALIMANTAN_BARAT"            // 123
	KALIMANTAN_TIMUR            = "KALIMANTAN_TIMUR"            // 124
	KALIMANTAN_TENGAH           = "KALIMANTAN_TENGAH"           // 125
	SULSELBAR                   = "SULSELBAR"                   // 126
	SULUT                       = "SULUT"                       // 127
	NUSA_TENGGARA_BARAT         = "NUSA_TENGGARA_BARAT"         // 128
	NUSA_TENGGARA_BARAT_UUS     = "NUSA_TENGGARA_BARAT_UUS"     // 128
	BALI                        = "BALI"                        // 129
	NUSA_TENGGARA_TIMUR         = "NUSA_TENGGARA_TIMUR"         // 130
	MALUKU                      = "MALUKU"                      // 131
	PAPUA                       = "PAPUA"                       // 132
	SULAWESI_TENGAH             = "SULAWESI_TENGAH"             // 134 not supported by xendit
	SULAWESI_UTARA              = "SULAWESI_UTARA"              // 135 not supported by xendit
	BANTEN                      = "BANTEN"                      // 137
	NUSANTARA_PARAHYANGAN       = "NUSANTARA_PARAHYANGAN"       // 145
	INDIA                       = "INDIA"                       // 146
	MUAMALAT                    = "MUAMALAT"                    // 147
	MESTIKA_DHARMA              = "MESTIKA_DHARMA"              // 151
	SHINHAN                     = "SHINHAN"                     // 152
	SINARMAS                    = "SINARMAS"                    // 153
	MASPION                     = "MASPION"                     // 157
	GANESHA                     = "GANESHA"                     // 161
	ICBC                        = "ICBC"                        // 164
	QNB_INDONESIA               = "QNB_INDONESIA"               // 167
	BTN                         = "BTN"                         // 200
	WOORI_SAUDARA               = "WOORI_SAUDARA"               // 212
	TABUNGAN_PENSIUNAN_NASIONAL = "TABUNGAN_PENSIUNAN_NASIONAL" // 213
	VICTORIA_SYR                = "VICTORIA_SYR"                // 405
	JABAR_BANTEN_SYARIAH        = "JABAR_BANTEN_SYARIAH"        // 425
	MEGA                        = "MEGA"                        // 426
	BUKOPIN                     = "BUKOPIN"                     // 441
	BUKOPIN_SYR                 = "BUKOPIN_SYR"                 // 441
	BSI                         = "BSI"                         // 451
	JASA_JAKARTA                = "JASA_JAKARTA"                // 472
	HANA                        = "HANA"                        // 484
	MNC_INTERNASIONAL           = "MNC_INTERNASIONAL"           // 485
	YUDHA_BHAKTI                = "YUDHA_BHAKTI"                // 490
	AGRONIAGA                   = "AGRONIAGA"                   // 494
	SBI_INDONESIA               = "SBI_INDONESIA"               // 498
	ROYAL                       = "ROYAL"                       // 501
	NATIONALNOBU                = "NATIONALNOBU"                // 503
	MEGA_SYR                    = "MEGA_SYR"                    // 506
	INA_PERDANA                 = "INA_PERDANA"                 // 513
	PRIMA_MASTER                = "PRIMA_MASTER"                // 520
	SAHABAT_SAMPOERNA           = "SAHABAT_SAMPOERNA"           // 523
	DINAR_INDONESIA             = "DINAR_INDONESIA"             // 526
	KESEJAHTERAAN_EKONOMI       = "KESEJAHTERAAN_EKONOMI"       // 535
	BCA_SYR                     = "BCA_SYR"                     // 536
	ARTOS                       = "ARTOS"                       // 542
	BTPN_SYARIAH                = "BTPN_SYARIAH"                // 547
	MULTI_ARTA_SENTOSA          = "MULTI_ARTA_SENTOSA"          // 548
	MAYORA                      = "MAYORA"                      // 553
	INDEX_SELINDO               = "INDEX_SELINDO"               // 555
	CNB                         = "CNB"                         // 559 not supported by xendit
	MANTAP                      = "MANTAP"                      // 564 not supported by xendit
	VICTORIA_INTERNASIONAL      = "VICTORIA_INTERNASIONAL"      // 566
	HARDA_INTERNASIONAL         = "HARDA_INTERNASIONAL"         // 567 not supported by xendit
	BPR_KS                      = "BPR_KS"                      // 688 not supported by xendit
	IBK                         = "IBK"                         // 945 not supported by xendit
	CTBC_INDONESIA              = "CTBC_INDONESIA"              // 949 not supported by xendit
	COMMONWEALTH                = "COMMONWEALTH"                // 950
	CCB                         = "CCB"                         // XXX ANTARDAERAH / WINDUR_KENTJANA not supported by fusindo
	ANTARDAERAH                 = "ANTARDAERAH"                 // 088 CCB
	WINDUR_KENTJANA             = "WINDUR_KENTJANA"             // 036 CCB
	DANA                        = "DANA"
	GOPAY                       = "GOPAY"
	SHOPEEPAY                   = "SHOPEEPAY"
	OVO                         = "OVO"
	LINK_AJA                    = "LINK_AJA"
	ALADIN                      = "ALADIN"
)

func ConvertBankCodeOY(bankName string) string {

	if bankName == BRI {
		return "002"
	}
	if bankName == MANDIRI {
		return "008"
	}
	if bankName == BNI {
		return "009"
	}
	if bankName == DANAMON {
		return "011"
	}
	if bankName == PERMATA {
		return "013"
	}
	if bankName == PERMATA_UUS {
		return "013"
	}
	if bankName == BCA {
		return "014"
	}
	if bankName == MAYBANK {
		return "016"
	}
	if bankName == MAYBANK_SYR {
		return "016"
	}
	if bankName == PANIN {
		return "019"
	}
	if bankName == PANIN_SYR {
		return "019"
	}
	if bankName == CIMB {
		return "022"
	}
	if bankName == UOB {
		return "023"
	}
	if bankName == OCBC {
		return "028"
	}
	if bankName == CITIBANK {
		return "031"
	}
	if bankName == MEGA {
		return "426"
	}
	if bankName == ANTARDAERAH {
		return "088"
	}
	if bankName == WINDUR_KENTJANA {
		return "036"
	}
	if bankName == ARTHA {
		return "037"
	}
	if bankName == TOKYO {
		return "042"
	}
	if bankName == DBS {
		return "046"
	}
	if bankName == STANDARD_CHARTERED {
		return "050"
	}
	if bankName == CAPITAL {
		return "054"
	}
	if bankName == ANZ {
		return "061"
	}
	if bankName == BOC {
		return "069"
	}
	if bankName == BUMI_ARTA {
		return "076"
	}
	if bankName == HSBC {
		return "087"
	}
	if bankName == RABOBANK {
		return "089"
	}
	if bankName == JTRUST {
		return "095"
	}
	if bankName == MAYAPADA {
		return "097"
	}
	if bankName == JAWA_BARAT {
		return "110"
	}
	if bankName == DKI {
		return "111"
	}
	if bankName == DAERAH_ISTIMEWA {
		return "112"
	}
	if bankName == JAWA_TENGAH {
		return "113"
	}
	if bankName == JAWA_TIMUR {
		return "114"
	}
	if bankName == JAMBI {
		return "115"
	}
	if bankName == JAMBI_UUS {
		return "115"
	}
	if bankName == ACEH {
		return "116"
	}
	if bankName == ACEH_UUS {
		return "116"
	}
	if bankName == SUMATERA_UTARA {
		return "117"
	}
	if bankName == NAGARI {
		return "118"
	}
	if bankName == RIAU_DAN_KEPRI {
		return "119"
	}
	if bankName == RIAU_DAN_KEPRI_UUS {
		return "119"
	}
	if bankName == SUMSEL_DAN_BABEL {
		return "120"
	}
	if bankName == SUMSEL_DAN_BABEL_UUS {
		return "120"
	}
	if bankName == LAMPUNG {
		return "121"
	}
	if bankName == KALIMANTAN_SELATAN {
		return "122"
	}
	if bankName == KALIMANTAN_BARAT {
		return "123"
	}
	if bankName == KALIMANTAN_TIMUR {
		return "124"
	}
	if bankName == KALIMANTAN_TENGAH {
		return "125"
	}
	if bankName == SULSELBAR {
		return "126"
	}
	if bankName == SULUT {
		return "127"
	}
	if bankName == NUSA_TENGGARA_BARAT {
		return "128"
	}
	if bankName == NUSA_TENGGARA_BARAT_UUS {
		return "128"
	}
	if bankName == BALI {
		return "129"
	}
	if bankName == NUSA_TENGGARA_TIMUR {
		return "130"
	}
	if bankName == MALUKU {
		return "131"
	}
	if bankName == PAPUA {
		return "132"
	}
	if bankName == SULAWESI_TENGAH {
		return "134"
	}
	if bankName == SULAWESI_UTARA {
		return "135"
	}
	if bankName == BANTEN {
		return "137"
	}
	if bankName == NUSANTARA_PARAHYANGAN {
		return "145"
	}
	if bankName == INDIA {
		return "146"
	}
	if bankName == MUAMALAT {
		return "147"
	}
	if bankName == MESTIKA_DHARMA {
		return "151"
	}
	if bankName == SHINHAN {
		return "152"
	}
	if bankName == SINARMAS {
		return "153"
	}
	if bankName == MASPION {
		return "157"
	}
	if bankName == GANESHA {
		return "161"
	}
	if bankName == ICBC {
		return "164"
	}
	if bankName == QNB_INDONESIA {
		return "167"
	}
	if bankName == BTN {
		return "200"
	}
	if bankName == WOORI_SAUDARA {
		return "212"
	}
	if bankName == TABUNGAN_PENSIUNAN_NASIONAL {
		return "213"
	}
	if bankName == VICTORIA_SYR {
		return "405"
	}
	if bankName == JABAR_BANTEN_SYARIAH {
		return "425"
	}
	if bankName == BUKOPIN {
		return "441"
	}
	if bankName == BUKOPIN_SYR {
		return "441"
	}
	if bankName == BSI {
		return "451"
	}
	if bankName == JASA_JAKARTA {
		return "472"
	}
	if bankName == HANA {
		return "484"
	}
	if bankName == MNC_INTERNASIONAL {
		return "485"
	}
	if bankName == YUDHA_BHAKTI {
		return "490"
	}
	if bankName == AGRONIAGA {
		return "494"
	}
	if bankName == SBI_INDONESIA {
		return "498"
	}
	if bankName == ROYAL {
		return "501"
	}
	if bankName == NATIONALNOBU {
		return "503"
	}
	if bankName == MEGA_SYR {
		return "506"
	}
	if bankName == INA_PERDANA {
		return "513"
	}
	if bankName == PRIMA_MASTER {
		return "520"
	}
	if bankName == SAHABAT_SAMPOERNA {
		return "523"
	}
	if bankName == DINAR_INDONESIA {
		return "526"
	}
	if bankName == KESEJAHTERAAN_EKONOMI {
		return "535"
	}
	if bankName == BCA_SYR {
		return "536"
	}
	if bankName == ARTOS {
		return "542"
	}
	if bankName == BTPN_SYARIAH {
		return "547"
	}
	if bankName == MULTI_ARTA_SENTOSA {
		return "548"
	}
	if bankName == MAYORA {
		return "553"
	}
	if bankName == INDEX_SELINDO {
		return "555"
	}
	if bankName == CNB {
		return "559"
	}
	if bankName == MANTAP {
		return "564"
	}
	if bankName == VICTORIA_INTERNASIONAL {
		return "566"
	}
	if bankName == HARDA_INTERNASIONAL {
		return "567"
	}
	if bankName == BPR_KS {
		return "688"
	}
	if bankName == IBK {
		return "945"
	}
	if bankName == CTBC_INDONESIA {
		return "949"
	}
	if bankName == COMMONWEALTH {
		return "950"
	}
	if bankName == DANA {
		return "dana"
	}
	if bankName == GOPAY {
		return "gopay"
	}
	if bankName == SHOPEEPAY {
		return "shopeepay"
	}
	if bankName == OVO {
		return "ovo"
	}
	if bankName == LINK_AJA {
		return "linkaja"
	}
	if bankName == ALADIN {
		return "947"
	}

	return ""
}

func ConvertBankCodeMMBC(bankName string) string {

	if bankName == BRI {
		return "2"
	}
	if bankName == MANDIRI {
		return "8"
	}
	if bankName == BNI {
		return "9"
	}
	if bankName == DANAMON {
		return "11"
	}
	if bankName == PERMATA {
		return "13"
	}
	if bankName == PERMATA_UUS {
		return "13"
	}
	if bankName == BCA {
		return "14"
	}
	if bankName == MAYBANK {
		return "16"
	}
	if bankName == MAYBANK_SYR {
		return "16"
	}
	if bankName == PANIN {
		return "19"
	}
	if bankName == PANIN_SYR {
		return "19"
	}
	if bankName == CIMB {
		return "22"
	}
	if bankName == UOB {
		return "23"
	}
	if bankName == OCBC {
		return "28"
	}
	if bankName == CITIBANK {
		return "31"
	}
	if bankName == MEGA {
		return "426"
	}
	if bankName == ANTARDAERAH {
		return "88"
	}
	if bankName == WINDUR_KENTJANA {
		return "36"
	}
	if bankName == ARTHA {
		return "37"
	}
	if bankName == TOKYO {
		return "42"
	}
	if bankName == DBS {
		return "46"
	}
	if bankName == STANDARD_CHARTERED {
		return "50"
	}
	if bankName == CAPITAL {
		return "54"
	}
	if bankName == ANZ {
		return "61"
	}
	if bankName == BOC {
		return "69"
	}
	if bankName == BUMI_ARTA {
		return "76"
	}
	if bankName == HSBC {
		return "87"
	}
	if bankName == RABOBANK {
		return "89"
	}
	if bankName == JTRUST {
		return "95"
	}
	if bankName == MAYAPADA {
		return "97"
	}
	if bankName == JAWA_BARAT {
		return "110"
	}
	if bankName == DKI {
		return "111"
	}
	if bankName == DAERAH_ISTIMEWA {
		return "112"
	}
	if bankName == JAWA_TENGAH {
		return "113"
	}
	if bankName == JAWA_TIMUR {
		return "114"
	}
	if bankName == JAMBI {
		return "115"
	}
	if bankName == JAMBI_UUS {
		return "115"
	}
	if bankName == ACEH {
		return "116"
	}
	if bankName == ACEH_UUS {
		return "116"
	}
	if bankName == SUMATERA_UTARA {
		return "117"
	}
	if bankName == NAGARI {
		return "118"
	}
	if bankName == RIAU_DAN_KEPRI {
		return "119"
	}
	if bankName == RIAU_DAN_KEPRI_UUS {
		return "119"
	}
	if bankName == SUMSEL_DAN_BABEL {
		return "120"
	}
	if bankName == SUMSEL_DAN_BABEL_UUS {
		return "120"
	}
	if bankName == LAMPUNG {
		return "121"
	}
	if bankName == KALIMANTAN_SELATAN {
		return "122"
	}
	if bankName == KALIMANTAN_BARAT {
		return "123"
	}
	if bankName == KALIMANTAN_TIMUR {
		return "124"
	}
	if bankName == KALIMANTAN_TENGAH {
		return "125"
	}
	if bankName == SULSELBAR {
		return "126"
	}
	if bankName == SULUT {
		return "127"
	}
	if bankName == NUSA_TENGGARA_BARAT {
		return "128"
	}
	if bankName == NUSA_TENGGARA_BARAT_UUS {
		return "128"
	}
	if bankName == BALI {
		return "129"
	}
	if bankName == NUSA_TENGGARA_TIMUR {
		return "130"
	}
	if bankName == MALUKU {
		return "131"
	}
	if bankName == PAPUA {
		return "132"
	}
	if bankName == SULAWESI_TENGAH {
		return "134"
	}
	if bankName == SULAWESI_UTARA {
		return "135"
	}
	if bankName == BANTEN {
		return "137"
	}
	if bankName == NUSANTARA_PARAHYANGAN {
		return "145"
	}
	if bankName == INDIA {
		return "146"
	}
	if bankName == MUAMALAT {
		return "147"
	}
	if bankName == MESTIKA_DHARMA {
		return "151"
	}
	if bankName == SHINHAN {
		return "152"
	}
	if bankName == SINARMAS {
		return "153"
	}
	if bankName == MASPION {
		return "157"
	}
	if bankName == GANESHA {
		return "161"
	}
	if bankName == ICBC {
		return "164"
	}
	if bankName == QNB_INDONESIA {
		return "167"
	}
	if bankName == BTN {
		return "200"
	}
	if bankName == WOORI_SAUDARA {
		return "212"
	}
	if bankName == TABUNGAN_PENSIUNAN_NASIONAL {
		return "213"
	}
	if bankName == VICTORIA_SYR {
		return "405"
	}
	if bankName == JABAR_BANTEN_SYARIAH {
		return "425"
	}
	if bankName == BUKOPIN {
		return "441"
	}
	if bankName == BUKOPIN_SYR {
		return "441"
	}
	if bankName == BSI {
		return "451"
	}
	if bankName == JASA_JAKARTA {
		return "472"
	}
	if bankName == HANA {
		return "484"
	}
	if bankName == MNC_INTERNASIONAL {
		return "485"
	}
	if bankName == YUDHA_BHAKTI {
		return "490"
	}
	if bankName == AGRONIAGA {
		return "494"
	}
	if bankName == SBI_INDONESIA {
		return "498"
	}
	if bankName == ROYAL {
		return "501"
	}
	if bankName == NATIONALNOBU {
		return "503"
	}
	if bankName == MEGA_SYR {
		return "506"
	}
	if bankName == INA_PERDANA {
		return "513"
	}
	if bankName == PRIMA_MASTER {
		return "520"
	}
	if bankName == SAHABAT_SAMPOERNA {
		return "523"
	}
	if bankName == DINAR_INDONESIA {
		return "526"
	}
	if bankName == KESEJAHTERAAN_EKONOMI {
		return "535"
	}
	if bankName == BCA_SYR {
		return "536"
	}
	if bankName == ARTOS {
		return "542"
	}
	if bankName == BTPN_SYARIAH {
		return "547"
	}
	if bankName == MULTI_ARTA_SENTOSA {
		return "548"
	}
	if bankName == MAYORA {
		return "553"
	}
	if bankName == INDEX_SELINDO {
		return "555"
	}
	if bankName == CNB {
		return "559"
	}
	if bankName == MANTAP {
		return "564"
	}
	if bankName == VICTORIA_INTERNASIONAL {
		return "566"
	}
	if bankName == HARDA_INTERNASIONAL {
		return "567"
	}
	if bankName == BPR_KS {
		return "688"
	}
	if bankName == IBK {
		return "945"
	}
	if bankName == CTBC_INDONESIA {
		return "949"
	}
	if bankName == COMMONWEALTH {
		return "950"
	}

	if bankName == DANA {
		return DANA
	}

	if bankName == GOPAY {
		return GOPAY
	}

	if bankName == SHOPEEPAY {
		return SHOPEEPAY
	}

	if bankName == OVO {
		return OVO
	}

	if bankName == LINK_AJA {
		return "LINKAJA"
	}

	if bankName == ALADIN {
		return "947"
	}

	return ""
}
