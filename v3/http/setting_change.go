package http

import (
	"fmt"
	"log"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"

	v1common "github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/v3/common"
)

func (s *Server) validateChangeEntry(e common.SettingChangeType, changeType common.ChangeType) error {
	var (
		err error
	)

	switch changeType {
	case common.ChangeTypeCreateAsset:
		err = s.checkCreateAssetParams(*(e.(*common.CreateAssetEntry)))
	case common.ChangeTypeUpdateAsset:
		err = s.checkUpdateAssetParams(*(e.(*common.UpdateAssetEntry)))
	case common.ChangeTypeCreateAssetExchange:
		err = s.checkCreateAssetExchangeParams(*(e.(*common.CreateAssetExchangeEntry)))
	case common.ChangeTypeUpdateAssetExchange:
		err = s.checkUpdateAssetExchangeParams(*(e.(*common.UpdateAssetExchangeEntry)))
	case common.ChangeTypeCreateTradingPair:
		_, _, err = s.checkCreateTradingPairParams(*(e.(*common.CreateTradingPairEntry)))
	case common.ChangeTypeCreateTradingBy:
		err = s.checkCreateTradingByParams(*e.(*common.CreateTradingByEntry))
	case common.ChangeTypeChangeAssetAddr:
		err = s.checkChangeAssetAddressParams(*e.(*common.ChangeAssetAddressEntry))
	case common.ChangeTypeUpdateExchange:
		return nil
	case common.ChangeTypeDeleteTradingPair:
		err = s.checkDeleteTradingPairParams(*e.(*common.DeleteTradingPairEntry))
	case common.ChangeTypeDeleteAssetExchange:
		err = s.checkDeleteAssetExchangeParams(*e.(*common.DeleteAssetExchangeEntry))
	default:
		return errors.Errorf("unknown type of setting change: %v", reflect.TypeOf(e))
	}
	return err
}

