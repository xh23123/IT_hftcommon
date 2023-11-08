package common

type ExchangeID string

type AccountIdx int

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

// ExchangeID
const BINANCEID ExchangeID = "1"
const OKEXID ExchangeID = "3"
const COINBASEID ExchangeID = "4"
const COINEX ExchangeID = "5"
const BYBIT ExchangeID = "6"
const ETHEREUMDEX ExchangeID = "7"
const KRAKEN ExchangeID = "8"
const MOCKEXCHANGE ExchangeID = "9"

// TransactionID
const MixID TransactionID = "0"
const SpotID TransactionID = "1"
const FutureID TransactionID = "2"
const MarginID TransactionID = "3"
const CoinFutureID TransactionID = "4"
const UniMarginID TransactionID = "5"

// DataID
const AccountUpdateID DataID = "0"
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

// OrderStatusID
const NEW OrderStatusID = "NEW"
const OPEN OrderStatusID = "OPEN"
const CANCELED OrderStatusID = "CANCELED"
const REJECTED OrderStatusID = "REJECTED"
const EXPIRED OrderStatusID = "EXPIRED"
const PARTIALLY_FILLED OrderStatusID = "PARTIALLY_FILLED"
const FILLED OrderStatusID = "FILLED"

// StatusID
const Pending StatusID = 0
const Normal StatusID = 1
const Canceling StatusID = 2

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
