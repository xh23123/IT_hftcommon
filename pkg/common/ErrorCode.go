package common

type ErrorCode int

const ERRORCODE_WRONGTYPE ErrorCode = -1          //not gostError
const ERRORCODE_CANCEL_REJECTED ErrorCode = 2     //cancel failed
const ERRORCODE_ORDER_REJECTED ErrorCode = 3      //order failed
const ERRORCODE_CANCEL_ALL_REJECTED ErrorCode = 4 //cancel all failed
const ERRORCODE_AMEND_REJECTED ErrorCode = 5      //amend failed
const ERRORCODE_RISK_REJECTED ErrorCode = 6       //risk control rejected

type ReasonCode int

const REASON_UNKNOWN ReasonCode = 0           //unknown reason,should let developer know and fix it
const REASON_LIMIT_BREACH ReasonCode = 1      //limit breach. strategy should stop and retry later
const REASON_TIMEOUT ReasonCode = 2           //rest time out. strategy could retry
const REASON_PARAM_INVALID ReasonCode = 3     //param invalid, e.g. price is 0,price exceed limit,symbol not exist, account not exist
const REASON_NOTENOUGH_BALANCE ReasonCode = 4 // not enough balance, e.g. not enough balance to buy, not enough balance to pay fee, not enough margin
const REASON_ORDER_NOT_EXIST ReasonCode = 5   // order not exist, e.g. cancel order not exist, amend order not exist
const REASON_SYSTEM ReasonCode = 6            // system error, e.g. network error, db error, etc. shoult not retry
const REASON_RISK_CONTROL ReasonCode = 7