func (s *Server) fillLiveInfoSettingChange(settingChange *common.SettingChange) error {
	assets, err := s.storage.GetAssets()
	if err != nil {
		return err
	}

	for _, o := range settingChange.ChangeList {
		switch o.Type {
		case common.ChangeTypeCreateAsset:
			asset := o.Data.(*common.CreateAssetEntry)
			for _, assetExchange := range asset.Exchanges {
				err = s.fillLiveInfoAssetExchange(assets, assetExchange.ExchangeID, assetExchange.TradingPairs, assetExchange.Symbol, assetExchange.AssetID)
				if err != nil {
					return err
				}
			}
		case common.ChangeTypeCreateTradingPair:
			entry := o.Data.(*common.CreateTradingPairEntry)
			baseSymbol, quoteSymbol, err := s.checkCreateTradingPairParams(*entry)
			if err != nil {
				return err
			}
			tradingPairSymbol := common.TradingPairSymbols{TradingPair: entry.TradingPair}
			tradingPairSymbol.BaseSymbol = baseSymbol
			tradingPairSymbol.QuoteSymbol = quoteSymbol
			tradingPairSymbol.ID = uint64(1)
			exhID := v1common.ExchangeID(entry.ExchangeID)
			centralExh, ok := s.supportedExchanges[exhID]
			if !ok {
				return errors.Errorf("exchange %s not supported", exhID)
			}
			exchangeInfo, err := centralExh.GetLiveExchangeInfos([]common.TradingPairSymbols{tradingPairSymbol})
			if err != nil {
				return err
			}
			info := exchangeInfo[1]
			entry.MinNotional = info.MinNotional
			entry.AmountLimitMax = info.AmountLimit.Max
			entry.AmountLimitMin = info.AmountLimit.Min
			entry.AmountPrecision = uint64(info.Precision.Amount)
			entry.PricePrecision = uint64(info.Precision.Price)
			entry.PriceLimitMax = info.PriceLimit.Max
			entry.PriceLimitMin = info.PriceLimit.Min
		case common.ChangeTypeCreateAssetExchange:
			assetExchange := o.Data.(*common.CreateAssetExchangeEntry)
			err = s.fillLiveInfoAssetExchange(assets, assetExchange.ExchangeID, assetExchange.TradingPairs, assetExchange.Symbol, assetExchange.AssetID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Server) fillLiveInfoAssetExchange(assets []common.Asset, exchangeID uint64, tradingPairs []common.TradingPair, assetSymbol string, assetID uint64) error {
	exhID := v1common.ExchangeID(exchangeID)
	centralExh, ok := s.supportedExchanges[exhID]
	if !ok {
		return errors.Errorf("exchange %s not supported", exhID)
	}
	var tps []common.TradingPairSymbols
	index := uint64(1)
	for idx, tradingPair := range tradingPairs {
		tradingPairSymbol := common.TradingPairSymbols{TradingPair: tradingPair}
		tradingPairSymbol.ID = index
		if tradingPair.Quote == 0 {
			tradingPairSymbol.QuoteSymbol = assetSymbol
			base, err := getAssetExchange(assets, tradingPair.Base, exchangeID)
			if err != nil {
				return err
			}
			tradingPairSymbol.BaseSymbol = base.Symbol
			if assetID != 0 {
				tradingPairs[idx].Quote = assetID
			}
		}
		if tradingPair.Base == 0 {
			tradingPairSymbol.BaseSymbol = assetSymbol
			quote, err := getAssetExchange(assets, tradingPair.Quote, exchangeID)
			if err != nil {
				return err
			}
			tradingPairSymbol.QuoteSymbol = quote.Symbol
			if assetID != 0 {
				tradingPairs[idx].Base = assetID
			}
		}
		tps = append(tps, tradingPairSymbol)
		index++
	}
	exchangeInfo, err := centralExh.GetLiveExchangeInfos(tps)
	if err != nil {
		return err
	}
	tradingPairID := uint64(1)
	for idx := range tradingPairs {
		if info, ok := exchangeInfo[tradingPairID]; ok {
			tradingPairs[idx].MinNotional = info.MinNotional
			tradingPairs[idx].AmountLimitMax = info.AmountLimit.Max
			tradingPairs[idx].AmountLimitMin = info.AmountLimit.Min
			tradingPairs[idx].AmountPrecision = uint64(info.Precision.Amount)
			tradingPairs[idx].PricePrecision = uint64(info.Precision.Price)
			tradingPairs[idx].PriceLimitMax = info.PriceLimit.Max
			tradingPairs[idx].PriceLimitMin = info.PriceLimit.Min
			tradingPairID++
		}
	}
	return nil
}

func (s *Server) createSettingChange(c *gin.Context) {
	var settingChange common.SettingChange
	if err := c.ShouldBindJSON(&settingChange); err != nil {
		log.Printf("cannot bind data to create setting_change from request err=%s", err.Error())
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	for i, o := range settingChange.ChangeList {
		if err := binding.Validator.ValidateStruct(o.Data); err != nil {
			msg := fmt.Sprintf("validate obj error at %d, err=%s", i, err)
			httputil.ResponseFailure(c, httputil.WithError(err), httputil.WithReason(msg))
			return
		}

		if err := s.validateChangeEntry(o.Data, o.Type); err != nil {
			msg := fmt.Sprintf("validate error at %d, err=%s", i, err)
			log.Println(msg)
			httputil.ResponseFailure(c, httputil.WithError(err), httputil.WithReason(msg))
			return
		}
	}
	if err := s.fillLiveInfoSettingChange(&settingChange); err != nil {
		msg := fmt.Sprintf("fill live info error=%s", err)
		log.Println(msg)
		httputil.ResponseFailure(c, httputil.WithError(err), httputil.WithReason(msg))
		return
	}

	id, err := s.storage.CreateSettingChange(settingChange)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	// test confirm
	err = s.storage.ConfirmSettingChange(id, false)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithField("id", id))
}

func (s *Server) getSettingChange(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		log.Println(err)
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	result, err := s.storage.GetSettingChange(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(result))
}

func (s *Server) getSettingChanges(c *gin.Context) {
	result, err := s.storage.GetSettingChanges()
	if err != nil {
		log.Printf("failed to get setting changes %v\n", err)
		httputil.ResponseFailure(c, httputil.WithError(err))
	}
	httputil.ResponseSuccess(c, httputil.WithData(result))
}

func (s *Server) rejectSettingChange(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	err := s.storage.RejectSettingChange(input.ID)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

func (s *Server) confirmSettingChange(c *gin.Context) {
	var input struct {
		ID uint64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&input); err != nil {
		log.Println(err)
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	err := s.storage.ConfirmSettingChange(input.ID, true)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}
