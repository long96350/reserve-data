package configuration

import "github.com/KyberNetwork/reserve-data/common"

//FeeConfigs store predefined fee configs of exchanges
var FeeConfigs = map[string]common.ExchangeFees{
	"binance": {
		Trading: common.TradingFee{
			"taker": 0.001,
			"maker": 0.001,
		},
		Funding: common.FundingFee{
			Withdraw: map[string]float64{
				"ADA":   1,
				"ADX":   4,
				"AE":    2.3,
				"AION":  1.6,
				"AMB":   9.4,
				"APPC":  7.8,
				"ARK":   0.1,
				"ARN":   2.3,
				"AST":   11.9,
				"BAT":   14,
				"BCC":   0.001,
				"BCD":   1,
				"BCPT":  4.8,
				"BCX":   1,
				"BLZ":   10.2,
				"BNB":   0.51,
				"BNT":   1.1,
				"BQX":   1.2,
				"BRD":   6.4,
				"BTC":   0.0005,
				"BTG":   0.001,
				"BTM":   5,
				"BTS":   1,
				"CDT":   72,
				"CHAT":  57.3,
				"CMT":   40,
				"CND":   40,
				"CTR":   7.5,
				"DAI":   0,
				"DASH":  0.002,
				"DGD":   0.02,
				"DLT":   16.3,
				"DNT":   56,
				"EDO":   2,
				"ELF":   4.4,
				"ENG":   2,
				"ENJ":   29,
				"EOS":   0.7,
				"ETC":   0.01,
				"ETF":   1,
				"ETH":   0.01,
				"EVX":   2.5,
				"FUEL":  44,
				"FUN":   82,
				"GAS":   0,
				"GTO":   15,
				"GVT":   0.15,
				"GXS":   0.3,
				"HCC":   0.0005,
				"HSR":   0.0001,
				"ICN":   2.9,
				"ICX":   1.5,
				"INS":   2.4,
				"IOST":  145.7,
				"IOTA":  0.5,
				"KMD":   0.002,
				"KNC":   2.6,
				"LEND":  60,
				"LINK":  8.7,
				"LRC":   9.4,
				"LSK":   0.1,
				"LTC":   0.01,
				"LUN":   0.24,
				"MANA":  46,
				"MCO":   0.56,
				"MDA":   4.1,
				"MOD":   2,
				"MTH":   36,
				"MTL":   0.9,
				"NANO":  0.01,
				"NAV":   0.2,
				"NCASH": 160.6,
				"NEBL":  0.01,
				"NEO":   0,
				"NULS":  1.4,
				"OAX":   5.7,
				"OMG":   0.3,
				"ONT":   1,
				"OST":   16,
				"PIVX":  0.02,
				"POA":   0.01,
				"POE":   159,
				"POWR":  9.1,
				"PPT":   0.25,
				"QSP":   22,
				"QTUM":  0.01,
				"RCN":   35,
				"RDN":   2,
				"REP":   0.1,
				"REQ":   17.3,
				"RLC":   3.8,
				"RPX":   1,
				"SALT":  1.3,
				"SBTC":  1,
				"SNGLS": 49,
				"SNM":   26,
				"SNT":   30,
				"STEEM": 0.01,
				"STORJ": 4.6,
				"STRAT": 0.1,
				"SUB":   12.3,
				"TNB":   108,
				"TNT":   47,
				"TUSD":  2.6,
				"TRIG":  15.7,
				"TRX":   110,
				"USDT":  11.5,
				"VEN":   1.1,
				"VIA":   0.01,
				"VIB":   20,
				"VIBE":  13.2,
				"WABI":  3.8,
				"WAVES": 0.002,
				"WINGS": 6.7,
				"WTC":   0.3,
				"XLM":   0.01,
				"XMR":   0.04,
				"XRP":   0.25,
				"XVG":   0.1,
				"XZC":   0.02,
				"YOYO":  42,
				"ZEC":   0.005,
				"ZIL":   100,
				"ZRX":   5.8,
			},
			Deposit: map[string]float64{
				"ADA":   0,
				"ADX":   0,
				"AE":    0,
				"AION":  0,
				"AMB":   0,
				"APPC":  0,
				"ARK":   0,
				"ARN":   0,
				"AST":   0,
				"BAT":   0,
				"BCC":   0,
				"BCD":   0,
				"BCPT":  0,
				"BCX":   0,
				"BLZ":   0,
				"BNB":   0,
				"BNT":   0,
				"BQX":   0,
				"BRD":   0,
				"BTC":   0,
				"BTG":   0,
				"BTM":   0,
				"BTS":   0,
				"CDT":   0,
				"CHAT":  0,
				"CMT":   0,
				"CND":   0,
				"CTR":   0,
				"DAI":   0,
				"DASH":  0,
				"DGD":   0,
				"DLT":   0,
				"DNT":   0,
				"EDO":   0,
				"ELF":   0,
				"ENG":   0,
				"ENJ":   0,
				"EOS":   0,
				"ETC":   0,
				"ETF":   0,
				"ETH":   0,
				"EVX":   0,
				"FUEL":  0,
				"FUN":   0,
				"GAS":   0,
				"GTO":   0,
				"GVT":   0,
				"GXS":   0,
				"HCC":   0,
				"HSR":   0,
				"ICN":   0,
				"ICX":   0,
				"INS":   0,
				"IOST":  0,
				"IOTA":  0,
				"KMD":   0,
				"KNC":   0,
				"LEND":  0,
				"LINK":  0,
				"LRC":   0,
				"LSK":   0,
				"LTC":   0,
				"LUN":   0,
				"MANA":  0,
				"MCO":   0,
				"MDA":   0,
				"MOD":   0,
				"MTH":   0,
				"MTL":   0,
				"NANO":  0,
				"NAV":   0,
				"NCASH": 0,
				"NEBL":  0,
				"NEO":   0,
				"NULS":  0,
				"OAX":   0,
				"OMG":   0,
				"ONT":   0,
				"OST":   0,
				"PIVX":  0,
				"POA":   0,
				"POE":   0,
				"POWR":  0,
				"PPT":   0,
				"QSP":   0,
				"QTUM":  0,
				"RCN":   0,
				"RDN":   0,
				"REP":   0,
				"REQ":   0,
				"RLC":   0,
				"RPX":   0,
				"SALT":  0,
				"SBTC":  0,
				"SNGLS": 0,
				"SNM":   0,
				"SNT":   0,
				"STEEM": 0,
				"STORJ": 0,
				"STRAT": 0,
				"SUB":   0,
				"TNB":   0,
				"TNT":   0,
				"TUSD":  0,
				"TRIG":  0,
				"TRX":   0,
				"USDT":  0,
				"VEN":   0,
				"VIA":   0,
				"VIB":   0,
				"VIBE":  0,
				"WABI":  0,
				"WAVES": 0,
				"WINGS": 0,
				"WTC":   0,
				"XLM":   0,
				"XMR":   0,
				"XRP":   0,
				"XVG":   0,
				"XZC":   0,
				"YOYO":  0,
				"ZEC":   0,
				"ZIL":   0,
				"ZRX":   0,
			},
		},
	},
	"bittrex": {
		Trading: common.TradingFee{
			"taker": 0.0025,
			"maker": 0.0025,
		},
		Funding: common.FundingFee{
			Withdraw: map[string]float64{
				"ADX":   3,
				"AE":    0,
				"ANT":   1.3,
				"APPC":  0,
				"AST":   0,
				"BAT":   10,
				"BNB":   0,
				"BNT":   8,
				"BQX":   0,
				"BTM":   0,
				"CFI":   27,
				"CTR":   0,
				"CVC":   8,
				"DAI":   0,
				"DGD":   0.038,
				"ELF":   0,
				"ENG":   1,
				"EOS":   0,
				"ETH":   0.006,
				"ETHOS": 0,
				"FUN":   49,
				"GNO":   0.02,
				"GNT":   2,
				"GTO":   0,
				"ICN":   0,
				"KNC":   0,
				"LINK":  0,
				"LRC":   0,
				"MANA":  36,
				"MCO":   0.5,
				"MLN":   0.035,
				"MTL":   1.35,
				"OMG":   0.35,
				"PAY":   2,
				"POWR":  5,
				"PPT":   0,
				"QRL":   2.5,
				"RCN":   16,
				"RDN":   0,
				"REP":   0.1,
				"REQ":   0,
				"RLC":   3.5,
				"SALT":  0.6,
				"SNGLS": 3.5,
				"SNT":   20,
				"SONM":  0,
				"STOX":  0,
				"TAAS":  0,
				"TRX":   1,
				"VERI":  0,
				"WINGS": 4,
				"WTC":   0,
				"ZRX":   1,
			},
			Deposit: map[string]float64{
				"ADX":   0,
				"AE":    0,
				"ANT":   0,
				"APPC":  0,
				"AST":   0,
				"BAT":   0,
				"BNB":   0,
				"BNT":   0,
				"BQX":   0,
				"BTM":   0,
				"CFI":   0,
				"CTR":   0,
				"CVC":   0,
				"DAI":   0,
				"DGD":   0,
				"ELF":   0,
				"ENG":   0,
				"EOS":   0,
				"ETH":   0,
				"FUN":   0,
				"GNO":   0,
				"GNT":   0,
				"GTO":   0,
				"ICN":   0,
				"KNC":   0,
				"LINK":  0,
				"LRC":   0,
				"MANA":  0,
				"MCO":   0,
				"MLN":   0,
				"MTL":   0,
				"OMG":   0,
				"PAY":   0,
				"POWR":  0,
				"PPT":   0,
				"QRL":   0,
				"RCN":   0,
				"RDN":   0,
				"REP":   0,
				"REQ":   0,
				"RLC":   0,
				"SALT":  0,
				"SNGLS": 0,
				"SNT":   0,
				"SONM":  0,
				"STOX":  0,
				"TAAS":  0,
				"TRX":   0,
				"VERI":  0,
				"WINGS": 0,
				"WTC":   0,
				"ZRX":   0,
			},
		},
	},
	"huobi": {
		Trading: common.TradingFee{
			"taker": 0.002,
			"maker": 0.002,
		},
		Funding: common.FundingFee{
			Withdraw: map[string]float64{
				"ABT":   2,
				"ACT":   0.01,
				"ADX":   0.5,
				"AIDOC": 10,
				"APPC":  0.5,
				"AST":   5,
				"BAT":   5,
				"BCH":   0.0001,
				"BLZ":   2,
				"BQX":   0,
				"BTC":   0.001,
				"BTM":   2,
				"CHAT":  2,
				"CMT":   20,
				"CVC":   2,
				"DAI":   0,
				"DASH":  0.002,
				"DAT":   10,
				"DBC":   10,
				"DGD":   0.01,
				"DTA":   100,
				"EDU":   500,
				"EKO":   20,
				"ELA":   0.005,
				"ELF":   5,
				"ENG":   0.5,
				"EOS":   0.5,
				"ETC":   0.01,
				"ETH":   0.01,
				"EVX":   0.5,
				"GAS":   0,
				"GNT":   5,
				"GNX":   5,
				"GTO":   0,
				"HSR":   0.0001,
				"HT":    1,
				"ICX":   0.2,
				"ITC":   2,
				"KNC":   1,
				"LBA":   10,
				"LET":   30,
				"LINK":  1,
				"LSK":   0.1,
				"LTC":   0.001,
				"LUN":   0.05,
				"MANA":  10,
				"MCO":   0.2,
				"MDS":   20,
				"MEE":   10,
				"MTL":   0.2,
				"MTN":   5,
				"MTX":   2,
				"NAS":   0.2,
				"NEO":   0,
				"OCN":   100,
				"OMG":   0.1,
				"OST":   1,
				"PAY":   0.5,
				"POLY":  1,
				"POWR":  2,
				"PROPY": 0.5,
				"QSB":   10,
				"QTUM":  0.01,
				"QUASH": 1,
				"QUN":   30,
				"RCN":   10,
				"RDN":   1,
				"REQ":   5,
				"RPX":   2,
				"RUFF":  20,
				"SALT":  0.1,
				"SMT":   50,
				"SNC":   5,
				"SNT":   50,
				"SOC":   10,
				"SRN":   0.5,
				"STK":   10,
				"STORJ": 2,
				"SWFTC": 100,
				"THETA": 10,
				"TNB":   50,
				"TNT":   20,
				"TOPC":  20,
				"TRX":   20,
				"USDT":  20,
				"UTK":   2,
				"VEN":   2,
				"WAX":   1,
				"WICC":  2,
				"WPR":   10,
				"XEM":   4,
				"XRP":   0.1,
				"YEE":   50,
				"ZEC":   0.001,
				"ZIL":   100,
				"ZLA":   1,
				"ZRX":   5,
			},
			Deposit: map[string]float64{
				"ABT":   0,
				"ACT":   0,
				"ADX":   0,
				"AIDOC": 0,
				"APPC":  0,
				"AST":   0,
				"BAT":   0,
				"BCH":   0,
				"BLZ":   0,
				"BQX":   0,
				"BTC":   0,
				"BTM":   0,
				"CHAT":  0,
				"CMT":   0,
				"CVC":   0,
				"DAI":   0,
				"DASH":  0,
				"DAT":   0,
				"DBC":   0,
				"DGD":   0,
				"DTA":   0,
				"EDU":   0,
				"EKO":   0,
				"ELA":   0,
				"ELF":   0,
				"ENG":   0,
				"EOS":   0,
				"ETC":   0,
				"ETH":   0,
				"EVX":   0,
				"GAS":   0,
				"GNT":   0,
				"GNX":   0,
				"GTO":   0,
				"HSR":   0,
				"HT":    0,
				"ICX":   0,
				"IOST":  0,
				"ITC":   0,
				"KNC":   0,
				"LBA":   0,
				"LET":   0,
				"LINK":  0,
				"LSK":   0,
				"LTC":   0,
				"LUN":   0,
				"MANA":  0,
				"MCO":   0,
				"MDS":   0,
				"MEE":   0,
				"MTL":   0,
				"MTN":   0,
				"MTX":   0,
				"NAS":   0,
				"NEO":   0,
				"OCN":   0,
				"OMG":   0,
				"OST":   0,
				"PAY":   0,
				"POLY":  0,
				"POWR":  0,
				"PROPY": 0,
				"QSB":   0,
				"QTUM":  0,
				"QUASH": 0,
				"QUN":   0,
				"RCN":   0,
				"RDN":   0,
				"REQ":   0,
				"RPX":   0,
				"RUFF":  0,
				"SALT":  0,
				"SMT":   0,
				"SNC":   0,
				"SNT":   0,
				"SOC":   0,
				"SRN":   0,
				"STK":   0,
				"STORJ": 0,
				"SWFTC": 0,
				"THETA": 0,
				"TNB":   0,
				"TNT":   0,
				"TOPC":  0,
				"TRX":   0,
				"USDT":  0,
				"UTK":   0,
				"VEN":   0,
				"WAX":   0,
				"WICC":  0,
				"WPR":   0,
				"XEM":   0,
				"XRP":   0,
				"YEE":   0,
				"ZEC":   0,
				"ZIL":   0,
				"ZLA":   0,
				"ZRX":   0,
			},
		},
	},
	"stable_exchange": {
		Trading: common.TradingFee{
			"taker": 0,
			"maker": 0,
		},
		Funding: common.FundingFee{
			Withdraw: map[string]float64{
				"ETH": 0,
				"DGX": 0,
			},
			Deposit: map[string]float64{
				"ETH": 0,
				"DGX": 0,
			},
		},
	},
}