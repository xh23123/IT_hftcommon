package common

type ErrorCode int

const ERRORCODE_WRONGTYPE ErrorCode = -1             //not gostError
const ERRORCODE_CANCEL_ORDER_NOT_EXIST ErrorCode = 1 //cancel failed due to no such order
const ERRORCODE_CANCEL_REJECTED ErrorCode = 2        //cancel failed
const ERRORCODE_ORDER_REJECTED ErrorCode = 3         //order failed
const ERRORCODE_CANCEL_ALL_REJECTED ErrorCode = 4    //order not exist

type ReasonCode int

const REASON_UNKNOWN ReasonCode = 0           //unknown reason,should let developer know and fix it
const REASON_LIMIT_BREACH ReasonCode = 1      //limit breach. strategy should stop and retry later
const REASON_TIMEOUT ReasonCode = 2           //rest time out. strategy could retry
const REASON_PARAM_INVALID ReasonCode = 3     //param invalid, e.g. price is 0,price exceed limit,symbol not exist, account not exist
const REASON_NOTENOUGH_BALANCE ReasonCode = 4 // not enough balance, e.g. not enough balance to buy, not enough balance to pay fee, not enough margin
