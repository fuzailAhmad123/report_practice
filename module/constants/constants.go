package constants

const GROUP_BY string = "groupby"
const METRICS string = "metrics"
const START string = "start"
const END string = "end"
const ORG_ID string = "org_id"
const AD_ID string = "ad_id"
const DATE string = "date"
const ID string = "_id"
const BETS string = "bets"
const WINS string = "wins"

// operators constants
const GREATER_THAN_EQUALS string = "$gte"
const LESSER_THAN_EQUALS string = "$lte"
const INCLUDES string = "$in"
const SUM string = "$sum"

const PROJECTION_PREFIX string = "$_id."

// allowerd values for reprting.
var ALLOWED_GROUP_BY_FOR_REPORTING []string = []string{ORG_ID, AD_ID, DATE}
var ALLOWED_METRICS_FOR_REPORTING []string = []string{BETS, WINS}

// allowed retriver source type
const CLICKHOUSE string = "clickhouse"
const MONGO_LIVE string = "mongo_live"
const MONGO_SNAP string = "mongo_snap"
const BIGQUERY string = "bigquery"
const REDIS string = "redis"
