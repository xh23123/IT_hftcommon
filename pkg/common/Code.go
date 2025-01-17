package common

type ExchangeID string

type AccountIdx int

type SymbolID string

type TransactionID string

type DataID string

type ActionID string

type PositionID string

type SideID string

type OrderTypeID string

type OrderStatusID string

type StatusID int

type ResetID string

type IntervalID string

type OrderidID string

type ClientOrderidID string

type AmendTypeID string

// ExchangeID
const BINANCEID ExchangeID = "1"
const OKEXID ExchangeID = "3"
const COINBASEID ExchangeID = "4"
const COINEXID ExchangeID = "5"
const BYBITID ExchangeID = "6"
const ETHEREUMDEXID ExchangeID = "7"
const KRAKENID ExchangeID = "8"
const MOCKEXCHANGEID ExchangeID = "9"
const GATEID ExchangeID = "10"
const BITGETID ExchangeID = "11"
const LBANKID ExchangeID = "12"
const WOOID ExchangeID = "13"
const KUCOINID ExchangeID = "14"
const PHEMEXID ExchangeID = "15"
const HUOBIID ExchangeID = "16"

// TransactionID
const SpotID TransactionID = "1"
const FutureID TransactionID = "2"
const MarginID TransactionID = "3"
const CoinFutureID TransactionID = "4"
const UniMarginID TransactionID = "5"

// DataID

const OrderUpdateID DataID = "1"
const TradeUpdateID DataID = "2"
const BookTickID DataID = "3"
const DepthID DataID = "4"
const TickID DataID = "5"
const KlineWsID DataID = "6"
const MarkPriceID DataID = "7"
const ErrorID DataID = "8"
const AggTradeID DataID = "9"
const UnknownDataID DataID = "10"
const RestDataID DataID = "11"
const TradeID DataID = "12"
const DexBookTickID DataID = "13"
const DexTradeID DataID = "14"
const OrderbookID DataID = "15"
const BalancesUpdateID DataID = "16"
const PositionsUpdateID DataID = "17"
const UpdateRestOrderID DataID = "18"
const MiscEventID DataID = "19"

// ActionID
const CREATE_LIMITTYPE_SPOT_ORDER ActionID = "1"
const CREATE_MARKET_SPOT_ORDER ActionID = "2"
const CANCEL_SPOT_ORDER ActionID = "3"
const CANCEL_ALL_SPOT_ORDER ActionID = "4"
const CREATE_LIMITTYPE_FUTURE_ORDER ActionID = "5"
const CREATE_MARKET_FUTURE_ORDER ActionID = "6"
const CANCEL_FUTURE_ORDER ActionID = "7"
const CANCEL_ALL_FUTURE_ORDER ActionID = "8"
const SET_MARGINTYPE ActionID = "9"
const SET_LEVERAGE ActionID = "10"
const SET_DUEL_SIDE_POSITION ActionID = "11"
const UNIVERSAL_TRANSFER ActionID = "12"
const RESET_USER_WS ActionID = "13"
const RESET_MARKET_WS ActionID = "14"
const RESET_TRIG_ONOFF ActionID = "15"
const SET_MULTIASSETMARGIN ActionID = "16"
const UNKNOWN_ACTION ActionID = "17"
const CREATE_LIMITTYPE_MARGIN_ORDER ActionID = "18"
const CANCEL_ALL_MARGIN_ORDER ActionID = "19"
const CANCEL_MARGIN_ORDER ActionID = "20"
const CREATE_LIMITTYPE_COIN_FUTURE_ORDER ActionID = "21"
const CANCEL_COIN_FUTURE_ORDER ActionID = "22"
const CANCEL_ALL_COIN_FUTURE_ORDER ActionID = "23"
const CREATE_AMEND_SPOT_ORDER ActionID = "24"
const CREATE_AMEND_FUTURE_ORDER ActionID = "25"
const CREATE_AMEND_COIN_FUTURE_ORDER ActionID = "26"
const CREATE_AMEND_MARGIN_ORDER ActionID = "27"
const CREATE_IOC_SPOT_ORDER ActionID = "28"
const CREATE_IOC_FUTURE_ORDER ActionID = "29"
const CREATE_IOC_COIN_FUTURE_ORDER ActionID = "30"
const CREATE_IOC_MARGIN_ORDER ActionID = "31"
const CREATE_MARKET_MARGIN_ORDER ActionID = "32"
const CREATE_MARKET_COINFUTURE_ORDER ActionID = "33"

// PositionID
const LONG PositionID = "LONG"
const SHORT PositionID = "SHORT"
const BOTH PositionID = "BOTH"
const NOPOSITIONID PositionID = ""

// SideID
const BUY SideID = "BUY"
const SELL SideID = "SELL"

// OrderTypeID
const OrderTypeLimit OrderTypeID = "LIMIT"
const OrderTypeMarket OrderTypeID = "MARKET"
const OrderTypeLimitMaker OrderTypeID = "LIMIT_MAKER"
const OrderTypeIOC OrderTypeID = "IOC"

// OrderStatusID (exchange side)
const NEW OrderStatusID = "NEW"         //The order has been accepted by the engine
const OPEN OrderStatusID = "OPEN"       // The order is open on the order book and is being worked
const AMENDED OrderStatusID = "AMENDED" //The amend order was accepted by the engine
const CANCELED OrderStatusID = "CANCELED"
const REJECTED OrderStatusID = "REJECTED"
const PARTIALLY_FILLED OrderStatusID = "PARTIALLY_FILLED"
const FILLED OrderStatusID = "FILLED"

// StatusID (gost side)
const Pending StatusID = 0 // the new order request was sent to exchange and waiting for response
const Normal StatusID = 1
const Canceling StatusID = 2
const Amending StatusID = 3 // the amend request was sent to exchange and waiting for response

// ResetID
const ResetSpotBookTick ResetID = "1"
const ResetSpotDepth ResetID = "2"
const ResetSpotTick ResetID = "3"
const ResetFutureBookTick ResetID = "4"
const ResetFutureDepth ResetID = "5"
const ResetFutureTick ResetID = "6"
const ResetSpotAccount ResetID = "7"
const ResetFutureAccount ResetID = "8"
const ResetAggTrade ResetID = "9"
const ResetMarginAccount ResetID = "10"
const ResetCoinFutureBookTick ResetID = "11"
const ResetCoinFutureAccount ResetID = "12"
const ResetSpotTrade ResetID = "13"
const ResetUniMarginAccount ResetID = "14"
const ResetFutureOrderbook ResetID = "15"
const ResetMarkPrice ResetID = "16"

// Interval
const Interval1s IntervalID = "1s"
const Interval1m IntervalID = "1m"
const Interval3m IntervalID = "3m"
const Interval5m IntervalID = "5m"
const Interval15m IntervalID = "15m"
const Interval30m IntervalID = "30m"
const Interval1h IntervalID = "1h"
const Interval2h IntervalID = "2h"
const Interval4h IntervalID = "4h"
const Interval6h IntervalID = "6h"
const Interval8h IntervalID = "8h"
const Interval12h IntervalID = "12h"
const Interval1d IntervalID = "1d"
const Interval3d IntervalID = "3d"
const Interval1w IntervalID = "1w"
const Interval1M IntervalID = "1M"

const CancelPlace AmendTypeID = "CancelPlace"   // cancel the order and place a new order
const Amendment AmendTypeID = "Amendment"       // amend the order with partially filled
const NotSupported AmendTypeID = "NotSupported" // don't support amend
